package mempool

import (
	"sync"
	"log"
	"time"
	"github.com/michain/dotcoin/chain"
	"sync/atomic"
	"fmt"
	"github.com/pkg/errors"
	"container/list"
	"encoding/hex"
	"bytes"
	"encoding/gob"

	"github.com/michain/dotcoin/logx"
	"github.com/michain/dotcoin/config/chainhash"
)


const (
	// DefaultBlockPrioritySize is the default size in bytes for high-
	// priority / low-fee transactions.  It is used to help determine which
	// are allowed into the mempool and consequently affects their relay and
	// inclusion when generating block templates.
	DefaultBlockPrioritySize = 50000

	// orphanTTL is the maximum amount of time an orphan is allowed to
	// stay in the orphan pool before it expires and is evicted during the
	// next scan.
	orphanTTL = time.Minute * 15

	// orphanExpireScanInterval is the minimum amount of time in between
	// scans of the orphan pool to evict expired transactions.
	orphanExpireScanInterval = time.Minute * 5

	//now it's free for new transaction!!
	defaultFee = 0

	defaultMaxOrphanTransactions = 100
	defaultMaxOrphanTxSize       = 100000
)

type orphanTx struct {
	tx         *chain.Transaction
	expiration time.Time
}


type TxPool struct {
	// The following variables must only be used atomically.
	lastUpdated int64 // last time pool was updated

	mtx           sync.RWMutex
	pool          map[string]*chain.Transaction
	orphans       map[string]*orphanTx

	bc 			  *chain.Blockchain

	orphansByPrev map[chain.OutPoint]map[string]*chain.Transaction
	outpoints     map[chain.OutPoint]*chain.Transaction

	// nextExpireScan is the time after which the orphan pool will be
	// scanned in order to evict orphans.  This is NOT a hard deadline as
	// the scan will only run when an orphan is added to the pool as opposed
	// to on an unconditional timer.
	nextExpireScan time.Time
}

// New returns a new memory pool for validating and storing standalone
// transactions until they are mined into a block.
func New(bc *chain.Blockchain) *TxPool {
	return &TxPool{
		bc:				bc,
		pool:           make(map[string]*chain.Transaction),
		orphans:        make(map[string]*orphanTx),
		nextExpireScan: time.Now().Add(orphanExpireScanInterval),
		outpoints:		make(map[chain.OutPoint]*chain.Transaction),
	}
}


// TxHashes returns a slice of hashes for all of the transactions in the memory
// pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) TxHashes() [][]byte {
	mp.mtx.RLock()
	hashes := make([][]byte, len(mp.pool))
	i := 0
	for hash := range mp.pool {
		hashCopy := hash
		hb, err := hex.DecodeString(hashCopy)
		if err != nil{
			continue
		}
		hashes = append(hashes, hb)
		i++
	}
	mp.mtx.RUnlock()

	return hashes
}

// TxDescs returns a slice of descriptors for all the transactions in the pool.
// The descriptors are to be treated as read only.
//
// This function is safe for concurrent access.
func (mp *TxPool) TxDescs() []*chain.Transaction {
	mp.mtx.RLock()
	descs := make([]*chain.Transaction, len(mp.pool))
	i := 0
	for _, desc := range mp.pool {
		descs[i] = desc
		i++
	}
	mp.mtx.RUnlock()

	return descs
}


// addOrphan adds an orphan transaction to the orphan pool.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) addOrphan(tx *chain.Transaction) {

	mp.orphans[tx.StringHash()] = &orphanTx{
		tx:         tx,
		expiration: time.Now().Add(orphanTTL),
	}
	logx.Debugf("Stored orphan transaction %v (total: %d)", tx.Hash(),
		len(mp.orphans))
}


