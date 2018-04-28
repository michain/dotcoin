package protocol

import "testing"

func TestNewMsgAddr(t *testing.T) {
	msg := NewMsgAddr()
	msg.SetNeedBroadcast(true)
	t.Log(msg.NeedBroadcast())
}
