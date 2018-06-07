package sync

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/logx"
	"fmt"
	"github.com/michain/dotcoin/util/hashx"
	"github.com/michain/dotcoin/chain"
)

func (manager *SyncManager) HandleMessage(msg protocol.Message){
	manager.msgChan <- msg
}

func (manager *SyncManager) loopHandle(){
	for {
		select {
		case m := <-manager.msgChan:
			switch msg := m.(type) {
			case *protocol.MsgVersion:
				manager.handleMsgVersion(msg)
			case *protocol.MsgGetAddr:
				manager.handleMsgGetAddr(msg)
			case *protocol.MsgAddr:
				manager.handleMsgAddr(msg)
			case *protocol.MsgInv:
				manager.handleMsgInv(msg)
			case *protocol.MsgGetBlocks:
				manager.handleMsgGetBlocks(msg)
			case *protocol.MsgGetData:
				manager.handleMsgGetData(msg)
			case *protocol.MsgBlock:
				manager.handleMsgBlock(msg)
			case *protocol.MsgTx:
				manager.handleMsgTx(msg)
			default:
				logx.Warnf("Invalid message type in sync msg chan: %T", msg)
			}

		case <-manager.quitSign:
			logx.Trace("SyncManager handle message done")
			return
		}
	}
}

// handleMsgAddr handles Addr messages from other node.
func (manager *SyncManager) handleMsgAddr(msg *protocol.MsgAddr){
	needSync := []string{}
	for _, addr:=range msg.AddrList{
		if !manager.addrManager.HasAddress(addr){
			manager.addrManager.AddAddress(addr)
			needSync = append(needSync, addr)
		}
	}
	logx.Infof("handleMsgAddr get Addresses %d from peer:%v", len(msg.AddrList), msg.GetFromAddr())

	//notify other node
	if len(needSync) < 0 {
		msgSend := protocol.NewMsgAddr()
		msgSend.SetFromAddr(msg.GetFromAddr())
		msgSend.AddrList = needSync
		msgSend.AddrFrom = msg.AddrFrom
		manager.peer.SendRouteMessage(msgSend)
	}
}

// handleMsgGetAddr handles GetAddr messages from other node.
func (manager *SyncManager) handleMsgGetAddr(msg *protocol.MsgGetAddr){
	// Get the current known addresses from the address manager.
	addrCache := manager.addrManager.GetAddresses()

	msgSend := protocol.NewMsgAddr()
	msgSend.SetFromAddr(msg.GetFromAddr())
	msgSend.AddrList = addrCache
	manager.peer.PushAddrMsg(msgSend)
	logx.Infof("handleMsgGetAddr Send Addresses %d from peer:%v", len(addrCache), msg.GetFromAddr())
}

// handleVerionMsg handles version messages from other node.
// check best block height
func (manager *SyncManager) handleMsgVersion(msg *protocol.MsgVersion){
	//TODO Add remote Timestamp -> AddTimeData
	manager.AddPeerState(msg.GetFromAddr())

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

// handleInvMsg handles inv messages from other node.
// handle the inventory message and act GetData message
func (manager *SyncManager) handleMsgInv(msg *protocol.MsgInv) {

	if len(msg.InvList) <= 0 {
		logx.Warnf("SyncManager:handleInvMsg received empty inv list")
		return
	}

	state:= manager.getPeerState(msg.GetFromAddr())

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
			logx.Errorf("SyncManager:handleInvMsg [%v] Unexpected failure when checking for existing inventory [%s]", "handleInvMsg", err)
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

// handleMsgGetBlocks handles getblocks messages from other node.
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

// handleMsgGetData handles getdata messages from other node.
func (manager *SyncManager) handleMsgGetData(msg *protocol.MsgGetData){
	logx.Debugf("SyncManager.handleMsgGetData peer:%v msg:%v", manager.peer.GetListenAddr(), *msg)
	for _, iv := range msg.InvList {
		var err error
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
		case protocol.InvTypeTx:
			tx, exists:=manager.txMemPool.GetTransaction(iv.Hash.String())
			if !exists{
				tx, err = manager.chain.FindTransaction(&iv.Hash)
				if err != nil{
					logx.Errorf("SyncManager.handleMsgGetData get inv tx error %v %v peer:%v", iv.Hash, err, manager.peer.GetListenAddr())
					//TODO add to requery queue
					continue
				}
			}

			msgSend := protocol.NewMsgTx(tx)
			msgSend.AddrFrom = msg.GetFromAddr()
			err = manager.peer.PushTx(msgSend)
		default:
			logx.Warnf("SyncManager.handleMsgGetData Unknown type in getdata request %d",
				iv.Type)
			continue
		}
	}
}

// handleMsgBlock handles block messages from other node.
func (manager *SyncManager) handleMsgBlock(msg *protocol.MsgBlock){
	logx.Debugf("SyncManager.handleMsgBlock peer:%v msg:%v txs:%v", manager.peer.GetListenAddr(), msg.Block.GetHash(), len(msg.Block.Transactions))
	hash := msg.Block.GetHash()

	// Add the block to the known inventory for the peer.
	iv := protocol.NewInvInfo(protocol.InvTypeBlock, *hash)
	peerState:=manager.getPeerState(msg.GetFromAddr())
	peerState.AddKnownInventory(iv)

	// Add block to block chain
	isMain, isOrphanBlock, err:=manager.chain.ProcessBlock(msg.Block)
	logx.Info("SyncManager.handleMsgBlock ProcessBlock ", isMain, isOrphanBlock, err)

	// Notify signal to stop current mining
	if err == nil{
		manager.chain.TerminationMine()
	}

	// if add block success, remove repetition txs in txmempool
	// must remove main tx and orphan tx
	if err == nil || err == chain.ErrorAlreadyExistsBlock{
		//remove same tx in txmempool
		for _, tx:=range msg.Block.Transactions{
			manager.txMemPool.RemoveTransaction(tx)
			manager.txMemPool.RemoveOrphan(tx)
		}
	}

	// Notify other node which related of current node
	if err == nil{
		msgSend := protocol.NewMsgInv()
		msgSend.AddInvInfo(protocol.NewInvInfo(protocol.InvTypeBlock, *msg.Block.GetHash()))
		msgSend.AddrFrom = msg.AddrFrom
		manager.peer.SendRouteMessage(msgSend)
	}

}

// handleMsgTx handles tx messages from other node.
func (manager *SyncManager) handleMsgTx(msg *protocol.MsgTx){
	logx.Debugf("SyncManager.handleMsgTx peer:%v msg:%v", manager.peer.GetListenAddr(), msg.Tx.GetHash())
	hash := msg.Tx.GetHash()

	// Add the tx to the known inventory for the peer.
	iv := protocol.NewInvInfo(protocol.InvTypeTx, *hash)

	peerState:=manager.getPeerState(msg.GetFromAddr())
	peerState.AddKnownInventory(iv)

	// Add block to block chain
	err:=manager.txMemPool.ProcessTransaction(msg.Tx, true)
	if err != nil{
		logx.Errorf("SyncManager.handleMsgTx ProcessTransaction tx:%v error", msg.Tx.GetHash())
	}else{
		logx.Infof("SyncManager.handleMsgTx ProcessTransaction tx:%v success", msg.Tx.GetHash())
	}

	// Notify other node which related of current node
	if err == nil{
		msgSend := protocol.NewMsgInv()
		msgSend.AddInvInfo(protocol.NewInvInfo(protocol.InvTypeTx, *msg.Tx.GetHash()))
		msgSend.AddrFrom = msg.AddrFrom
		manager.peer.SendRouteMessage(msgSend)
	}

}