// removeOrphan is the internal function which implements the public
// RemoveOrphan.  See the comment for RemoveOrphan for more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) removeOrphan(tx *chain.Transaction, removeRedeemers bool) {
	// Nothing to do if passed tx is not an orphan.
	txHash := tx.StringHash()
	otx, exists := mp.orphans[txHash]
	if !exists {
		return
	}

	// Remove the reference from the previous orphan index.
	for _, txIn := range otx.tx.Inputs {
		orphans, exists := mp.orphansByPrev[txIn.PreviousOutPoint]
		if exists {
			delete(orphans, txHash)

			// Remove the map entry altogether if there are no
			// longer any orphans which depend on it.
			if len(orphans) == 0 {
				delete(mp.orphansByPrev, txIn.PreviousOutPoint)
			}
		}
	}


	// Remove the transaction from the orphan pool.
	delete(mp.orphans, txHash)
}

// RemoveOrphan removes the passed orphan transaction from the orphan pool and
// previous orphan index.
//
// This function is safe for concurrent access.
func (mp *TxPool) RemoveOrphan(tx *chain.Transaction) {
	mp.mtx.Lock()
	mp.removeOrphan(tx, false)
	mp.mtx.Unlock()
}

// maybeAddOrphan potentially adds an orphan to the orphan pool.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) maybeAddOrphan(tx *chain.Transaction) error {
	// Ignore orphan transactions that are too large.
	//
	// this equates to a maximum memory used of
	// mp.cfg.Policy.MaxOrphanTxSize * mp.cfg.Policy.MaxOrphanTxs (which is ~5MB
	// using the default values at the time this comment was written).
	serializedLen := len(chain.SerializeTransaction(tx))
	if serializedLen > defaultMaxOrphanTxSize {
		str := fmt.Sprintf("orphan transaction size of %d bytes is "+
			"larger than max allowed size of %d bytes",
			serializedLen, defaultMaxOrphanTxSize)
		return errors.New(str)
	}

	// Add the orphan if the none of the above disqualified it.
	mp.addOrphan(tx)

	return nil
}


// isTransactionInPool returns whether or not the passed transaction already
// exists in the main pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) isTransactionInPool(hash string) bool {
	if _, exists := mp.pool[hash]; exists {
		return true
	}
	return false
}

// IsTransactionInPool returns whether or not the passed transaction already
// exists in the main pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) IsTransactionInPool(hash string) bool {
	// Protect concurrent access.
	mp.mtx.RLock()
	inPool := mp.isTransactionInPool(hash)
	mp.mtx.RUnlock()

	return inPool
}
// isOrphanInPool returns whether or not the passed transaction already exists
// in the orphan pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) isOrphanInPool(hash string) bool {
	if _, exists := mp.orphans[hash]; exists {
		return true
	}

	return false
}

// IsOrphanInPool returns whether or not the passed transaction already exists
// in the orphan pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) IsOrphanInPool(hash string) bool {
	// Protect concurrent access.
	mp.mtx.RLock()
	inPool := mp.isOrphanInPool(hash)
	mp.mtx.RUnlock()

	return inPool
}

// haveTransaction returns whether or not the passed transaction already exists
// in the main pool or in the orphan pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) haveTransaction(hash string) bool {
	return mp.isTransactionInPool(hash) || mp.isOrphanInPool(hash)
}

// HaveTransaction returns whether or not the passed transaction already exists
// in the main pool or in the orphan pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) HaveTransaction(hash []byte) bool {
	// Protect concurrent access.
	mp.mtx.RLock()
	haveTx := mp.haveTransaction(hex.EncodeToString(hash))
	mp.mtx.RUnlock()

	return haveTx
}

// removeTransaction is the internal function which implements the public
// RemoveTransaction.  See the comment for RemoveTransaction for more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) removeTransaction(tx *chain.Transaction, removeRedeemers bool) {
	txHash := tx.StringHash()

	// Remove the transaction if needed.
	_, exists := mp.pool[txHash]
	if exists {
		delete(mp.pool, txHash)
		atomic.StoreInt64(&mp.lastUpdated, time.Now().Unix())
	}
	logx.Debugf("Remove transaction %v (exists: %v)", tx.Hash(), exists)
}

// RemoveTransaction removes the passed transaction from the mempool.
//
// This function is safe for concurrent access.
func (mp *TxPool) RemoveTransaction(tx *chain.Transaction, removeRedeemers bool) {
	// Protect concurrent access.
	mp.mtx.Lock()
	mp.removeTransaction(tx, removeRedeemers)
	mp.mtx.Unlock()
}


