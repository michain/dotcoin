package chain

import (
	"sync"
	"math/big"
	"github.com/michain/dotcoin/util/hashx"
	"github.com/boltdb/bolt"
	"encoding/binary"
	"github.com/michain/dotcoin/storage"
	"bytes"
	"encoding/gob"
	"log"
)
// approxNodesPerWeek is an approximation of the number of new blocks there are
// in a week on average.
const approxNodesPerWeek = 6 * 24 * 7


type ChainIndex struct {
	db  *bolt.DB
	sync.RWMutex
	indexes []*blockIndex
	hashIndexes map[hashx.Hash]*blockIndex
	dirtyIndexes map[*blockIndex]struct{}
}

// newChainIndex returns a new chain index for the given tip block index.
// The tip can be updated at any time via the setTip function.
func newChainIndex(db *bolt.DB, index *blockIndex) *ChainIndex {
	bi :=&ChainIndex{
		db :			db,
		indexes:	[]*blockIndex{},
		hashIndexes:	make(map[hashx.Hash]*blockIndex),
		dirtyIndexes:   make(map[*blockIndex]struct{}),
	}
	bi.SetTip(index)
	return bi
}


// addIndex adds the provided index to the chain index
// This can be used while initializing the block index.
//
// This function is NOT safe for concurrent access.
func (bi *ChainIndex) addIndex(index *blockIndex) {
	bi.hashIndexes[index.Hash] = index
}


// tip returns the current tip block index for the chain index.
// It will return nil if there is no tip.
//
// This function MUST be called with the view mutex locked (for reads).
func (bi *ChainIndex) tip() *blockIndex {
	if len(bi.indexes) == 0 {
		return nil
	}

	return bi.indexes[len(bi.indexes)-1]
}

// setTip sets use the provided block index as the current tip
//
// This function MUST be called with the view mutex locked (for writes).
func (bi *ChainIndex) setTip(index *blockIndex) {
	if index == nil {
		// Keep the backing array around for potential future use.
		bi.indexes = bi.indexes[:0]
		return
	}

	// Create or resize the slice that will hold the block indexes to the
	// provided tip height.
	needed := index.Height + 1
	if int32(cap(bi.indexes)) < needed {
		nodes := make([]*blockIndex, needed, needed+approxNodesPerWeek)
		copy(nodes, bi.indexes)
		bi.indexes = nodes
	} else {
		prevLen := int32(len(bi.indexes))
		bi.indexes = bi.indexes[0:needed]
		for i := prevLen; i < needed; i++ {
			bi.indexes[i] = nil
		}
	}

	for index != nil && bi.indexes[index.Height] != index {
		bi.indexes[index.Height] = index
		index = index.Parent
	}
}

// height returns the height of the tip of the chain index.
// It will return -1 if there is no tip
//
// This function MUST be called with the view mutex locked (for reads).
func (bi *ChainIndex) height() int32 {
	return int32(len(bi.indexes) - 1)
}

// nodeByHeight returns the block index at the specified height.
// Nil will be returned if the height does not exist.
//
// This function MUST be called with the view mutex locked (for reads).
func (bi *ChainIndex) nodeByHeight(height int32) *blockIndex {
	if height < 0 || height >= int32(len(bi.indexes)) {
		return nil
	}

	return bi.indexes[height]
}

// Tip returns the current tip block index for the chain index.
// It will return nil if there is no tip.
//
// This function is safe for concurrent access.
func (bi *ChainIndex) Tip() *blockIndex {
	bi.Lock()
	tip := bi.tip()
	bi.Unlock()
	return tip
}

// SetTip sets use the provided block index as the current tip
//
// This function is safe for concurrent access.
func (bi *ChainIndex) SetTip(node *blockIndex) {
	bi.Lock()
	bi.setTip(node)
	bi.Unlock()
}

