package protocol

import "fmt"

const defaultInvListAlloc = 1000

type MsgGetData struct {
	NetMessage
	InvList []*InvInfo
}

// AddInvInfo adds an inventory to the message.
func (msg *MsgGetData) AddInvInfo(iv *InvInfo) error {
	if len(msg.InvList)+1 > MaxInvPerMsg {
		str := fmt.Sprintf("too many inventory in message [max %v]",
			MaxInvPerMsg)
		return messageError("MsgGetData.AddInvVect", str)
	}

	msg.InvList = append(msg.InvList, iv)
	return nil
}

func (msg *MsgGetData) Command() string {
	return CmdGetData
}

func NewMsgGetData() *MsgGetData {
	return &MsgGetData{
		InvList: make([]*InvInfo, 0, defaultInvListAlloc),
	}
}

