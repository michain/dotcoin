package chain

import (
	"github.com/boltdb/bolt"
	"fmt"
	"crypto/ecdsa"
	"log"
	"errors"
	"os"
	"github.com/michain/dotcoin/util"
	"github.com/michain/dotcoin/storage"
	"github.com/michain/dotcoin/wallet"
	"github.com/michain/dotcoin/util/hashx"
	"sync"
	"encoding/hex"
)

const genesisCoinbaseData = "The Times 15/April/2018 for my 35th birthday!"

const (
	defaultNonce = 0
	blockDefaultDifficult = 10
)



var ErrorBlockChainNotFount = errors.New("blockchain is not found")
var ErrorNoExistsAnyBlock = errors.New("not exists any block")

// Blockchain implements interactions with a DB
type Blockchain struct {
	lastBlockHash []byte
	db  *bolt.DB
	chainLock *sync.RWMutex

	orphanLock   *sync.RWMutex
	orphanBlocks map[hashx.Hash]*Block

	// previous hash index for faster lookups
	prevOrphanBlocks map[hashx.Hash][]*Block
}

// CreateBlockchain creates a new blockchain with genesisBlock
func CreateBlockchain(isGenesisNode bool, address, nodeID string) *Blockchain {
	dbFile := storage.GetDBFileName(nodeID)
	if util.ExitFile(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	fmt.Println("CreateBlockchain Begin")


	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic("Open db error", err)
	}

	//create bolt block bucket
	err = storage.CreateBlockBucket(db)
	if err != nil {
		log.Panic("CreateBlockBucket error", err)
	}

	var lastBlockHash []byte
	if isGenesisNode {
		genesis := NewGenesisBlock(address)

		err =storage.SaveBlock(db, genesis.Hash, SerializeBlock(genesis))
		if err != nil {
			log.Panic("SaveBlock error", err)
		}else{
			lastBlockHash = genesis.Hash
		}
	}



	bc := Blockchain{
		lastBlockHash:lastBlockHash,
		db:db,
		chainLock:new(sync.RWMutex),
		orphanLock:new(sync.RWMutex),
		orphanBlocks:make(map[hashx.Hash]*Block),
	}

	fmt.Println("CreateBlockchain Success!")
	fmt.Println(fmt.Sprintf("lastBlockHash %x", bc.lastBlockHash))

	if isGenesisNode {
		//Rebuild UTXO data
		bc.GetUTXOSet().Rebuild()
	}

	return &bc
}

// LoadBlockChain load Blockchain with nodeID from bolt
func LoadBlockChain(nodeID string) (*Blockchain, error) {
	dbFile := storage.GetDBFileName(nodeID)
	if !util.ExitFile(dbFile) {
		fmt.Println("No existing blockchain found. Create one first.")
		return nil, ErrorBlockChainNotFount
	}

	var lastBlockHash []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	lastBlockHash, _, err = storage.GetLastBlock(db)
	if err != nil {
		return nil, err
	}

	bc := Blockchain{
		lastBlockHash:lastBlockHash,
		db:db,
		chainLock:new(sync.RWMutex),
		orphanLock:new(sync.RWMutex),
		orphanBlocks:make(map[hashx.Hash]*Block),
	}

	return &bc, nil
}

// GetStorageDB get storage db
func (bc *Blockchain) GetStorageDB() *bolt.DB {
	return bc.db
}



// addOrphanBlock add block into orphan blocks
func (bc *Blockchain) addOrphanBlock(block *Block){
	bc.orphanLock.Lock()
	defer bc.orphanLock.Unlock()
	bc.orphanBlocks[*block.GetHash()] = block

	// Add to previous hash index for faster lookups.
	prevHash := block.GetPrevHash()
	bc.prevOrphanBlocks[*prevHash] = append(bc.prevOrphanBlocks[*prevHash], block)
}

