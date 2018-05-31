package server

import (
	"fmt"
	"time"
	"github.com/michain/dotcoin/chain"
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/logx"
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
		var txs []*chain.Transaction

		//reward miningAddress in this node
		cbTx := chain.NewCoinbaseTX(s.minerAddress, "", coinbaseReward)
		txs = append(txs, cbTx)

		for _, tx := range s.TXMemPool.TxDescs() {
			if s.BlockChain.VerifyTransaction(tx) {
				txs = append(txs, tx)
			}
		}

		if len(txs) == 0 {
			//TODO log err info
			return nil, ErrorAllTXInvalid
		}


		//rebuild utxo set
		isSuccess := false
		newBlock, isSuccess = s.BlockChain.MineBlock(txs)

		if !isSuccess{
			fmt.Println("MineBlock failde")
		}else {

			s.BlockChain.GetUTXOSet().Rebuild()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				if !tx.IsCoinBase() {
					s.TXMemPool.RemoveTransaction(tx, false)
				}
			}

			// Broadcast inv message to other node
			hash := newBlock.GetHash()
			inv := protocol.NewInvInfo(protocol.InvTypeBlock, *hash)
			msgSend := protocol.NewMsgInv()
			msgSend.AddInvInfo(inv)
			s.Peer.BroadcastMessage(msgSend)
			logx.Infof("Server Mining Broadcast block [%v] inv message", hash.String())

			//TODO: start next mining
			/*if s.TXMemPool.Count() > 0 {
				goto MineTransactions
			}*/
		}
	}else{
		fmt.Println("no tx to mine")
	}

	return newBlock, nil
}
