package protocol


type MsgGetAddr struct {
	NetMessage
}

// Command returns the protocol command string
func (msg *MsgGetAddr) Command() string {
	return CmdGetAddr
}

// NewMsgGetAddr returns a new getaddr message
func NewMsgGetAddr() *MsgGetAddr {
	return &MsgGetAddr{}
}
