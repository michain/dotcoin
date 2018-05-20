package peer

import (
	"github.com/michain/dotcoin/config/chainhash"
	"github.com/michain/dotcoin/protocol"
	"math/rand"
	"github.com/michain/dotcoin/logx"
	"reflect"
)

// Peer extends the node to maintain state shared by the server
type Peer struct {
	node *Node
	boardcastQueue chan interface{}
	singleQueue chan *SingleRequest
	receiveQueue chan *Request
	messageHandler MessageHandle
	continueHash   *chainhash.Hash
}


func NewPeer(listenAddr, seedAddr string, msgHandler MessageHandle) *Peer{
	p := new(Peer)
	p.singleQueue = make(chan *SingleRequest, 10)
	p.boardcastQueue = make(chan interface{}, 10)
	p.receiveQueue = make(chan *Request, 10)
	p.messageHandler = msgHandler
	p.node = NewNode(listenAddr, seedAddr, p.boardcastQueue, p.receiveQueue, p.singleQueue)
	p.messageHandler.SetPeer(p)
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

// ReciveMessage recive message from net node
func (p *Peer) ReciveMessage() {
	for {
		req := <-p.receiveQueue
		//logx.DevPrintf("Received msgData type:%v", reflect.TypeOf(req.Data))
		if p.messageHandler == nil {
			logx.Error("Peer's messageHandler is nil!")
			break
		}
		switch msg :=  req.Data.(type) {
		case protocol.MsgAddr:
			msg.AddrFrom = req.From
			p.messageHandler.OnAddr(&msg)
		case protocol.MsgInv:
			msg.AddrFrom = req.From
			p.messageHandler.OnInv(&msg)
		default:
			logx.Errorf("Received unhandled message of type %v "+
				"from %v [%v]", reflect.TypeOf(req.Data), p, msg)
		}
	}
}

func (p *Peer) PushAddrMsg(addresses []string) error {
	addressCount := len(addresses)

	// Nothing to send.
	if addressCount == 0 {
		return nil
	}

	msg := protocol.NewMsgAddr()
	msg.AddrList = make([]string, addressCount)
	copy(msg.AddrList, addresses)

	// Randomize the addresses sent if there are more than the maximum allowed.
	if addressCount > protocol.MaxAddrPerMsg {
		// Shuffle the address list.
		for i := 0; i < protocol.MaxAddrPerMsg; i++ {
			j := i + rand.Intn(addressCount-i)
			msg.AddrList[i], msg.AddrList[j] = msg.AddrList[j], msg.AddrList[i]
		}

		// Truncate it to the maximum size.
		msg.AddrList = msg.AddrList[:protocol.MaxAddrPerMsg]
	}

	//set single send
	msg.SetNeedBroadcast(false)

	p.SendSingleMessage(msg)
	return nil
}
