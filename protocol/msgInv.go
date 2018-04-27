package protocol

import (
	"github.com/michain/dotcoin/config/chainhash"
	"fmt"
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

type InvVect struct {
	Type string        // Type of data
	Hash chainhash.Hash // Hash of the data
}

// NewInvVect returns a new InvVect using the provided type and hash.
func NewInvVect(typ string, hash *chainhash.Hash) *InvVect {
	return &InvVect{
		Type: typ,
		Hash: *hash,
	}
}

type MsgInv struct {
	netMessage
	InvList []*InvVect
}

// AddInvVect adds an inventory vector to the message.
func (msg *MsgInv) AddInvVect(iv *InvVect) error {
	if len(msg.InvList)+1 > MaxInvPerMsg {
		str := fmt.Sprintf("too many invvect in message [max %v]",
			MaxInvPerMsg)
		return messageError("MsgInv.AddInvVect", str)
	}

	msg.InvList = append(msg.InvList, iv)
	return nil
}

// Command returns the protocol command string
func (msg *MsgInv) Command() string {
	return CmdInv
}

