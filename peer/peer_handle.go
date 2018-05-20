package peer

import (
	"github.com/michain/dotcoin/protocol"
)
type MessageHandle interface {

	SetPeer(p *Peer)

	// OnGetAddr is invoked when a peer receives a getaddr message.
	OnGetAddr(msg *protocol.MsgGetAddr)

	// OnAddr is invoked when a peer receives an addr message.
	OnAddr(msg *protocol.MsgAddr)

	// OnInv is invoked when a peer receives an inv message.
	OnInv(msg *protocol.MsgInv)

	// OnVersion is invoked when a peer receives an ver message
	OnVersion(msg *protocol.MsgVersion)
}
