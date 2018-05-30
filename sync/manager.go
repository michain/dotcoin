package sync

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/chain"
	"github.com/michain/dotcoin/mempool"
	"github.com/michain/dotcoin/util/hashx"
	"github.com/michain/dotcoin/addr"
)

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	Chain        *chain.Blockchain
	TxMemPool    *mempool.TxPool
	Peer		 *peer.Peer
	AddrManager  *addr.AddrManager
	MaxPeers     int
}

type SyncManager struct{
	chain          *chain.Blockchain
	txMemPool      *mempool.TxPool
	peer			*peer.Peer
	addrManager     *addr.AddrManager
	peerStates      map[string]*peerSyncState

	requestedTxs   map[hashx.Hash]struct{}
	requestedBlocks map[hashx.Hash]struct{}

	msgChan        chan interface{}
	quitSign       chan struct{}
}

// New return a new SyncManager with sync config
func New(config *Config) (*SyncManager, error) {
	sm := SyncManager{
		chain:           	config.Chain,
		txMemPool:       	config.TxMemPool,
		peer:				config.Peer,
		addrManager:		config.AddrManager,
		peerStates: 	 	make(map[string]*peerSyncState),
		requestedTxs:	 	make(map[hashx.Hash]struct{}),
		requestedBlocks:	make(map[hashx.Hash]struct{}),
		msgChan:			make(chan interface{}, config.MaxPeers * 5),
	}
	return &sm, nil
}

// StartSync start loop sync handle
func (manager *SyncManager) StartSync(){
	manager.loopHandle()
}

// haveInventory check inv is exists
func (manager *SyncManager) haveInventory(inv *protocol.InvInfo) (bool, error) {
	switch inv.Type {
	case protocol.InvTypeBlock:
		return manager.chain.HaveBlock(&inv.Hash)
	case protocol.InvTypeTx:
		//check tx-mempool
		if manager.txMemPool.HaveTransaction(inv.Hash.String()) {
			return true, nil
		}

		// Check if the transaction exists from the point of view of the
		// end of the main chain.
		entry, err := manager.chain.FindTransaction(&inv.Hash)
		if err != nil {
			if err == chain.ErrorNotFoundTransaction{
				return false, nil
			}
			return false, err
		}

		return entry != nil, nil
	}
	//unsupported type
	return false, nil
}

// getPeerState get peerState with peer's Addr
func (manager *SyncManager) getPeerState(peerAddr string) *peerSyncState{
	state, exists := manager.peerStates[peerAddr]
	if !exists{
		 state= &peerSyncState{
			setInventoryKnown: newInventorySet(maxInventorySize),
			requestedTxns:     make(map[hashx.Hash]struct{}),
			requestedBlocks:   make(map[hashx.Hash]struct{}),
		}
		manager.peerStates[peerAddr] = state
	}
	return state
}

// AddPeerState create new remote peer state
// if already exists, not cover it
func (manager *SyncManager) AddPeerState(remoteAddr string){
	if _, exists:=manager.peerStates[remoteAddr];!exists {
		manager.peerStates[remoteAddr] = &peerSyncState{
			setInventoryKnown: newInventorySet(maxInventorySize),
			requestedTxns:     make(map[hashx.Hash]struct{}),
			requestedBlocks:   make(map[hashx.Hash]struct{}),
		}
	}
}
