package server

import (
	"testing"
	"fmt"
	"os"
	"github.com/michain/dotcoin/chain"
	"bytes"
	"log"
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
