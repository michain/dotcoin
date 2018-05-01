package server

import (
	"testing"
	"fmt"
	"os"
	"github.com/michain/dotcoin/chain"
	"bytes"
	"log"
	"time"
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/config/chainhash"
)

const nodeID = "3eb456d086f34118925793496cd20945"
var fromAddress = "16KKFgx41SEi2YwfCcJTroENrndeNoSYL7"


func init(){
	nodeID := "3eb456d086f34118925793496cd20945"
	err := initServer(nodeID, "", true)
	if err!=nil{
		fmt.Println(err)
		os.Exit(-1)
	}
}


func Test_StartPeer(t *testing.T){
	var seed = "127.0.0.1:2398"
	var pnode1 = "127.0.0.1:2391"
	var pnode1_1 = "127.0.0.1:2392"
	var pnode2 = "127.0.0.1:2491"
	var pnode2_1 = "127.0.0.1:2492"

	go func() {
		p := peer.NewPeer(seed, "", NewMessageHandler())

		err := p.StartListen()
		if err != nil {
			t.Error("Seed Peer start error", err)
		} else {
			t.Log("Seed Peer start success")
		}
	}()

	go func() {
		p := peer.NewPeer(pnode1, seed, NewMessageHandler())
		err := p.StartListen()
		if err != nil {
			t.Error("pnode1 Peer start error", err)
		} else {
			t.Log("pnode1 Peer start success")
		}
	}()

	go func() {
		p := peer.NewPeer(pnode1_1, pnode1, NewMessageHandler())
		err := p.StartListen()
		if err != nil {
			t.Error("pnode1_1 Peer start error", err)
		} else {
			t.Log("pnode1_1 Peer start success")
		}
	}()


	go func() {
		p := peer.NewPeer(pnode2, seed, NewMessageHandler())
		err := p.StartListen()
		if err != nil {
			t.Error("pnode2 Peer start error", err)
		} else {
			t.Log("pnode2 Peer start success")
		}
	}()

	var p_2_1 *peer.Peer
	go func() {
		var err error
		p_2_1 = peer.NewPeer(pnode2_1, pnode2, NewMessageHandler())
		err = p_2_1.StartListen()
		if err != nil {
			t.Error("pnode2_1 Peer start error", err)
		} else {
			t.Log("pnode2_1 Peer start success")
		}
	}()

	go func(){
		time.Sleep(time.Second * 3)
		msg := protocol.NewMsgAddr()
		msg.AddrList = []string{"127.0.0.1:1", "127.0.0.1:2", "127.0.0.1:3"}
		p_2_1.BroadcastMessage(msg)
	}()

	go func(){
		time.Sleep(time.Second * 6)
		iv := protocol.NewInvInfo(protocol.InvTypeTx, chainhash.ZeroHash())
		msg := protocol.NewMsgInv()
		msg.AddInvInfo(iv)

		p_2_1.BroadcastMessage(msg)
	}()


	for{
		select{}
	}


}

func Test_runMining(t *testing.T){
	//add tx

	fromWallet := currentWallets.GetWallet(fromAddress)
	if fromWallet  == nil{
		fmt.Println("not exists [from] address")
		os.Exit(-1)
	}

	to := currentWallets.CreateWallet().GetStringAddress()
	tx := chain.NewUTXOTransaction(fromWallet, to, 1, currentBlockChain.GetUTXOSet())
	fmt.Println("NewUTXOTransaction", tx.ID, tx.StringHash())
	//add TX to mempool
	_, err := txPool.MaybeAcceptTransaction(tx, true, true)
	if err != nil{
		fmt.Println(err)
	}

	/*
	logx.Debugf("Begin Range TxPool")
	for _, tx:=range txPool.TxDescs(){
		fmt.Println(tx.ID, tx.StringHash())
	}
	logx.Debugf("End Range TxPool")
	*/

	block, err := runMining(currentBlockChain)
	if err != nil{
		log.Panic("Mining err", err)
		return
	}

	fmt.Println(fromAddress, "Balance", currentBlockChain.GetBalance(fromAddress))

	lastBlock, err := currentBlockChain.GetLastBlock()
	if err != nil{
		fmt.Println("GetLastBlock err", err)
	}else{
		if bytes.Compare(block.Hash, lastBlock.Hash) != 0{
			fmt.Println("LastBlock not equal with MineBlock")
		}else {
			fmt.Println(lastBlock)
		}
	}
}
