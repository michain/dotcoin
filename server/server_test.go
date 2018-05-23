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
)

const genNodeID = "3eb456d086f34118925793496cd20945"
var genServer *Server
var fromAddress = "16KKFgx41SEi2YwfCcJTroENrndeNoSYL7"

const NodeID_1 = "1eb456d086f34118925793496cd20945"
const NodeID_2 = "2eb456d086f34118925793496cd20945"
var Server_1 *Server
var Server_2 *Server

const genAddr = "127.0.0.1:2398"
const Addr1 = "127.0.0.1:2491"
const Addr2 = "127.0.0.1:2492"


func init(){
	var err error
	genServer, err = initServer(genNodeID, "", genAddr, genAddr, true)
	if err!=nil{
		fmt.Println(err)
		os.Exit(-1)
	}else{
		//start sync loop
		go genServer.SyncManager.StartSync()
		fmt.Println("Genesis server start", genServer.NodeID)
	}

	Server_1, err = initServer(NodeID_1, "", Addr1, genAddr, false)
	if err!=nil{
		fmt.Println(err)
		os.Exit(-1)
	}else{
		//start sync loop
		go Server_1.SyncManager.StartSync()
		fmt.Println("Server_1 server start", Server_1.NodeID)
	}

	Server_2, err = initServer(NodeID_2, "", Addr2, Addr1, false)
	if err!=nil{
		fmt.Println(err)
		os.Exit(-1)
	}else{
		//start sync loop
		go Server_2.SyncManager.StartSync()
		fmt.Println("Server_2 server start", Server_2.NodeID)
	}
}


func Test_StartPeer(t *testing.T){

	var err error

	go func() {
		genServer.listenPeer()
		if err != nil {
			t.Error("Seed Peer start error", err)
		} else {
			t.Log("Seed Peer start success")
		}
	}()

	go func() {
		time.Sleep(time.Second * 1)
		Server_1.listenPeer()
		if err != nil {
			t.Error("pnode1 Peer start error", err)
		} else {
			t.Log("pnode1 Peer start success")
		}
	}()

	go func() {
		time.Sleep(time.Second * 2)
		Server_2.listenPeer()
		if err != nil {
			t.Error("pnode2 Peer start error", err)
		} else {
			t.Log("pnode2 Peer start success")
		}
	}()

	/*go func(){
		time.Sleep(time.Second * 3)
		msg := protocol.NewMsgAddr()
		msg.AddrList = []string{"127.0.0.1:1", "127.0.0.1:2", "127.0.0.1:3"}
		p_2_1.BroadcastMessage(msg)
	}()*/

	go func(){
		time.Sleep(time.Second * 5)
		msg := protocol.NewMsgVersion(Server_1.BlockChain.GetBestHeight())
		Server_1.Peer.PushVersion(msg)



		/*iv := protocol.NewInvInfo(protocol.InvTypeTx, chainhash.ZeroHash())
		msgInv := protocol.NewMsgInv()
		msgInv.AddInvInfo(iv)
		p_2_1.BroadcastMessage(msgInv)*/

	}()




	for{
		select{}
	}


}

func Test_runMining(t *testing.T){
	//add tx

	fromWallet := genServer.Wallets.GetWallet(fromAddress)
	if fromWallet  == nil{
		fmt.Println("not exists [from] address")
		os.Exit(-1)
	}

	to := genServer.Wallets.CreateWallet().GetStringAddress()
	tx := chain.NewUTXOTransaction(fromWallet, to, 1, genServer.BlockChain.GetUTXOSet())
	fmt.Println("NewUTXOTransaction", tx.ID, tx.StringHash())
	//add TX to mempool
	_, err := genServer.TXMemPool.MaybeAcceptTransaction(tx, true, true)
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

	block, err := runMining(genServer)
	if err != nil{
		log.Panic("Mining err", err)
		return
	}

	fmt.Println(fromAddress, "Balance", genServer.BlockChain.GetBalance(fromAddress))

	lastBlock, err := genServer.BlockChain.GetLastBlock()
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
