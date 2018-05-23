package peer

import (
	"github.com/michain/dotcoin/logx"
	"github.com/michain/dotcoin/protocol"
	"reflect"
)

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
		case protocol.MsgVersion:
			msg.AddrFrom = req.From
			p.messageHandler.OnVersion(&msg)
		case protocol.MsgGetBlocks:
			msg.AddrFrom = req.From
			p.messageHandler.OnGetBlocks(&msg)
		case protocol.MsgGetData:
			msg.AddrFrom = req.From
			p.messageHandler.OnGetData(&msg)
		default:
			logx.Errorf("Received unhandled message of type %v "+
				"from %v [%v]", reflect.TypeOf(req.Data), p, msg)
		}
	}
}
