package protocol

import "time"

type MsgVersion struct {
	netMessage
	ProtocolVersion int32

	// Time the message was generated.  This is encoded as an int64 on the wire.
	Timestamp time.Time

	// Last block Height seen by the generator of the version message.
	LastBlockHeight int32
}

func NewMsgVersion(lastBlockHeight int32) *MsgVersion{
	return &MsgVersion{
		ProtocolVersion: int32(ProtocolVersion),
		Timestamp:       time.Unix(time.Now().Unix(), 0),
		LastBlockHeight:       lastBlockHeight,
	}
}

// Command returns the protocol command string
func (msg *MsgVersion) Command() string {
	return CmdVersion
}
