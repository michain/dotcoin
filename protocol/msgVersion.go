package protocol

import "time"

type MsgVersion struct {
	NetMessage
	ProtocolVersion int32

	// Time the message was generated.  This is encoded as an int64 on the wire.
	Timestamp time.Time

	// Last block Height seen by the generator of the version message.
	LastBlockHeight int32

	// Last block hash seen by the generator of the version message.
	LastBlockHash []byte

	// Last block prev hash seen by the generator of the version message.
	LastBlockPrevHash []byte
}

func NewMsgVersion(lastBlockHeight int32, lastBlockHash, lastBlockPrevHash []byte) *MsgVersion{
	return &MsgVersion{
		ProtocolVersion: int32(ProtocolVersion),
		Timestamp:       time.Unix(time.Now().Unix(), 0),
		LastBlockHeight: lastBlockHeight,
		LastBlockHash:	 lastBlockHash,
		LastBlockPrevHash:lastBlockPrevHash,
	}
}

// Command returns the protocol command string
func (msg *MsgVersion) Command() string {
	return CmdVersion
}
