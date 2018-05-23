package sync

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/chain"
	"github.com/michain/dotcoin/mempool"
	"github.com/michain/dotcoin/logx"
	"github.com/michain/dotcoin/util/hashx"
)

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	Chain        *chain.Blockchain
	TxMemPool    *mempool.TxPool
	Peer		 *peer.Peer
	MaxPeers     int
}

type SyncManager struct{
	chain          *chain.Blockchain
	txMemPool      *mempool.TxPool
	peer			*peer.Peer
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
		peerStates: 	 	make(map[string]*peerSyncState),
		requestedTxs:	 	make(map[hashx.Hash]struct{}),
		requestedBlocks:	make(map[hashx.Hash]struct{}),
		msgChan:			make(chan interface{}, config.MaxPeers * 5),
	}
	return &sm, nil
}

func (manager *SyncManager) StartSync(){
	for {
		select {
		case m := <-manager.msgChan:
			switch msg := m.(type) {
			case *protocol.MsgVersion:
				manager.handleMsgVersion(msg)
			case *protocol.MsgInv:
				manager.handleMsgInv(msg)
			case *protocol.MsgGetBlocks:
				manager.handleMsgGetBlocks(msg)
			default:
				logx.Warnf("Invalid message type in sync msg chan: %T", msg)
			}

		case <-manager.quitSign:
			logx.Trace("SyncManager handle message done")
			return
		}
	}
}

// haveInventory check inv is exists
func (manager *SyncManager) haveInventory(inv *protocol.InvInfo) (bool, error) {
	switch inv.Type {
	case protocol.InvTypeBlock:
		return manager.chain.HaveBlock(inv.Hash)
	case protocol.InvTypeTx:
		//check tx-mempool
		if manager.txMemPool.HaveTransaction(inv.Hash) {
			return true, nil
		}

		// Check if the transaction exists from the point of view of the
		// end of the main chain.
		entry, err := manager.chain.FindTransaction(&inv.Hash)
		if err != nil {
			return false, err
		}

		return entry != nil, nil
	}
	//unsupported type
	return false, nil
}




func (manager *SyncManager) HandleInv(msg *protocol.MsgInv) {
	manager.msgChan <- msg
}

func (manager *SyncManager) HandleVersion(msg *protocol.MsgVersion){
	manager.msgChan <- msg
}

func (manager *SyncManager) HandleGetBlocks(msg *protocol.MsgGetBlocks){
	manager.msgChan <- msg
}