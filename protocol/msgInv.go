package protocol

import (
	"fmt"
	"github.com/michain/dotcoin/util/hashx"
)

// Map of service flags back to their constant names for pretty printing.
const(
	InvTypeError = "ERROR"
	InvTypeTx = "MSG_TX"
	InvTypeBlock = "MSG_BLOCK"
)

const(
	// MaxInvPerMsg is the maximum number of inventory vectors that can be in a single inv message.
	MaxInvPerMsg = 50000
)

type InvInfo struct {
	Type string        // Type of data
	Hash hashx.Hash // Hash of the data
}

// NewInvInfo returns a new InvVect using the provided type and hash.
func NewInvInfo(typ string, hash hashx.Hash) *InvInfo {
	return &InvInfo{
		Type: typ,
		Hash: hash,
	}
}

type MsgInv struct {
	NetMessage
	InvList []*InvInfo
}

func NewMsgInv() *MsgInv{
	return &MsgInv{}
}

// AddInvVect adds an inventory info to the message.
func (msg *MsgInv) AddInvInfo(iv *InvInfo) error {
	if len(msg.InvList)+1 > MaxInvPerMsg {
		str := fmt.Sprintf("too many invinfo in message [max %v]",
			MaxInvPerMsg)
		return messageError("MsgInv.AddInvInfo", str)
	}

	msg.InvList = append(msg.InvList, iv)
	return nil
}

// Command returns the protocol command string
func (msg *MsgInv) Command() string {
	return CmdInv
}

