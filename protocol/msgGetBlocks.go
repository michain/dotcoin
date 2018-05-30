package protocol

import "github.com/michain/dotcoin/util/hashx"

// MaxBlocksPerMsg is the max number of block hashes allowed
// per message.
const MaxBlocksPerMsg = 500

type MsgGetBlocks struct {
	NetMessage
	HashStop        hashx.Hash
}

func (msg *MsgGetBlocks) Command() string {
	return CmdGetBlocks
}

// NewMsgGetBlocks returns a new getblocks message
func NewMsgGetBlocks(hashStop hashx.Hash) *MsgGetBlocks {
	msg := &MsgGetBlocks{
		HashStop:           hashStop,
	}
	msg.ProtocolVersion = ProtocolVersion
	return msg
}