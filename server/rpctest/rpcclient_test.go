package rpctest

import "testing"

func TestCallCreateWallet(t *testing.T) {
	CallCreateWallet()
}


func TestCallListAddress(t *testing.T) {
	CallListAddress()
}

func TestCallSendTX(t *testing.T) {
	from:="16dE6XG9F6KQDAEduaeTxozNtfqcKfh3EG"
	to:="1H8iWFoaEN324YM1rLUTqYz32zkSt5hYft"

	errFrom := "16dE6XG9F6KQDAEduaeTxozNtfqcKfh3E2"
	err := CallSendTX(from, to)
	if err != nil{
		t.Error(err)
	}
	err = CallSendTX(errFrom, to)
	if err!= nil{
		t.Log("")
	}else{
		t.Error("no validate address!")
	}
}