package peer

import (
	"time"
)

type Packet struct {
	Addr   string
	retrys int
}

var sendPackets = make(map[int64][]*Packet)
var sendDatas = make(map[int64]Request)

// the outer application send messages
func localSend(node *Node) {
	for {
		select {
		case raw := <-node.send:
			now := time.Now().UnixNano()
			r := Request{
				ID:      now,
				Command: NormalRequest,
				Data:    raw,
				From:    node.nodeAddr,
			}

			lock.Lock()
			sendPackets[r.ID] = make([]*Packet, 0)
			sendDatas[r.ID] = r
			lock.Unlock()
			n := 0
			if node.seedAddr != "" {
				// send to the seed
				WriteConnRequest(node.seedConn, r)
				lock.Lock()
				sendPackets[r.ID] = append(sendPackets[r.ID], &Packet{
					Addr: node.seedAddr,
				})
				lock.Unlock()
				n++
			}

			// send to the downstream
			for addr, conn := range node.downstreamNodes {
				WriteConnRequest(conn, r)
				lock.Lock()
				sendPackets[r.ID] = append(sendPackets[r.ID], &Packet{
					Addr: addr,
				})
				lock.Unlock()
				n++
			}

			// nothing happend, do some sweeping work.
			if n == 0 {
				lock.Lock()
				delete(sendPackets, r.ID)
				delete(sendDatas, r.ID)
				lock.Unlock()
			}
		}
	}
}

// when receive remote node's messages, we will route to other nodes and the outer application
func routeSend(node *Node, r *Request) {
	now := time.Now().UnixNano()
	newReq := Request{
		ID:      now,
		Command: NormalRequest,
		Data:    r.Data,
		From:    node.nodeAddr,
	}

	lock.Lock()
	sendPackets[newReq.ID] = make([]*Packet, 0)
	sendDatas[newReq.ID] = newReq
	lock.Unlock()

	n := 0

	//if message not from seed node, we will send message to our seed node
	if r.From != node.seedAddr && node.seedAddr != "" {
		WriteConnRequest(node.seedConn, newReq)
		lock.Lock()
		sendPackets[newReq.ID] = append(sendPackets[newReq.ID], &Packet{
			Addr: node.seedAddr,
		})
		lock.Unlock()
		n++
	}

	//send message to all downstream nodes
	for addr, conn := range node.downstreamNodes {
		if r.From != addr && addr != "" {
			WriteConnRequest(conn, newReq)
			lock.Lock()
			sendPackets[newReq.ID] = append(sendPackets[newReq.ID], &Packet{
				Addr: addr,
			})
			lock.Unlock()
			n++
		}
	}

	// nothing happend, do some sweeping work.
	if n == 0 {
		lock.Lock()
		delete(sendPackets, newReq.ID)
		delete(sendDatas, newReq.ID)
		lock.Unlock()
	}

	// send to the outer application
	node.recv <- r
}

func deleteResend(rid int64, from string) {
	lock.Lock()
	ps, ok := sendPackets[rid]
	lock.Unlock()
	if !ok {
		return
	}

	for i, p := range ps {
		if p.Addr == from {
			ps = append(ps[:i], ps[i+1:]...)
			break
		}
	}
	if len(ps) != 0 {
		lock.Lock()
		sendPackets[rid] = ps
		lock.Unlock()
		return
	}

	lock.Lock()
	delete(sendPackets, rid)
	delete(sendDatas, rid)
	lock.Unlock()
}

// periodically resend the messages
func resend(node *Node) {
	for {
		now := time.Now().Unix()
		lock.Lock()
		for rid, ps := range sendPackets {
			// if the message stays too long,we will delete it directly
			if now-(rid/1e9) > maxResendStayTime {
				delete(sendPackets, rid)
				delete(sendDatas, rid)
				continue
			}
			// the message must stays for some time to resend
			if now-(rid/1e9) > minResendStayTime {
				r, ok := sendDatas[rid]
				if ok {
					for i, p := range ps {
						conn := node.getConnByAddr(p.Addr)
						if conn == nil {
							// the conn is empty,delete the message
							// TODO log the error addr info
							ps = append(ps[:i], ps[i+1:]...)
							continue
						}
						err := WriteConnRequest(conn, r)
						if err != nil {
							// the conn is broken, delete the message
							ps = append(ps[:i], ps[i+1:]...)
							continue
						}
					}
				}
			}

			if len(ps) == 0 {
				delete(sendPackets, rid)
				delete(sendDatas, rid)
			} else {
				sendPackets[rid] = ps
			}
		}
		lock.Unlock()

		time.Sleep(10 * time.Second)
	}
}



