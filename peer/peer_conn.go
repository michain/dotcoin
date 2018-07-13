package peer

import (
	"fmt"
	"net"
	"strconv"
	"github.com/michain/dotcoin/logx"
	"encoding/gob"
)

// Request 节点之间交换的数据结构
type Request struct {
	ID      int64
	Command int
	From    string
	Data    interface{}
}

func init(){
	gob.Register(Request{})
}

const (
	SendRequest 		 = 1
	ServerPing            	 = 2 // ping to the seed
	ServerPong            	 = 3 // pong to the ping
	SyncBackupSeeds       	 = 4 // query for the backup seeds
	BackupSeeds           	 = 5 // return the backup seeds
	RequestReceived 		 = 6 // ack
)

func (r *Request) handleConn(node *Node, conn net.Conn) (string, error) {
	switch r.Command {
	case RequestReceived:
		// delete the message from resend queue
		deleteResend(r.ID, r.From)
	case SendRequest:
		//logx.Tracef("peer.handleConn.SendRequest [%v] => [%v]", r.From, node.listenAddr)
		// send to the outer application
		node.recv <- r
		//send ack message
		sendNormalRequestReceived(r, node, conn)
	case SyncBackupSeeds:
		// the address of the requester
		fromAddr := r.Data.(string)

		// filter the adjacent nodes from current seed and the downstream nodes
		// Avoid forming a dead loop, seed addr mustn't be equal to from addr
		var addrs []string
		if node.seedAddr != "" && node.seedAddr != fromAddr {
			addrs = append(addrs, node.seedAddr)
		}
		for addr := range node.downstreamNodes {
			if len(addrs) < maxBackupSeedLen && addr != fromAddr {
				addrs = append(addrs, addr)
			}
		}

		WriteConnRequest(conn, Request{
			Command: BackupSeeds,
			Data:    addrs,
		})
	case BackupSeeds:
		if r.Data == nil{

		}else {
			addrs := r.Data.([]string)

			for _, addr1 := range addrs {
				if addr1 == "" {
					continue
				}
				// the strategy of seeds update
				// if the upper limit of the seedBackup is not reached, we can append the new addr to the seedBackup
				// otherwise, we need to replace those nodes whose connection retries bigger than the maxRetry,
				// with the new seed
				exist := false
				maxRetry := 0
				for _, seed := range node.seedBackup {
					if seed.retry > maxRetry {
						maxRetry = seed.retry
					}
					if addr1 == seed.addr {
						exist = true
						break
					}
				}

				if !exist {
					if len(node.seedBackup) >= maxBackupSeedLen {
						if maxRetry <= maxSeedFailedRetry {
							break
						}
						for i, seed := range node.seedBackup {
							if seed.retry > maxSeedFailedRetry {
								node.seedBackup[i] = &Seed{
									addr:  addr1,
									retry: 0,
								}
							}
						}
					} else {
						node.seedBackup = append(node.seedBackup, &Seed{
							addr:  addr1,
							retry: 0,
						})
					}
				}

			}
		}

		logx.Debugf("peer.BackupSeeds source: %s,current：%s,backup：%v,downsteam：%v", node.sourceAddr, node.seedAddr, getSeedAddrs(node.seedBackup), node.downstreamNodes)
	case ServerPing:
		//logx.DevDebugf("peer.ServerPing [%v] => [%v] %v", r.From, node.listenAddr, *r)
		// a downstream node sends its address to its seed node
		addr, ok := r.Data.(string)
		if ok {
			// we need to add the downstream node
			lock.Lock()
			node.downstreamNodes[addr] = conn
			lock.Unlock()
			// a node can't appear in downstream and seedBackup at the same time
			for i, seed := range node.seedBackup {
				if seed.addr == addr {
					node.seedBackup = append(node.seedBackup[:i], node.seedBackup[i+1:]...)
					break
				}
			}

			WriteConnRequest(conn, Request{
				Command: ServerPong,
				From:    node.listenAddr,
			})

			return addr, nil
		}

	case ServerPong:
		//logx.DevDebugf("peer.seedPingPong [%v] => [%v] %v", r.From, node.listenAddr, *r)
		node.seedPingPong = true
	default:
		fmt.Println("unrecognized message type：", r.Command)
	}

	return "", nil
}

func getSeedAddrs(seeds []*Seed) []string {
	addrs := make([]string, len(seeds))
	for i, seed := range seeds {
		addrs[i] = seed.addr + "/" + strconv.Itoa(seed.retry)
	}

	return addrs
}
