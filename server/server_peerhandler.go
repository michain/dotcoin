package server

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/logx"
)

type MessageHandler struct{
	server *Server
}

func NewMessageHandler(s *Server) *MessageHandler{
	return &MessageHandler{server:s}
}

// OnGetAddr is invoked when a peer receives a getaddr message
func (handler *MessageHandler) OnGetAddr(msg *protocol.MsgGetAddr) {
	logx.Debugf("messageHandler OnGetAddr %v", msg)
	// Get the current known addresses from the address manager.
	addrCache := handler.server.AddrManager.GetAddresses()
	// Push the addresses.
	handler.server.Peer.PushAddrMsg(addrCache)
}

// OnAddr is invoked when a peer receives an addr message.
func (handler *MessageHandler) OnAddr(msg *protocol.MsgAddr) {
 	logx.DevPrintf("MessageHandler OnAddr peer:%v from:%v peers:%v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, handler.server.AddrManager.GetAddresses())
	for _, addr:=range msg.AddrList{
		handler.server.AddrManager.AddAddress(addr)
	}
}

// OnInv is invoked when a peer receives an inv message.
func (handler *MessageHandler) OnInv(msg *protocol.MsgInv) {
	logx.DevPrintf("messageHandler OnInv peer:%v from:%v invs:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.InvList[0])
	if len(msg.InvList) > 0 {
		handler.server.SyncManager.HandleInv(msg)
	}
	return
}

// OnVersion is invoked when a peer receives an ver message
func (handler *MessageHandler) OnVersion(msg *protocol.MsgVersion){
	logx.DevPrintf("messageHandler OnVersion peer:%v from:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	//add addrManager
	handler.server.AddrManager.AddAddress(msg.GetFromAddr())
	handler.server.SyncManager.HandleVersion(msg)
}

// OnGetBlocks is invoked when a peer receives an getblocks message
func (handler *MessageHandler) OnGetBlocks(msg *protocol.MsgGetBlocks){
	logx.DevPrintf("messageHandler OnGetBlocks peer:%v from:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	//add addrManager
	handler.server.SyncManager.HandleGetBlocks(msg)
}