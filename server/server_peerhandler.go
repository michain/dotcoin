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
	logx.DevPrintf("MessageHandler OnGetAddr peer:%v remote:%v peers:%v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, handler.server.AddrManager.GetAddresses())
	// Get the current known addresses from the address manager.
	addrCache := handler.server.AddrManager.GetAddresses()
	// Push the addresses.
	handler.server.Peer.PushAddrMsg(addrCache)
}

// OnAddr is invoked when a peer receives an addr message.
func (handler *MessageHandler) OnAddr(msg *protocol.MsgAddr) {
 	logx.DevPrintf("MessageHandler OnAddr peer:%v remote:%v peers:%v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, handler.server.AddrManager.GetAddresses())
	for _, addr:=range msg.AddrList{
		handler.server.AddrManager.AddAddress(addr)
	}
}

// OnInv is invoked when a peer receives an inv message.
func (handler *MessageHandler) OnInv(msg *protocol.MsgInv) {
	logx.DevPrintf("messageHandler OnInv peer:%v remote:%v invs:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, len(msg.InvList))
	if len(msg.InvList) > 0 {
		handler.server.SyncManager.HandleMessage(msg)
	}
	return
}

// OnVersion is invoked when a peer receives an ver message
func (handler *MessageHandler) OnVersion(msg *protocol.MsgVersion){
	logx.DevPrintf("messageHandler OnVersion peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	//add addrManager
	handler.server.AddrManager.AddAddress(msg.GetFromAddr())
	handler.server.SyncManager.HandleMessage(msg)
}

// OnGetBlocks is invoked when a peer receives an getblocks message
func (handler *MessageHandler) OnGetBlocks(msg *protocol.MsgGetBlocks){
	logx.DevPrintf("messageHandler OnGetBlocks peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	handler.server.SyncManager.HandleMessage(msg)
}

// OnGetData is invoked when a peer receives an getdata message
func (handler *MessageHandler) OnGetData(msg *protocol.MsgGetData){
	logx.DevPrintf("messageHandler OnGetData peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	handler.server.SyncManager.HandleMessage(msg)
}

// OnBlock is invoked when a peer receives an block message
func (handler *MessageHandler) OnBlock(msg *protocol.MsgBlock){
	logx.DevPrintf("messageHandler OnBlock peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	handler.server.SyncManager.HandleMessage(msg)
}
