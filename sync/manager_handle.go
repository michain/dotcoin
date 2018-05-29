package sync

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/logx"
	"fmt"
	"github.com/michain/dotcoin/util/hashx"
)

func (manager *SyncManager) HandleMessage(msg protocol.Message){
	manager.msgChan <- msg
}

// handleInvMsg handles inv messages from other peer.
// handle the inventory message and act GetData message
func (manager *SyncManager) handleMsgInv(msg *protocol.MsgInv) {

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
	//TODO why calc lastBlock?
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

// handleVerionMsg handles version messages from other peer.
// check best block height
func (manager *SyncManager) handleMsgVersion(msg *protocol.MsgVersion){
	//TODO Add remote Timestamp -> AddTimeData
	manager.peerStates[msg.GetFromAddr()] = &peerSyncState{
		setInventoryKnown: newInventorySet(maxInventorySize),
		requestedTxns:   make(map[hashx.Hash]struct{}),
		requestedBlocks: make(map[hashx.Hash]struct{}),
	}
	if manager.chain.GetBestHeight() < msg.LastBlockHeight  {
		hashStop := hashx.ZeroHash()
		if  manager.chain.GetBestHeight() > 0 {
			//send getblocks message
			block, err := manager.chain.GetLastBlock()
			if err != nil {
				//TODO log get last block err
				logx.Error("handleMsgVersion::GetLastBlock error", err)
				return
			}
			hashStop = block.GetHash()
		}
		msgSend := protocol.NewMsgGetBlocks(*hashStop)
		msgSend.AddrFrom = msg.GetFromAddr()
		manager.peer.PushGetBlocks(msgSend)
	} else if manager.chain.GetBestHeight() > msg.LastBlockHeight   {
		//send version message
		msgSend := protocol.NewMsgVersion(manager.chain.GetBestHeight())
		msgSend.AddrFrom = msg.GetFromAddr()
		manager.peer.PushVersion(msgSend)
	}else{
		//TODO nothing to do?
		logx.Debug("handleMsgVersion: it's same height, so nothing to do")
	}
}

// handleMsgGetBlocks handles getblocks messages from other peer.
func (manager *SyncManager) handleMsgGetBlocks(msg *protocol.MsgGetBlocks){
	logx.Debugf("SyncManager.handleMsgGetBlocks peer:%v msg:%v", manager.peer.GetListenAddr(), *msg)
	block, err := manager.chain.GetLastBlock()
	if err != nil{
		//TODO log get last block err
		return
	}
	h := block.GetHash()
	hashes, err:= manager.chain.GetBlockHashes(h, msg.HashStop, protocol.MaxBlocksPerMsg)
	if err != nil{
		//TODO log get block hashes err
		return
	}

	//send blocks inv
	msgInv := protocol.NewMsgInv()
	msgInv.AddrFrom = msg.GetFromAddr()
	for _, hash:=range hashes{
		msgInv.AddInvInfo(protocol.NewInvInfo(protocol.InvTypeBlock, *hash))
	}
	manager.peer.SendSingleMessage(msgInv)
}

// handleMsgGetData handles getdata messages from other peer.
func (manager *SyncManager) handleMsgGetData(msg *protocol.MsgGetData){
	logx.Debugf("SyncManager.handleMsgGetData peer:%v msg:%v", manager.peer.GetListenAddr(), *msg)
	for _, iv := range msg.InvList {
		switch iv.Type {
		case protocol.InvTypeBlock:
			block, err :=manager.chain.GetBlock(iv.Hash.CloneBytes())
			if err != nil{
				logx.Errorf("SyncManager.handleMsgGetData get inv block error %v %v", iv.Hash, err)
				//TODO add to requery queue
				continue
			}
			msgSend := protocol.NewMsgBlock(block)
			msgSend.AddrFrom = msg.GetFromAddr()
			err = manager.peer.PushBlock(msgSend)
		default:
			logx.Warnf("SyncManager.handleMsgGetData Unknown type in getdata request %d",
				iv.Type)
			continue
		}
	}
}

// handleMsgGetData handles getdata messages from other peer.
func (manager *SyncManager) handleMsgBlock(msg *protocol.MsgBlock){
	logx.Debugf("SyncManager.handleMsgBlock peer:%v msg:%v txs:%v", manager.peer.GetListenAddr(), msg.Block.GetHash(), len(msg.Block.Transactions))
	hash := msg.Block.GetHash()

	// Add the block to the known inventory for the peer.
	iv := protocol.NewInvInfo(protocol.InvTypeBlock, *hash)
	peerState, exists:=manager.peerStates[msg.GetFromAddr()]
	if !exists{
		logx.Warnf("SyncManager.handleMsgBlock Received block message from unknown peer %s", msg.GetFromAddr())
		return
	}
	peerState.AddKnownInventory(iv)


	// Add block to block chain
	isMain, isOrphanBlock, err:=manager.chain.ProcessBlock(msg.Block)
	logx.Info("SyncManager.handleMsgBlock ProcessBlock ", isMain, isOrphanBlock, err)

}