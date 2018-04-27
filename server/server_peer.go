package server

import (
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/protocol"
)

type serverPeer struct{

}

// OnGetAddr is invoked when a peer receives a getaddr bitcoin message
// and is used to provide the peer with known addresses from the address
// manager.
func (sp *serverPeer) OnGetAddr(p *peer.Peer, msg *protocol.MsgGetAddr) {
	// Get the current known addresses from the address manager.
	// TODO use AddrManager replace slice
	addrCache := knownNodes

	// Push the addresses.
	p.PushAddrMsg(addrCache)
}
