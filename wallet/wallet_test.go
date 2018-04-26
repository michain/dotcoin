package wallet

import (
	"testing"
	"fmt"
	"github.com/michain/dotcoin/base58"
)

func Test_newKeyPair(t *testing.T){
	private, pubKey := newKeyPair()
	fmt.Println(private)
	fmt.Println(base58.Encode(pubKey))
}


func Test_newWallet(t *testing.T){
	w := newWallet()
	fmt.Println(string(w.GetAddress()))
}