package sync

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/config/chainhash"
)

const(
	// maxKnownInventory is the maximum item's number to keep in the known inventory set.
	maxKnownInventory = 1000
)


type peerSyncState struct {
	syncCandidate   bool
	setInventoryKnown *inventorySet
	requestInvQueue    []*protocol.InvInfo
	requestedTxns   map[chainhash.Hash]struct{}
	requestedBlocks map[chainhash.Hash]struct{}
}

// AddKnownInventory adds the passed inventory to the cache of known inventory for the peer.
func (p *peerSyncState) AddKnownInventory(invVect *protocol.InvInfo) {
	p.setInventoryKnown.Add(invVect)
}


func newPeerSyncState() *peerSyncState{
	return &peerSyncState{
		setInventoryKnown:newInventorySet(maxKnownInventory),
		requestedTxns:make(map[chainhash.Hash]struct{}),
		requestedBlocks:make(map[chainhash.Hash]struct{}),
	}
}
