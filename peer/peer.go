package peer

import (
	"net"
	"time"
	"sync"
	"fmt"
	"io"
	"errors"
)

var currentNode *Node
var lock = &sync.Mutex{}
var reqID int64

const maxBackupSeedLen = 10 // the max length of the seed backups

const maxSeedFailedRetry = 3 // the max retry times when a seed failed to connect

const syncBackupSeedInterval = 30 // the seed backup interval

const pingInterval = 30 // second
const maxPingErrorAllowed = 8

const maxBackupSeedAlive = 240

const maxResendStayTime = 120
const minResendStayTime = 20

type(
	// Node is the local server
	Node struct {
		nodeAddr string

		// the source seed addr
		sourceAddr string

		// the current seed
		seedAddr string
		seedConn net.Conn

		// ping、pong status with the current seed
		// when ping failed a few times, we need to connect to another seed node
		seedPingPong bool

		// backup seeds list
		// when failed to connect to source seed,will try backup seed
		seedBackup []*Seed

		// the nodes which use our node as the current seed
		downstreamNodes map[string]net.Conn

		// outer application channels
		send chan interface{}
		recv chan interface{}
	}

	// Seed structure
	Seed struct {
	addr  string
	retry int
	}
)

// NewNode return new &node with laddr, saddr, send, recv
func NewNode(laddr, saddr string, send, recv chan interface{}) *Node{
	return &Node{
		nodeAddr:    laddr,
		sourceAddr:  saddr,
		downstreamNodes: make(map[string]net.Conn),
		send:        send,
		recv:        recv,
	}
}



/*
laddr: our node's listen addr
saddr: the source seed addr
send: outer application pushes messages to this channel
recv: outer application receives messages from this channel
*/
func StartNode(laddr, saddr string, send, recv chan interface{}) error {
	if laddr == "" {
		return errors.New("must set local listen addr")
	}

	currentNode = NewNode(laddr, saddr, send, recv)

	// start tcp listening
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return err
	}

	// listen laddr wait for new message
	// wait for downstream nodes to connect
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("accept downstream node error：", err)
				continue
			}
			go currentNode.receiveMessage(conn, false, false)
		}
	}()

	// receive outer application's message,and route to the seed node and the downstream nodes
	go localSend(currentNode)

	// resend the unsent messages(these messages didn't receive a matched ack from target node)
	go resend(currentNode)

	// the main logic of seed manage
	if saddr != "" {
		// start to loop ping with the seed
		// when ping failed a few times, we need to connect to another seed node
		go currentNode.loopPing()

		err := currentNode.connectSeed(saddr)
		if err != nil {
			fmt.Printf("failed to connecto the seed%s：%v\n", saddr, err)
			return err
		}

		// start to sync the backup seed from current seed
		//the backup seeds are those nodes who directly connected with the current seed
		go currentNode.syncBackupSeed()

		// start to receive messages from the current seed
		currentNode.receiveMessage(currentNode.seedConn, true, false)

	SourceSeedTry:
	// although we disconnected with the source seed
	// but,here we want retry source seed for a few times(n) first
		n := 0
		for {
			if n > maxSeedFailedRetry {
				break
			}

			err := currentNode.connectSeed(saddr)
			if err != nil {
				n++
				goto CONTINUE
			}
			currentNode.receiveMessage(currentNode.seedConn, true, false)
			//when successfully connected, the counter will be reset to 0
			n = 0
		CONTINUE:
			time.Sleep(3 * time.Second)
		}

		// after retry several times with the source seed,now we want connect with our backup seeds
		for {
			if len(currentNode.seedBackup) <= 0 {
				// if there is no backup seed,we will go back to the source seed
				fmt.Printf("no backup seed exist now\n")
				break
			}

			// here is one important thing to notice
			//if stepBack is setted to 'true', we will go back to source seed retrys again

			// why?
			//because, at times, the big cluster will divided into few smaller clusters, the smaller
			// ones will not perceive each other, so we need a way to combine smaller ones to a
			// larger one, this is why we will go back to retry the source seed after some time.

			//and this stepBack action only happend when we has connected to backup seeds
			stepBack := currentNode.connectBackSeeds()
			if stepBack {
				fmt.Println("step back to the source seed")
				goto SourceSeedTry
			}
		}

		// go back to try source seed
		goto SourceSeedTry
	}

	select {}
}

