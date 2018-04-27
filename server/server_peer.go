package server

import (
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/protocol"
)

type serverHandle struct{

}

// OnGetAddr is invoked when a peer receives a getaddr message
func (sh *serverHandle) OnGetAddr(p *peer.Peer, msg *protocol.MsgGetAddr) {
	// Get the current known addresses from the address manager.
	addrCache := currentAddrManager.GetAddresses()

	// Push the addresses.
	p.PushAddrMsg(addrCache)
}
