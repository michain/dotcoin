package chain

import (
	"math/big"
	"bytes"
	"github.com/michain/dotcoin/util/hashx"
)

// ValidateBlock validate block data
func (bc *Blockchain)ValidateBlock(block *Block, powLimit *big.Int) error {
	//TODO: check ProofOfWork
	//TODO: check block time

	//check transaction's count
	//must have at least one
	numTx := len(block.Transactions)
	if numTx == 0 {
		return ErrBlockNoTransactions
	}

	// check max block payload is bigger than limit.
	if numTx > MaxBlockBaseSize {
		return ErrBlockTooManyTransaction
	}

	//TODO check max block's byte size

	// The first transaction in a block must be a coinbase.
	transactions := block.Transactions
	if transactions[0].IsCoinBase() {
		return ErrFirstTxNotCoinbase
	}

	// check coinbase transaction count
	// count == 1
	for _, tx := range transactions[1:] {
		if tx.IsCoinBase() {
			return ErrMultipleCoinbases
		}
	}

	// TODO validate each transaction
	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx){
			return ErrNotVerifyTransaction
		}
	}

	// check merkleRoot
	merkleRoot := block.HashTransactions()
	if bytes.Compare(block.MerkleRoot, merkleRoot) != 0{
		return ErrBlockBadMerkleRoot
	}

	// Check for duplicate transactions.
	existingTxHashes := make(map[hashx.Hash]struct{})
	for _, tx := range transactions {
		hash := tx.Hash()
		if _, exists := existingTxHashes[hash]; exists {
			return ErrBlockDuplicateTx
		}
		existingTxHashes[hash] = struct{}{}
	}

	return nil
}

