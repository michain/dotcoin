package peer

import (
	"github.com/michain/dotcoin/config/chainhash"
	"github.com/michain/dotcoin/protocol"
	"math/rand"
)

// Peer extends the node to maintain state shared by the server
type Peer struct {
	node *Node
	sendQueue chan interface{}
	receiveQueue chan interface{}
	messageHandle *MessageHandle
	continueHash   *chainhash.Hash
}


func StartPeer(listenAddr, seedAddr string) (*Peer, error){
	p := new(Peer)
	p.sendQueue = make(chan interface{}, 1)
	p.receiveQueue = make(chan interface{}, 1)
	p.node = NewNode(listenAddr, seedAddr, p.sendQueue, p.receiveQueue)


	return p, nil
}


// SendData send message to all downstream nodes and seed node
func (p *Peer) SendData(msg protocol.Message){
	p.sendQueue <- msg
}

// ReciveMessage recive message from net node
func (p *Peer) ReciveMessage() protocol.Message{
	var msgData interface{}
	for {
		msgData =<- p.receiveQueue
		if msg, ok := msgData.(protocol.Message); !ok {
			//TODO log error msg
		}else{
			return msg
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

	p.SendData(msg)
	return nil
}