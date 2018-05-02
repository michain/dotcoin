package protocol

import "time"

type MsgVersion struct {
	netMessage
	ProtocolVersion int32

	// Time the message was generated.  This is encoded as an int64 on the wire.
	Timestamp time.Time

	// Unique value associated with message that is used to detect self
	// connections.
	Nonce uint64

	// Last block seen by the generator of the version message.
	LastBlock int32
}

func NewMsgVersion(nonce uint64, lastBlock int32) *MsgVersion{
	return &MsgVersion{
		ProtocolVersion: int32(ProtocolVersion),
		Timestamp:       time.Unix(time.Now().Unix(), 0),
		Nonce:           nonce,
		LastBlock:       lastBlock,
	}
}

// Command returns the protocol command string
func (msg *MsgVersion) Command() string {
	return CmdVersion
}
