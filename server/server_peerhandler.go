package server

import (
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/logx"
)

type MessageHandler struct{
	peer *peer.Peer
}

func NewMessageHandler(p *peer.Peer) *MessageHandler{
	return &MessageHandler{peer:p}
}

// OnGetAddr is invoked when a peer receives a getaddr message
func (handler *MessageHandler) OnGetAddr(msg *protocol.MsgGetAddr) {
	logx.Debugf("messageHandler OnGetAddr %v", msg)
	// Get the current known addresses from the address manager.
	addrCache := curAddrManager.GetAddresses()
	// Push the addresses.
	handler.peer.PushAddrMsg(addrCache)
}

// OnAddr is invoked when a peer receives an addr message.
func (handler *MessageHandler) OnAddr(msg *protocol.MsgAddr) {
	logx.DevPrintf("MessageHandler OnAddr peer:%v from:%v peers:%v", handler.peer.GetListenAddr(), msg.AddrFrom, curAddrManager.GetAddresses())
	for _, addr:=range msg.AddrList{
		curAddrManager.AddAddress(addr)
	}
}

// OnInv is invoked when a peer receives an inv message.
func (handler *MessageHandler) OnInv(msg *protocol.MsgInv) {
	logx.DevPrintf("messageHandler OnInv peer:%v from:%v invs:%+v", handler.peer.GetListenAddr(), msg.AddrFrom, msg.InvList[0])
	if len(msg.InvList) > 0 {
		curSyncManager.HandleInv(msg)
	}
	return
}
