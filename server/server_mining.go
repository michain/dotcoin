package server

import (
	"fmt"
	"time"
	"github.com/michain/dotcoin/chain"
)

func (s *Server) LoopMining(){
	return
	for{
		fmt.Println("mining begin...")
		b, err := runMining(s)
		if err != nil{
			//TODO log err info
		}
		fmt.Println("mining end", "["+string(b.Hash) + ", "+string(b.PrevBlockHash)+"]")
		fmt.Println("wait for next mining...")
		time.Sleep(180 * time.Second)
	}
}

// RunMining run mine block with mempool
func runMining(s *Server) (*chain.Block, error){
	var newBlock *chain.Block
	if s.TXMemPool.Count() >= 1 && len(s.minerAddress) > 0 {
	MineTransactions:
		var txs []*chain.Transaction

		for _, tx := range s.TXMemPool.TxDescs() {
			if s.BlockChain.VerifyTransaction(tx) {
				txs = append(txs, tx)
			}
		}

		if len(txs) == 0 {
			//TODO log err info
			return nil, ErrorAllTXInvalid
		}

		//reward miningAddress in this node
		cbTx := chain.NewCoinbaseTX(s.minerAddress, "", coinbaseReward)
		txs = append(txs, cbTx)

		//rebuild utxo set
		newBlock = s.BlockChain.MineBlock(txs)

		s.BlockChain.GetUTXOSet().Rebuild()

		fmt.Println("New block is mined!")

		for _, tx := range txs {
			if !tx.IsCoinBase() {
				s.TXMemPool.RemoveTransaction(tx, false)
			}
		}

		for _, node := range s.AddrManager.GetAddresses() {
			if node != s.ListenAddress {
				//TODO: send inv?
				hash := newBlock.Hash
				fmt.Println("sendInv block", hash)
				//sendInv(node, "block", [][]byte{newBlock.Hash})
			}
		}

		if s.TXMemPool.Count() > 0 {
			goto MineTransactions
		}
	}else{
		fmt.Println("no tx to mine")
	}

	return newBlock, nil
}
