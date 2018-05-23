package sync

import (
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/util/hashx"
)

const(
	// maxKnownInventory is the maximum item's number to keep in the known inventory set.
	maxKnownInventory = 1000
)


type peerSyncState struct {
	syncCandidate   bool
	setInventoryKnown *inventorySet
	requestInvQueue    []*protocol.InvInfo
	requestedTxns   map[hashx.Hash]struct{}
	requestedBlocks map[hashx.Hash]struct{}
}

// AddKnownInventory adds the passed inventory to the cache of known inventory for the peer.
func (p *peerSyncState) AddKnownInventory(invVect *protocol.InvInfo) {
	p.setInventoryKnown.Add(invVect)
}


func newPeerSyncState() *peerSyncState{
	return &peerSyncState{
		setInventoryKnown:newInventorySet(maxKnownInventory),
		requestedTxns:make(map[hashx.Hash]struct{}),
		requestedBlocks:make(map[hashx.Hash]struct{}),
	}
}