// AddBlock add the block into the blockchain
// save to bolt, update LastBlockHash
func (bc *Blockchain) AddBlock(block *Block) {
	bc.chainLock.Lock()
	defer bc.chainLock.Unlock()

	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(storage.BoltBlocksBucket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := SerializeBlock(block)
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte(storage.BoltLastHashKey), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.lastBlockHash = block.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// HaveBlock check block hash exists
func (bc *Blockchain) HaveBlock(blockHash *hashx.Hash) (bool, error){
	b, err:=bc.GetBlock(blockHash.CloneBytes())
	if err != nil{
		if err == storage.ErrorBlockNotFount{
			return false, nil
		}
		return false, err
	}else{
		return b != nil, nil
	}
}

// GetBlock finds a block by its hash and returns it
func (bc *Blockchain) GetBlock(blockHash []byte) (*Block, error) {
	var block *Block
	blockData, err := storage.GetBlock(bc.db, blockHash)
	if err != nil{
		return nil, err
	}

	block = DeserializeBlock(blockData)
	return block, err
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash, lastBlockData []byte
	var lastHeight int32
	var err error

	for _, tx := range transactions {
		// TODO: ignore transaction if it's not valid
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	lastHash, lastBlockData, err = storage.GetLastBlock(bc.db)
	if err != nil{
		log.Panic(err)
	}
	lastBlock := DeserializeBlock(lastBlockData)
	lastHeight = lastBlock.Height

	//run pow and create block
	newBlock := NewBlock(transactions, lastHash, lastHeight+1)

	//save block to db
	err = storage.SaveBlock(bc.db, newBlock.Hash, SerializeBlock(newBlock))
	if err != nil {
		log.Panic(err)
	}
	bc.lastBlockHash = newBlock.Hash

	return newBlock
}

// GetUTXOSet get current bc's UTXOSet wrapper
func (bc *Blockchain) GetUTXOSet() *UTXOSet{
	return &UTXOSet{bc}
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() map[string][]TXOutput {
	UTXO := make(map[string][]TXOutput)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := tx.StringID()

		Outputs:
			for outIdx, out := range tx.Outputs {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs = append(outs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinBase() == false {
				for _, in := range tx.Inputs {
					inTxID := in.PreviousOutPoint.StringHash()
					spentTXOs[inTxID] = append(spentTXOs[inTxID], int(in.PreviousOutPoint.Index))
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}


	return UTXO
}

// ListBlockHashs list println block's Hash and PrevBlockHash
func (bc *Blockchain) ListBlockHashs(){
	bci := bc.Iterator()
	for {
		block := bci.Next()
		if len(block.PrevBlockHash) != 0 {
			fmt.Println("ListBlockHashs", "prevhash:", hex.EncodeToString(block.PrevBlockHash), "hash:", hex.EncodeToString(block.Hash))
		}else{
			break
		}
	}
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Inputs {
		prevTX, err := bc.FindTransaction(&vin.PreviousOutPoint.Hash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[prevTX.StringID()] = *prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID *hashx.Hash) (*Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if tx.ID.IsEqual(ID){
				return tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return nil, ErrorNotFound
}

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}

	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.Inputs {
		prevTX, err := bc.FindTransaction(&vin.PreviousOutPoint.Hash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[prevTX.StringID()] = *prevTX
	}

	return tx.Verify(prevTXs)
}

// GetBalance
func (bc *Blockchain) GetBalance(address string) int{
	balance := 0
	pubKeyHash := wallet.GetPubKeyHashFromAddress([]byte(address))
	UTXOs := bc.GetUTXOSet().FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}
	return balance
}

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetBestHeight() int32 {
	var lastBlock *Block
	_, lastBlockData, err := storage.GetLastBlock(bc.db)
	if err != nil{
		return 0
	}

	if lastBlockData == nil{
		return 0
	}

	lastBlock = DeserializeBlock(lastBlockData)

	return lastBlock.Height
}

// GetLastBlock returns the latest block
func (bc *Blockchain) GetLastBlock() (*Block, error){
	var lastBlock *Block
	_, lastBlockData, err := storage.GetLastBlock(bc.db)
	if err != nil{
		return nil, err
	}

	if lastBlockData == nil{
		return nil, ErrorNoExistsAnyBlock
	}

	lastBlock = DeserializeBlock(lastBlockData)

	return lastBlock, nil
}

// GetBlockHashes returns a list of hashes with beginHash and maxNum limit
func (bc *Blockchain) GetBlockHashes(beginHash *hashx.Hash, stopHash hashx.Hash, maxNum int) ([]*hashx.Hash, error) {
	var blocks []*hashx.Hash
	bci := bc.Iterator()
	err := bci.LocationHash(beginHash.CloneBytes())
	if err != nil{
		return nil, err
	}

	getCount := 0
	for {
		block := bci.Next()

		fmt.Println("GetBlockHashes", block.GetHash(), hex.EncodeToString(block.Hash), hex.EncodeToString(block.GetHash().CloneBytes()))

		h := block.GetHash()

		if stopHash.IsEqual(h){
			break
		}

		blocks = append(blocks, h)
		getCount += 1

		if len(block.PrevBlockHash) == 0 {
			break
		}

		if getCount >= maxNum{
			break
		}
	}

	return blocks, nil
}

// Iterator returns a BlockchainIterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.lastBlockHash, bc.db}

	return bci
}



