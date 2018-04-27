package peer

import (
	"github.com/michain/dotcoin/protocol"
)
type MessageHandle struct {
	// OnGetAddr is invoked when a peer receives a getaddr bitcoin message.
	OnGetAddr func(p *Peer, msg *protocol.MsgGetAddr)

	// OnAddr is invoked when a peer receives an addr bitcoin message.
	OnAddr func(p *Peer, msg *protocol.MsgAddr)

	// OnInv is invoked when a peer receives an inv bitcoin message.
	OnInv func(p *Peer, msg *protocol.MsgInv)
}
