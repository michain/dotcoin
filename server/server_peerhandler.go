package server

import (
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/logx"
)

type MessageHandler struct{

}

func NewMessageHandler() *MessageHandler{
	return new(MessageHandler)
}

// OnGetAddr is invoked when a peer receives a getaddr message
func (mh *MessageHandler) OnGetAddr(p *peer.Peer, msg *protocol.MsgGetAddr) {
	logx.Debugf("messageHandler OnGetAddr %v", msg)
	// Get the current known addresses from the address manager.
	addrCache := currentAddrManager.GetAddresses()
	// Push the addresses.
	p.PushAddrMsg(addrCache)
}

// OnAddr is invoked when a peer receives an addr message.
func (mh *MessageHandler) OnAddr(p *peer.Peer, msg *protocol.MsgAddr) {
	for _, addr:=range msg.AddrList{
		currentAddrManager.AddAddress(addr)
	}
	logx.DevPrintf("MessageHandler OnAddr peer:%v from:%v peers:%v", p.GetListenAddr(), msg.AddrFrom, currentAddrManager.GetAddresses())
}

// OnInv is invoked when a peer receives an inv message.
func (mh *MessageHandler) OnInv(p *peer.Peer, msg *protocol.MsgInv) {
	logx.Debugf("messageHandler OnInv [%v]", msg)
}