// height returns the height of the tip of the chain index.
// It will return -1 if there is no tip
//
// This function is safe for concurrent access.
func (bi *ChainIndex) Height() int32 {
	bi.Lock()
	height := bi.height()
	bi.Unlock()
	return height
}

// IndexByHeight returns the block index at the specified height.
// Nil will be returned if the height does not exist.
//
// This function is safe for concurrent access.
func (bi *ChainIndex) IndexByHeight(height int32) *blockIndex {
	bi.Lock()
	node := bi.nodeByHeight(height)
	bi.Unlock()
	return node
}

// HaveBlock returns whether or not the block index contains the provided hash.
//
// This function is safe for concurrent access.
func (bi *ChainIndex) HaveBlock(hash *hashx.Hash) bool {
	bi.RLock()
	_, hasBlock := bi.hashIndexes[*hash]
	bi.RUnlock()
	return hasBlock
}

// LookupNode returns the block node identified by the provided hash.
// It will return nil if there is no entry for the hash.
//
// This function is safe for concurrent access.
func (bi *ChainIndex) LookupNode(hash *hashx.Hash) *blockIndex {
	bi.RLock()
	node := bi.hashIndexes[*hash]
	bi.RUnlock()
	return node
}

// AddIndex adds the provided index to the chain index
//
// This function is safe for concurrent access.
func (bi *ChainIndex) AddIndex(index *blockIndex) {
	bi.Lock()
	bi.addIndex(index)
	bi.dirtyIndexes[index] = struct{}{}
	bi.Unlock()
}

// flushToDB writes all dirty block nodes to the database. If all writes
// succeed, this clears the dirty set.
func (bi *ChainIndex) flushToDB() error {
	bi.Lock()
	if len(bi.dirtyIndexes) == 0 {
		bi.Unlock()
		return nil
	}
	var err error

	//TODO need transaction all save
	for index := range bi.dirtyIndexes {
		key := blockIndexKey(&index.Hash, uint32(index.Height))
		storage.SaveBlockIndex(bi.db, key, SerializeBlockIndex(index))
	}

	// If write was successful, clear the dirty set.
	if err == nil {
		bi.dirtyIndexes = make(map[*blockIndex]struct{})
	}

	bi.Unlock()
	return err
}

func blockIndexKey(blockHash *hashx.Hash, blockHeight uint32) []byte {
	indexKey := make([]byte, hashx.HashSize+4)
	binary.BigEndian.PutUint32(indexKey[0:4], blockHeight)
	copy(indexKey[4:hashx.HashSize+4], blockHash[:])
	return indexKey
}



type blockIndex struct{
	// parent is the parent block for this index.
	Parent *blockIndex
	// hash is the double sha 256 of the block.
	Hash hashx.Hash

	// workSum is the total amount of work in the chain up to and including
	// this node.
	WorkSum *big.Int
	// height is the position in the block chain.
	Height int32

	Difficult       uint32
	Nonce      int64
	Timestamp  int64
	MerkleRoot []byte
}

// newBlockIndex returns a new block index for the given block and parent index
// calculating the height and workSum
//
// This function is NOT safe for concurrent access.
func newBlockIndex(block *Block, parent *blockIndex) *blockIndex {
	index := blockIndex{
		Hash:       *block.GetHash(),
		WorkSum:    CalcWork(block.Difficult),
		Difficult:  block.Difficult,
		Nonce:      block.Nonce,
		Timestamp:  block.Timestamp.Unix(),
		MerkleRoot: block.MerkleRoot,
	}
	if parent != nil {
		index.Parent = parent
		index.Height = parent.Height + 1
		index.WorkSum = index.WorkSum.Add(parent.WorkSum, index.WorkSum)
	}
	return &index
}
// SerializeBlock serializes the block
func SerializeBlockIndex(b *blockIndex) []byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlockIndex deserializes a block index
func DeserializeBlockIndex(d []byte) *blockIndex {
	var index blockIndex

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&index)
	if err != nil {
		log.Panic(err)
	}

	return &index
}