// receiveMessage receive messages from remote node
func (node *Node) receiveMessage(conn net.Conn, isSeed bool, needStepBack bool) bool {
	var addr string
	// close the connection
	defer func() {
		conn.Close()
		// if the node is in downstream, then remove
		if addr != "" {
			fmt.Printf("remote downstream node %s close the connection\n", addr)
			node.deleteDownstreamNode(addr)
		}

		// if the node is the seed, then reset
		if isSeed {
			fmt.Printf("remote seed node %s close the connection\n", node.seedAddr)
			node.seedConn = nil
			node.seedAddr = ""
		}
	}()

	// the step back has been mentioned above
	start := time.Now().Unix()
	for {
		if needStepBack {
			now := time.Now().Unix()
			// A connection which connected to backup seed ,will maintain no more than 240 second
			if now-start > maxBackupSeedAlive {
				return true
			}
		}

		r, err:=ReadConnRequest(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Println("decode message error：", err)
			}
			break
		}
		a, err := r.handle(node, conn)
		if err != nil {
			fmt.Println("handle message error：", err)
			break
		}
		// update the downstream node's listen addr
		if a != "" {
			addr = a
		}
	}

	return false
}

// getConnByAddr get conn by addr
func (node *Node) getConnByAddr(addr string) net.Conn {
	if addr == node.seedAddr {
		return node.seedConn
	}

	conn, ok := node.downstreamNodes[addr]
	if ok {
		return conn
	}

	return nil
}


// delete node from downstream nodes
func (node *Node) deleteDownstreamNode(addr string) {
	lock.Lock()
	delete(node.downstreamNodes, addr)
	lock.Unlock()
}

// dial to remote node
func (node *Node) dialToNode(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// connect to the seed node
func (node *Node) connectSeed(addr string) error {
	conn, err := node.dialToNode(addr)
	if err != nil {
		return err
	}

	node.seedConn = conn
	node.seedAddr = addr
	fmt.Printf("connect to the seed %s successfully\n", addr)
	return nil
}

// sync backup seed from current seed
// the backup seeds are those nodes who directly connected with the current seed
func (node *Node) syncBackupSeed() {
	// waiting for node's initialization
	time.Sleep(100 * time.Millisecond)
	go func() {
		for {
			if node.seedConn != nil {
				r := &Request{
					Command: SyncBackupSeeds,
					Data:    node.nodeAddr,
				}
				WriteConnRequest(node.seedConn, r)
			}
			time.Sleep(syncBackupSeedInterval * time.Second)
		}
	}()
}


func (node *Node) loopPing() {
	pingNum := 0
	for {
		if node.seedPingPong {
			pingNum = 0
			node.seedPingPong = false
			continue
		}
		if pingNum >= maxPingErrorAllowed {
			// when the ping failed several times, we will choose another seed to connnect
			if node.seedConn != nil {
				node.seedConn.Close()
				node.seedConn = nil
				node.seedAddr = ""
			}
			pingNum = 0
			continue
		}

		if node.seedConn != nil {
			r := &Request{
				Command: ServerPing,
				Data:    node.nodeAddr,
			}
			WriteConnRequest(node.seedConn, r)
			pingNum++
		}
		time.Sleep(pingInterval * time.Second)
	}
}


func (node *Node) connectBackSeeds() bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("a critical error happens when connecto to backup seed", err)
		}
	}()

	for i, seed := range node.seedBackup {
		exist := false
		var err error
		var stepBack bool

		// a node can't  appear in seedBackup and downstream at the same time
		for addr := range node.downstreamNodes {
			if addr == seed.addr {
				node.seedBackup = append(node.seedBackup[:i], node.seedBackup[i+1:]...)
				exist = true
			}
		}
		if exist {
			fmt.Printf("a conflict between backupSeeds and downstream,so the backup seed is deleted：%s\n", seed.addr)
			continue
		}

		// seed connection retries can't exceed the upper limit
		if seed.retry > maxSeedFailedRetry {
			fmt.Printf("seed %sretries exceed the limit\n", seed.addr)
			node.seedBackup = append(node.seedBackup[:i], node.seedBackup[i+1:]...)
			goto CONTINUE
		}
		err = node.connectSeed(seed.addr)
		if err != nil {
			seed.retry++
			fmt.Printf("reconnect to seed %v error: %v\n", seed, err)
			goto CONTINUE
		}

		stepBack = node.receiveMessage(node.seedConn, true, true)
		// go back to source seed
		if stepBack {
			return true
		}
		// if a seed was successfully connected, the retry counter will be reset to 0
		seed.retry = 0
	CONTINUE:
		time.Sleep(3 * time.Second)
	}

	return false
}