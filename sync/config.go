package sync

import (
	"github.com/michain/dotcoin/mempool"
	"github.com/michain/dotcoin/chain"
)

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	Chain        *chain.Blockchain
	TxMemPool    *mempool.TxPool
	MaxPeers     int
}

