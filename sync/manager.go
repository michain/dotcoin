package sync

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/chain"
	"github.com/michain/dotcoin/mempool"
	"github.com/michain/dotcoin/logx"
)

type SyncManager struct{
	chain          *chain.Blockchain
	txMemPool      *mempool.TxPool
}

// New return a new SyncManager with sync config
func New(config *Config) (*SyncManager, error) {
	sm := SyncManager{
		chain:           config.Chain,
		txMemPool:       config.TxMemPool,
	}
	return &sm, nil
}



// handleInvMsg handles inv messages from other peer.
// examine the inventory message and act getdata message
func (sm *SyncManager) handleInvMsg(p *peer.Peer, msg *protocol.MsgInv) {
	// Attempt to find the final block in the inventory list
	lastBlock := -1
	invInfos := msg.InvList
	for i := len(invInfos) - 1; i >= 0; i-- {
		if invInfos[i].Type == protocol.InvTypeBlock {
			lastBlock = i
			break
		}
	}

	for i, iv := range invInfos {
		// Ignore unsupported inventory types.
		switch iv.Type {
		case protocol.InvTypeBlock:
		case protocol.InvTypeTx:
		default:
			continue
		}

		//TODO check is exists inv
	}

	// Request getdata command
	gdmsg := wire.NewMsgGetData()
	logx.DevPrintf("handleInvMsg %v, %v", lastBlock, gdmsg)
}
