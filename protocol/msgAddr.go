package protocol

import "fmt"

// MaxAddrPerMsg is the maximum number of addresses that can be in a single
// bitcoin addr message (MsgAddr).
const MaxAddrPerMsg = 1000

type MsgAddr struct {
	netMessage
	AddrList []string
}

// AddAddress adds a known active peer to the message.
func (msg *MsgAddr) AddAddress(na string) error {
	if len(msg.AddrList)+1 > MaxAddrPerMsg {
		str := fmt.Sprintf("too many addresses in message [max %v]",
			MaxAddrPerMsg)
		return messageError("MsgAddr.AddAddress", str)
	}

	msg.AddrList = append(msg.AddrList, na)
	return nil
}

// AddAddresses adds multiple known active peers to the message.
func (msg *MsgAddr) AddAddresses(netAddrs ...string) error {
	for _, na := range netAddrs {
		err := msg.AddAddress(na)
		if err != nil {
			return err
		}
	}
	return nil
}

// ClearAddresses removes all addresses from the message.
func (msg *MsgAddr) ClearAddresses() {
	msg.AddrList = []string{}
}

// Command returns the protocol command string
func (msg *MsgAddr) Command() string {
	return CmdAddr
}

// NewMsgAddr returns a new addr message
func NewMsgAddr() *MsgAddr {
	return  &MsgAddr{
		AddrList: make([]string, 0, MaxAddrPerMsg),
	}
}
