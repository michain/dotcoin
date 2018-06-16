package storage

import (
	"github.com/boltdb/bolt"
	"fmt"
	"log"
	"errors"
)

const(
	BoltLastHashKey = "l"
	BoltFileFormat = "dotchain_%s.db"
	BoltBlocksBucket = "blocks"
	BoltUTXOBucket = "chainstate"
	BoltTXMemPool = "txmempool"
)

var ErrorBlockNotFount = errors.New("Block is not found")


// getDBFileName get dbfile's name with NodeID
// like dotchain_XXXXXX.db
func GetDBFileName(nodeID string) string{
	return fmt.Sprintf(BoltFileFormat, nodeID)
}

// RemoveBlock remove block
func RemoveBlock(db *bolt.DB, blockHash []byte) error{
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBlocksBucket))
		err := b.Delete(blockHash)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// CreateBlockBucket
func CreateBlockBucket(db *bolt.DB) error{
	err := db.Update(func(tx *bolt.Tx) error {
		_, errc := tx.CreateBucket([]byte(BoltBlocksBucket))
		return errc
	})
	return err
}

func SaveBlock(db *bolt.DB, blockHash, blockData []byte) error{
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBlocksBucket))
		err := b.Put(blockHash, blockData)
		if err != nil {
			return err
		}

		err = b.Put([]byte(BoltLastHashKey),blockHash)
		if err != nil {
			return err
		}
		return nil
	})
	//TODO:log db operate
	return err
}

// GetBlock query block data with block hash
// if not exists, return ErrorBlockNotFount
func GetBlock(db *bolt.DB, blockHash []byte)(blockData []byte, err error){
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBlocksBucket))

		blockData = b.Get(blockHash)

		if blockData == nil {
			return ErrorBlockNotFount
		}
		return nil
	})
	return
}


func GetLastBlock(db *bolt.DB) (lastHash, lastBlockData []byte, err error){
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBlocksBucket))
		lastHash = b.Get([]byte(BoltLastHashKey))
		lastBlockData = b.Get(lastHash)
		return nil
	})
	return
}


func GetLashBlockHash(db *bolt.DB)(lastHash []byte, err error){
	lastHash, _, err = GetLastBlock(db)
	return
}

func GetTXMemPool(db *bolt.DB) ([]byte, error){
	var txPool []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBlocksBucket))

		txPool = b.Get([]byte(BoltTXMemPool))
		return nil
	})
	return txPool, err
}

// SaveTXMemPool save mempool into db
func SaveTXMemPool(db *bolt.DB, txPool []byte) error{
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltBlocksBucket))
		err := b.Put([]byte(BoltTXMemPool), txPool)
		if err != nil {
			return err
		}
		return nil
	})
	//TODO:log db operate
	return err
}


func CountTransactions(db *bolt.DB) int {
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BoltUTXOBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}
