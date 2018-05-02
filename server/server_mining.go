package server

import (
	"fmt"
	"time"
	"github.com/michain/dotcoin/chain"
)

func LoopMining(bc *chain.Blockchain){
	return
	for{
		fmt.Println("mining begin...")
		b, err := runMining(bc)
		if err != nil{
			//TODO log err info
		}
		fmt.Println("mining end", "["+string(b.Hash) + ", "+string(b.PrevBlockHash)+"]")
		fmt.Println("wait for next mining...")
		time.Sleep(180 * time.Second)
	}
}

// RunMining run mine block with mempool
func runMining(bc *chain.Blockchain) (*chain.Block, error){
	var newBlock *chain.Block
	if curTXMemPool.Count() >= 1 && len(minerAddress) > 0 {
	MineTransactions:
		var txs []*chain.Transaction

		for _, tx := range curTXMemPool.TxDescs() {
			if bc.VerifyTransaction(tx) {
				txs = append(txs, tx)
			}
		}

		if len(txs) == 0 {
			//TODO log err info
			return nil, ErrorAllTXInvalid
		}

		//reward miningAddress in this node
		cbTx := chain.NewCoinbaseTX(minerAddress, "", coinbaseReward)
		txs = append(txs, cbTx)

		//rebuild utxo set
		newBlock = bc.MineBlock(txs)

		bc.GetUTXOSet().Rebuild()

		fmt.Println("New block is mined!")

		for _, tx := range txs {
			if !tx.IsCoinBase() {
				curTXMemPool.RemoveTransaction(tx, false)
			}
		}

		for _, node := range curAddrManager.GetAddresses() {
			if node != listenAddress {
				//TODO: send inv?
				hash := newBlock.Hash
				fmt.Println("sendInv block", hash)
				//sendInv(node, "block", [][]byte{newBlock.Hash})
			}
		}

		if curTXMemPool.Count() > 0 {
			goto MineTransactions
		}
	}else{
		fmt.Println("no tx to mine")
	}

	return newBlock, nil
}
