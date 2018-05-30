package chain

import "github.com/pkg/errors"

var ErrorNotFoundTransaction = errors.New("not found the transaction")

var ErrBlockNoTransactions = errors.New("block does not contain any transactions")

var ErrBlockTooManyTransactions = errors.New("block has too many transactions")

var ErrFirstTxNotCoinbase = errors.New("first transaction in block is not a coinbase")

var ErrNotVerifyTransaction = errors.New("block has not verfiy transaction")

var ErrMultipleCoinbases = errors.New("block contains second coinbase transaction")

var ErrBlockBadMerkleRoot = errors.New("block merkle root is invalid")

var ErrBlockDuplicateTx = errors.New("block contains duplicate transaction")
