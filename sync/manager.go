package sync

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/chain"
	"github.com/michain/dotcoin/mempool"
	"github.com/michain/dotcoin/logx"
	"github.com/michain/dotcoin/config/chainhash"
	"fmt"
)

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	Chain        *chain.Blockchain
	TxMemPool    *mempool.TxPool
	Peer			*peer.Peer
	MaxPeers     int
}

type SyncManager struct{
	chain          *chain.Blockchain
	txMemPool      *mempool.TxPool
	peer			*peer.Peer
	peerStates      map[string]*peerSyncState

	requestedTxs   map[chainhash.Hash]struct{}
	requestedBlocks map[chainhash.Hash]struct{}

	msgChan        chan interface{}
	quitSign       chan struct{}
}

// New return a new SyncManager with sync config
func New(config *Config) (*SyncManager, error) {
	sm := SyncManager{
		chain:           	config.Chain,
		txMemPool:       	config.TxMemPool,
		peer:config.Peer,
		peerStates: 	 	make(map[string]*peerSyncState),
		requestedTxs:	 	make( map[chainhash.Hash]struct{}),
		requestedBlocks:	 make( map[chainhash.Hash]struct{}),
	}
	return &sm, nil
}

func (manager *SyncManager) Start(){
	for {
		select {
		case m := <-manager.msgChan:
			switch msg := m.(type) {
			case *protocol.MsgVersion:
				manager.handleVerionMsg(msg)
			case *protocol.MsgInv:
				manager.handleInvMsg(msg)
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
		return manager.chain.HaveBlock(inv.Hash.CloneBytes())
	case protocol.InvTypeTx:
		//check tx-mempool
		if manager.txMemPool.HaveTransaction(inv.Hash.CloneBytes()) {
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


// handleInvMsg handles inv messages from other peer.
// handle the inventory message and act GetData message
func (manager *SyncManager) handleInvMsg(msg *protocol.MsgInv) {

	state, exists := manager.peerStates[msg.GetFromAddr()]
	if !exists {
		logx.Warnf("Received inv message from unknown peer %s", msg.GetFromAddr())
		return
	}


	// Attempt to find the final block in the inventory list
	lastBlock := -1
	invInfos := msg.InvList
	for i := len(invInfos) - 1; i >= 0; i-- {
		if invInfos[i].Type == protocol.InvTypeBlock {
			lastBlock = i
			break
		}
	}
	fmt.Println("SyncManager:handleInvMsg", lastBlock)

	for _, iv := range invInfos {
		// Ignore unsupported inventory types.
		switch iv.Type {
		case protocol.InvTypeBlock:
		case protocol.InvTypeTx:
		default:
			continue
		}

		state.AddKnownInventory(iv)

		haveInv, err := manager.haveInventory(iv)
		if err != nil {
			logx.Errorf("[%v] Unexpected failure when checking for existing inventory [%s]", "handleInvMsg", err)
			continue
		}

		if !haveInv{
			if iv.Type == protocol.InvTypeTx {
				//TODO if  transaction has been rejected, skip it
			}
			// Add inv to the request inv queue.
			state.requestInvQueue = append(state.requestInvQueue, iv)
			continue
			if iv.Type == protocol.InvTypeBlock {

			}
		}
	}

	numRequestInvs := 0
	requestQueue := state.requestInvQueue
	logx.DevPrintf("handleInvMsg requestQueue %v", requestQueue)
	// Request GetData command
	getDataMsg := protocol.NewMsgGetData()
	getDataMsg.AddrFrom = msg.GetFromAddr()
	for _, iv:=range state.requestInvQueue{
		switch iv.Type {
		case protocol.InvTypeBlock:
			if _, exists := manager.requestedBlocks[iv.Hash]; !exists {
				manager.requestedBlocks[iv.Hash] = struct{}{}
				err := getDataMsg.AddInvInfo(iv)
				if err != nil{
					break
				}
				numRequestInvs++
			}
		case protocol.InvTypeTx:
			if _, exists := manager.requestedTxs[iv.Hash]; !exists {
				manager.requestedBlocks[iv.Hash] = struct{}{}
				err := getDataMsg.AddInvInfo(iv)
				if err != nil{
					break
				}
				numRequestInvs++
			}
		}
		if numRequestInvs >= protocol.MaxInvPerMsg {
			break
		}
	}


	state.requestInvQueue = []*protocol.InvInfo{}
	if len(getDataMsg.InvList) > 0 {
		manager.peer.SendSingleMessage(getDataMsg)
	}
}

func (manager *SyncManager) handleVerionMsg(msg *protocol.MsgVersion){
	//TODO Add remote Timestamp
	manager.peerStates[msg.GetFromAddr()] = &peerSyncState{
		setInventoryKnown: newInventorySet(maxInventorySize),
		requestedTxns:   make(map[chainhash.Hash]struct{}),
		requestedBlocks: make(map[chainhash.Hash]struct{}),
	}

}

func (manager *SyncManager) HandleInv(msg *protocol.MsgInv) {
	manager.msgChan <- msg
}

func (manager *SyncManager) HandleVersion(msg *protocol.MsgVersion){
	manager.msgChan <- msg
}