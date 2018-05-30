package chain

import (
	"github.com/michain/dotcoin/util/hashx"
)

type TxPool interface{
	HaveTransaction(hash string) bool
	MaybeAcceptTransaction(tx *Transaction) ([]*hashx.Hash, error)
}
