package chain

import (
	"testing"
	"fmt"
	"os"

	"github.com/michain/dotcoin/wallet"
)

var(
	god *wallet.Wallet
	from *wallet.Wallet
	to *wallet.Wallet
	wallets *wallet.WalletSet
	godBlockChain *Blockchain
	godNodeID string
	err error
)

func init(){
	godNodeID = "god"

	wallets, err = wallet.LoadWallets(godNodeID)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if len(wallets.Wallets) == 0{
		god = wallets.CreateWallet()
		from = wallets.CreateWallet()
		to = wallets.CreateWallet()
		wallets.SaveToFile()

		fmt.Println("god", string(god.GetAddress()))
		fmt.Println("from", string(from.GetAddress()))
		fmt.Println("to", string(to.GetAddress()))
	}else{
		godAddress := "19m3x9GjCtpNjCi1k7RdZz1W9HkxAun6Xv"
		fromAddress := "1AwHUCD6RDtszfQ7uTF5EmZgkdMtm2Usgi"
		toAddress := "1LMirU374Spvx3jg1MVTQ3qttwnypJRkPM"
		god = wallets.GetWallet(godAddress)
		if god == nil{
			fmt.Println("create god wallet error")
			os.Exit(-1)
		}
		from = wallets.GetWallet(fromAddress)
		if from == nil{
			fmt.Println("create from wallet error")
			os.Exit(-1)
		}
		to = wallets.GetWallet(toAddress)
		if to == nil{
			fmt.Println("create to wallet error")
			os.Exit(-1)
		}
	}

	if godBlockChain, err = LoadBlockChain(godNodeID); err!=nil{
		fmt.Println("LoadBlockChain error", err)
		return
	}
	if godBlockChain == nil{
		godBlockChain = CreateBlockchain(true, string(god.GetAddress()), godNodeID)
	}
}

func TestGetBalance(t *testing.T){
	needBalance := 100
	balance := 0
	pubKeyHash := wallet.GetPubKeyHashFromAddress(god.GetAddress())
	UTXOs := godBlockChain.GetUTXOSet().FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Println("GetBalance", string(god.GetAddress()), balance)
	if balance == needBalance{
		t.Log(balance)
	}else{
		t.Error("get balance not match!", needBalance, balance)
	}
}

func TestNewUTXOTransaction(t *testing.T) {
	trans := NewUTXOTransaction(god, from.GetStringAddress(), 10, godBlockChain.GetUTXOSet(), make(map[string]*Transaction))
	fmt.Println(*trans)
}

func TestFindTransaction(t *testing.T){
}

func TestNewBlock(t *testing.T) {

}