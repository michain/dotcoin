package protocol

type Message struct{
	AddrSource string
	AddrFrom string
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
	CmdPing        = "ping"
	CmdPong        = "pong"
	CmdAlert       = "alert"
	CmdMemPool     = "mempool"
	CmdMerkleBlock = "merkleblock"
	CmdSendHeaders = "sendheaders"
)