// addTransaction adds the passed transaction to the memory pool.  It should
// not be called directly as it doesn't perform any validation.  This is a
// helper for maybeAcceptTransaction.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) addTransaction(tx *chain.Transaction, height int32, fee int64) error {

	//FeePerKB: fee * 1000 / int64(tx.MsgTx().SerializeSize()),
	//StartingPriority: mining.CalcPriority(tx.MsgTx(), utxoView, height),
	mp.pool[tx.StringHash()] = tx
	for _, txIn := range tx.Inputs {
		mp.outpoints[txIn.PreviousOutPoint] = tx
	}

	atomic.StoreInt64(&mp.lastUpdated, time.Now().Unix())

	return nil
}


// maybeAcceptTransaction is the internal function which implements the public
// MaybeAcceptTransaction.  See the comment for MaybeAcceptTransaction for
// more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) maybeAcceptTransaction(tx *chain.Transaction, isNew, rateLimit, rejectDupOrphans bool) ([]*chainhash.Hash, error) {
	txHash := tx.StringHash()

	// Don't accept the transaction if it already exists in the pool.  This
	// applies to orphan transactions as well when the reject duplicate
	// orphans flag is set.  This check is intended to be a quick check to
	// weed out duplicates.
	if mp.isTransactionInPool(txHash) || (rejectDupOrphans &&
		mp.isOrphanInPool(txHash)) {

		str := fmt.Sprintf("already have transaction %v", txHash)
		return nil, errors.New(str)
	}

	// A standalone transaction must not be a coinbase transaction.
	if tx.IsCoinBase() {
		str := fmt.Sprintf("transaction %v is an individual coinbase",
			txHash)
		return nil,  errors.New(str)
	}

	// Get the current height of the main chain.  A standalone transaction
	// will be mined into the next block at best, so its height is at least
	// one more than the current height.
	bestHeight := mp.bc.GetBestHeight()

	//check double spend in pool
	err := mp.checkDoubleSpendInPool(tx)
	if err != nil {
		return nil, err
	}


	// Verify crypto signatures for each input and reject the transaction if
	// any don't verify.
	isVerify := mp.bc.VerifyTransaction(tx)
	if !isVerify {
		return nil, errors.New("Verify transaction sign failed")
	}

	// Add to transaction pool.
	mp.addTransaction(tx, int32(bestHeight), defaultFee)
	//TODO: log addTransaction err info

	logx.Debugf("Accepted transaction %v (pool size: %v)", txHash, len(mp.pool))

	return nil, nil
}


// MaybeAcceptTransaction is the main workhorse for handling insertion of new
// free-standing transactions into a memory pool.

// If the transaction is an orphan (missing parent transactions), the
// transaction is NOT added to the orphan pool, but each unknown referenced
// parent is returned.  Use ProcessTransaction instead if new orphans should
// be added to the orphan pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) MaybeAcceptTransaction(tx *chain.Transaction, isNew, rateLimit bool) ([]*chainhash.Hash, error) {
	// Protect concurrent access.
	mp.mtx.Lock()
	hashes, err := mp.maybeAcceptTransaction(tx, isNew, rateLimit, true)
	mp.mtx.Unlock()

	return hashes, err
}

