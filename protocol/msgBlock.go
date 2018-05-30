package protocol

import "github.com/michain/dotcoin/chain"

type MsgBlock struct {
	NetMessage
	Block *chain.Block
}

func (msg *MsgBlock) Command() string {
	return CmdBlock
}

func NewMsgBlock(block *chain.Block) *MsgBlock {
	return &MsgBlock{
		Block:block,
	}
}