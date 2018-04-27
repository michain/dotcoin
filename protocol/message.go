package protocol

type netMessage struct{
	AddrSource string
	AddrFrom string
}
// Message is an interface that describes a message.
type Message interface {
	Command() string
}

const (
	CmdVersion     = "version"
	CmdGetAddr     = "getaddr"
	CmdAddr        = "addr"
	CmdGetBlocks   = "getblocks"
	CmdInv         = "inv"
	CmdGetData     = "getdata"
	CmdNotFound    = "notfound"
	CmdBlock       = "block"
	CmdTx          = "tx"
	CmdGetHeaders  = "getheaders"
	CmdHeaders     = "headers"
	CmdAlert       = "alert"
	CmdMemPool     = "mempool"
	CmdMerkleBlock = "merkleblock"
	CmdSendHeaders = "sendheaders"
)