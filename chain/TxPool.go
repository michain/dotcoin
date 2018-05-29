package chain

type TxPool interface{
	HaveTransaction(hash string) bool
}
