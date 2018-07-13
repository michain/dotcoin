package server

import (
	"time"
	"github.com/michain/dotcoin/chain"
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/logx"
)

const(
	MinimumTxInNewBlock = 1
	ZeroTxInNewBlock = 0
)


func (s *Server) LoopMining(){
	for{
		runMining(s)
		time.Sleep(10 * time.Second)
	}
}

// RunMining run mine block with mempool
func runMining(s *Server) (*chain.Block, error){
	var newBlock *chain.Block

	minimumTxCount := MinimumTxInNewBlock
	if s.Config.EnabledEmptyTxInNewBlock{
		minimumTxCount = ZeroTxInNewBlock
	}
	if s.TXMemPool.Count() >= minimumTxCount && len(s.minerAddress) > 0 {
		var txs []*chain.Transaction

		//reward miningAddress in this node
		cbTx := chain.NewCoinbaseTX(s.minerAddress, "", coinbaseReward)
		txs = append(txs, cbTx)

		for _, tx := range s.TXMemPool.TxDescs() {
			if s.BlockChain.VerifyTransaction(tx) {
				txs = append(txs, tx)
			}
		}

		if minimumTxCount != ZeroTxInNewBlock && len(txs) == 0 {
			//TODO log err info
			return nil, ErrorAllTXInvalid
		}


		//rebuild utxo set
		isSuccess := false
		newBlock, isSuccess = s.BlockChain.MineBlock(txs)

		if !isSuccess{
			logx.Warnf("MineBlock failde")
		}else {

			s.BlockChain.GetUTXOSet().Rebuild()

			logx.Info("MineBlock Success", " hash: ", newBlock.GetHash().String(), " reward: ", coinbaseReward, " txs: ", len(newBlock.Transactions))

			for _, tx := range txs {
				if !tx.IsCoinBase() {
					s.TXMemPool.RemoveTransaction(tx)
					s.TXMemPool.RemoveOrphan(tx)
				}
			}

			// Broadcast inv message to other node
			hash := newBlock.GetHash()
			inv := protocol.NewInvInfo(protocol.InvTypeBlock, *hash)
			msgSend := protocol.NewMsgInv()
			msgSend.AddInvInfo(inv)
			s.Peer.BroadcastMessage(msgSend)
			logx.Debugf("Server Mining Broadcast block [%v] inv message", hash.String())
		}
	}else{
		logx.Tracef("MineBlock failde: no tx to mine")
	}

	return newBlock, nil
}
