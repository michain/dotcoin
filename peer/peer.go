package peer

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/util/hashx"
)

// Peer extends the node to maintain state shared by the server
type Peer struct {
	node *Node
	boardcastQueue chan interface{}
	singleQueue chan *SingleRequest
	receiveQueue chan *Request
	messageHandler MessageHandle
	continueHash   *hashx.Hash
}


func NewPeer(listenAddr, seedAddr string, msgHandler MessageHandle) *Peer{
	p := new(Peer)
	p.singleQueue = make(chan *SingleRequest, 10)
	p.boardcastQueue = make(chan interface{}, 10)
	p.receiveQueue = make(chan *Request, 10)
	p.messageHandler = msgHandler
	p.node = NewNode(listenAddr, seedAddr, p.boardcastQueue, p.receiveQueue, p.singleQueue)
	return p
}

func (p *Peer) StartListen() error{
	go p.ReciveMessage()
	return p.node.startNode()
}

func (p *Peer) GetSeedAddr() string{
	return p.node.seedAddr
}

// GetListenAddr get peer listen addr
func (p *Peer) GetListenAddr() string{
	return p.node.listenAddr
}

// BroadcastMessage send message to all downstream nodes and seed node
func (p *Peer) SendSingleMessage(msg protocol.Message){
	p.singleQueue <- &SingleRequest{
		Data:msg,
		FromAddr:msg.GetFromAddr(),
	}
}

// BroadcastMessage send message to all downstream nodes and seed node
func (p *Peer) BroadcastMessage(msg protocol.Message){
	p.boardcastQueue <- msg
}



