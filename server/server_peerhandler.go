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
	logx.Tracef("ServerHandler OnGetAddr peer:%v remote:%v peers:%v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, handler.server.AddrManager.GetAddresses())
	handler.server.SyncManager.HandleMessage(msg)
}

// OnAddr is invoked when a peer receives an addr message.
func (handler *MessageHandler) OnAddr(msg *protocol.MsgAddr) {
 	logx.Tracef("ServerHandler OnAddr peer:%v remote:%v peers:%v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, handler.server.AddrManager.GetAddresses())
	handler.server.SyncManager.HandleMessage(msg)
}

// OnInv is invoked when a peer receives an inv message.
func (handler *MessageHandler) OnInv(msg *protocol.MsgInv) {
	logx.Tracef("ServerHandler OnInv peer:%v remote:%v invs:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, len(msg.InvList))
	handler.server.SyncManager.HandleMessage(msg)
}

// OnVersion is invoked when a peer receives an ver message
func (handler *MessageHandler) OnVersion(msg *protocol.MsgVersion){
	logx.Tracef("ServerHandler OnVersion peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	//add addrManager
	handler.server.AddrManager.AddAddress(msg.GetFromAddr())
	handler.server.SyncManager.HandleMessage(msg)
}

// OnGetBlocks is invoked when a peer receives an getblocks message
func (handler *MessageHandler) OnGetBlocks(msg *protocol.MsgGetBlocks){
	logx.Tracef("ServerHandler OnGetBlocks peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	handler.server.SyncManager.HandleMessage(msg)
}

// OnGetData is invoked when a peer receives an getdata message
func (handler *MessageHandler) OnGetData(msg *protocol.MsgGetData){
	logx.Tracef("ServerHandler OnGetData peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	handler.server.SyncManager.HandleMessage(msg)
}

// OnBlock is invoked when a peer receives an block message
func (handler *MessageHandler) OnBlock(msg *protocol.MsgBlock){
	logx.Tracef("ServerHandler OnBlock peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	handler.server.SyncManager.HandleMessage(msg)
}

// OnTx is invoked when a peer receives an tx message
func (handler *MessageHandler) OnTx(msg *protocol.MsgTx){
	logx.Tracef("messageHandler OnTx peer:%v remote:%v version:%+v", handler.server.Peer.GetListenAddr(), msg.AddrFrom, msg.ProtocolVersion)
	handler.server.SyncManager.HandleMessage(msg)
}
