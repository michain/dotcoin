package peer

import (
	"github.com/michain/dotcoin/protocol"
)
type MessageHandle interface {
	// OnGetAddr is invoked when a peer receives a getaddr message.
	OnGetAddr(msg *protocol.MsgGetAddr)

	// OnAddr is invoked when a peer receives an addr message.
	OnAddr(msg *protocol.MsgAddr)

	// OnInv is invoked when a peer receives an inv message.
	OnInv(msg *protocol.MsgInv)
}
