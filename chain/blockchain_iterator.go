package chain

import (
	"github.com/boltdb/bolt"
	"log"
	"github.com/michain/dotcoin/storage"
)

// BlockchainIterator implement a iterator for blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(storage.BoltBlocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

// LocationHash locate current hash
func (i *BlockchainIterator) LocationHash(locateHash []byte) error{
	_, err:=storage.GetBlock(i.db, locateHash)
	if err == nil{
		i.currentHash = locateHash
	}
	return err
}