// processOrphans is the internal function which implements the public
// ProcessOrphans.  See the comment for ProcessOrphans for more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) processOrphans(acceptedTx *chain.Transaction) error {

	// Start with processing at least the passed transaction.
	processList := list.New()
	processList.PushBack(acceptedTx)
	for processList.Len() > 0 {
		// Pop the transaction to process from the front of the list.
		firstElement := processList.Remove(processList.Front())
		processItem := firstElement.(*chain.Transaction)

		prevOut := chain.OutPoint{Hash: processItem.Hash()}
		for txOutIdx  := range processItem.Outputs {
			// Look up all orphans that redeem the output that is
			// now available.  This will typically only be one, but
			// it could be multiple if the orphan pool contains
			// double spends.  While it may seem odd that the orphan
			// pool would allow this since there can only possibly
			// ultimately be a single redeemer, it's important to
			// track it this way to prevent malicious actors from
			// being able to purposely constructing orphans that
			// would otherwise make outputs unspendable.
			//
			// Skip to the next available output if there are none.
			prevOut.Index = txOutIdx
			orphans, exists := mp.orphansByPrev[prevOut]
			if !exists {
				continue
			}

			// Potentially accept an orphan into the tx pool.
			for _, tx := range orphans {
				missing, err := mp.maybeAcceptTransaction(tx, true, true, false)
				if err != nil {
					// The orphan is now invalid, so there
					// is no way any other orphans which
					// redeem any of its outputs can be
					// accepted.  Remove them.
					mp.removeOrphan(tx, true)
					break
				}

				// Transaction is still an orphan.  Try the next
				// orphan which redeems this output.
				if len(missing) > 0 {
					continue
				}

				mp.removeOrphan(tx, false)
				processList.PushBack(tx)

				// Only one transaction for this outpoint can be
				// accepted, so the rest are now double spends
				// and are removed later.
				break
			}
		}
	}

	return nil
}

// ProcessTransaction is the main workhorse for handling insertion of new
// free-standing transactions into the memory pool.  It includes functionality
// such as rejecting duplicate transactions, ensuring transactions follow all
// rules, orphan transaction handling, and insertion into the memory pool.
//
// It returns a slice of transactions added to the mempool.  When the
// error is nil, the list will include the passed transaction itself along
// with any additional orphan transaactions that were added as a result of
// the passed one being accepted.
//
// This function is safe for concurrent access.
func (mp *TxPool) ProcessTransaction(tx *chain.Transaction, allowOrphan, rateLimit bool) error {
	logx.Tracef("Processing transaction %v", tx.Hash())

	// Protect concurrent access.
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	// Potentially accept the transaction to the memory pool.
	missingParents, err := mp.maybeAcceptTransaction(tx, true, rateLimit, true)
	if err != nil {
		return err
	}

	if len(missingParents) == 0 {
		// Accept any orphan transactions that depend on this
		// transaction (they may no longer be orphans if all inputs
		// are now available) and repeat for those accepted
		// transactions until there are no more.
		return mp.processOrphans(tx)
	}

	// The transaction is an orphan (has inputs missing).  Reject
	// it if the flag to allow orphans is not set.
	if !allowOrphan {
		// Only use the first missing parent transaction in
		// the error message.
		//
		// NOTE: RejectDuplicate is really not an accurate
		// reject code here, but it matches the reference
		// implementation and there isn't a better choice due
		// to the limited number of reject codes.  Missing
		// inputs is assumed to mean they are already spent
		// which is not really always the case.
		str := fmt.Sprintf("orphan transaction %v references "+
			"outputs of unknown or fully-spent "+
			"transaction %v", tx.Hash(), missingParents[0])
		return errors.New(str)
	}

	// Potentially add the orphan transaction to the orphan pool.
	err = mp.maybeAddOrphan(tx)
	return err
}


// checkDoubleSpendInPool checks whether or not the passed transaction is
// attempting to spend coins already spent by other transactions in the pool.
// Note it does not check for double spends against transactions already in the
// main chain.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) checkDoubleSpendInPool(tx *chain.Transaction) error {
	for _, txIn := range tx.Inputs{
		if txR, exists := mp.outpoints[txIn.PreviousOutPoint]; exists {
			str := fmt.Sprintf("output %v already spent by "+
				"transaction %v in the memory pool",
				txIn.PreviousOutPoint, txR.Hash())
			return errors.New("DoubleSpend " +str)
		}
	}

	return nil
}

// Count returns the number of transactions in the main pool.  It does not
// include the orphan pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) Count() int {
	mp.mtx.RLock()
	count := len(mp.pool)
	mp.mtx.RUnlock()

	return count
}


// SerializeTxPool serializes a txpool for []byte
func SerializeTxPool(pool *TxPool) []byte{
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(pool)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// DeserializeTxPool deserializes a txpool
func DeserializeTxPool(data []byte) (*TxPool, error) {
	var pool TxPool

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&pool)
	return &pool, err
}

