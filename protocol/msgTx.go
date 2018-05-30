package protocol

import "github.com/michain/dotcoin/chain"

type MsgTx struct {
	NetMessage
	Tx *chain.Transaction
}

func (msg *MsgTx) Command() string {
	return CmdBlock
}

func NewMsgTx(tx *chain.Transaction) *MsgTx {
	return &MsgTx{
		Tx:tx,
	}
}
