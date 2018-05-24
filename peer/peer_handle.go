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

	// OnVersion is invoked when a peer receives an ver message
	OnVersion(msg *protocol.MsgVersion)

	// OnGetBlocks is invoked when a peer receives an getblocks message
	OnGetBlocks(msg *protocol.MsgGetBlocks)

	// OnGetData is invoked when a peer receives an getdata message
	OnGetData(msg *protocol.MsgGetData)

	// OnBlock is invoked when a peer receives an block message
	OnBlock(msg *protocol.MsgBlock)
}
