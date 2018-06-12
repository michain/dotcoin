package rpctest

import "testing"

func TestCallCreateWallet(t *testing.T) {
	CallCreateWallet()
}

func TestCallListMemPool(t *testing.T) {
	CallListMemPool()
}

func TestCallListAddress(t *testing.T) {
	CallListAddress()
}

func TestCallListBlocks(t *testing.T) {
	CallListBlocks()
}

func TestCallSendTX(t *testing.T) {
	from:="1HCGY3WD5UCFNxQyLodoPvSwhZDUQu3kCn"
	to:="1BozgKkxFjtW5RnozCLzmnKc1zsfYdd95q"

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