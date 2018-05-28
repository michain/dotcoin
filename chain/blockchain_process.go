package chain

import (
	"github.com/michain/dotcoin/util/hashx"
	"github.com/michain/dotcoin/logx"
	"fmt"
	"errors"
	"github.com/labstack/gommon/log"
	"math/big"
)

func (bc *Blockchain) maybeAcceptBlock(block *Block) (bool, error){
	logx.Infof("Blockchain maybeAcceptBlock block %v", block.GetHash())
	//check block validate
	err := bc.ValidateBlock(block, big.NewInt(int64(block.Difficult)))
	if err != nil{
		return false, err
	}
	bc.AddBlock(block)
	return true, nil
}

func (bc *Blockchain) processOrphans(hash *hashx.Hash) error {
	// Start with processing at least the passed hash.
	processHashes := make([]*hashx.Hash, 0, 10)
	processHashes = append(processHashes, hash)
	for len(processHashes) > 0 {
		// Pop the first hash to process from the slice.
		processHash := processHashes[0]
		processHashes[0] = nil // Prevent GC leak.
		processHashes = processHashes[1:]

		// Look up all orphans that are parented by the block we just accepted.
		for i := 0; i < len(bc.prevOrphanBlocks[*processHash]); i++ {
			orphan := bc.prevOrphanBlocks[*processHash][i]
			if orphan == nil {
				logx.Warnf("Found a nil entry at index %d in the "+
					"orphan dependency list for block %v", i,
					processHash)
				continue
			}

			// Remove the orphan from the orphan pool.
			orphanHash := orphan.GetHash()
			bc.removeOrphanBlock(orphan)
			i--

			//Add to block chain
			_, err := bc.maybeAcceptBlock(orphan)
			if err != nil{
				return err
			}

			// Add this block to the list of blocks to process so
			// any orphan blocks that depend on this block are
			// handled too.
			processHashes = append(processHashes, orphanHash)
		}
	}
	return nil
}


// removeOrphanBlock removes the orphan block from the orphan pool and previous orphan index.
func (bc *Blockchain) removeOrphanBlock(orphan *Block) {
	// Protect concurrent access.
	bc.orphanLock.Lock()
	defer bc.orphanLock.Unlock()

	// Remove the orphan block.
	orphanHash := orphan.GetHash()
	delete(bc.orphanBlocks, *orphanHash)

	// Remove the reference from the previous orphan index too
	prevHash := orphan.GetPrevHash()
	prevOrphans := bc.prevOrphanBlocks[*prevHash]
	for i := 0; i < len(prevOrphans); i++ {
		hash := prevOrphans[i].GetHash()
		if hash.IsEqual(orphanHash) {
			copy(prevOrphans[i:], prevOrphans[i+1:])
			prevOrphans[len(prevOrphans)-1] = nil
			prevOrphans = prevOrphans[:len(prevOrphans)-1]
			i--
		}
	}
	bc.prevOrphanBlocks[*prevHash] = prevOrphans

	// Remove the map elem if it's empty
	if len(bc.prevOrphanBlocks[*prevHash]) == 0 {
		delete(bc.prevOrphanBlocks, *prevHash)
	}
}


// ProcessBlock handling new block into chain
// return value: IsMainChain, IsOrphanBlock, error
func (bc *Blockchain) ProcessBlock(block *Block)(bool, bool, error){
	bc.chainLock.Lock()
	defer bc.chainLock.Unlock()

	blockHash := block.GetHash()
	logx.Tracef("Blockchain Processing block %v", blockHash)

	// The block must not already exist in the chain.
	exists, err := bc.HaveBlock(blockHash)
	if err != nil {
		logx.Errorf("Block Processing check have block error %v %v", err, blockHash)
		return false, false, err
	}
	if exists {
		str := fmt.Sprintf("Block Processing already have block %v", blockHash)
		return false, false, errors.New(str)
	}

	// The block must not already exist as an orphan.
	if _, exists := bc.orphanBlocks[*blockHash]; exists {
		str := fmt.Sprintf("Block Processing already have block (orphan) %v", blockHash)
		return false, false, errors.New(str)
	}

	//TODO checkBlockSanity

	//check prevHash, if not exists, add to orphanBlocks\
	if block.Height != 0 {
		prevHash := block.GetPrevHash()
		prevHashExists, err := bc.HaveBlock(prevHash)
		if err != nil {
			logx.Errorf("Block Processing check have block for prevhash error %v %v", err, prevHash)
			return false, false, err
		}
		if !prevHashExists {
			logx.Infof("Block Processing Adding orphan block %v with parent %v", blockHash, prevHash)
			bc.addOrphanBlock(block)
			return false, true, nil
		}
	}

	//add block into chain
	isMainChain, err := bc.maybeAcceptBlock(block)
	if err != nil {
		return false, false, err
	}

	// check and accept orphans which parented by this block
	// the accept blocks will remove from orphan pool
	err = bc.processOrphans(blockHash)
	if err != nil {
		return false, false, err
	}

	log.Debugf("Accepted block %v", blockHash)
	return isMainChain, false, nil
}