package protocol

type NetMessage struct{
	ProtocolVersion    uint32
	AddrFrom string
}

func (m *NetMessage) GetFromAddr() string{
	return m.AddrFrom
}

func (m *NetMessage) SetFromAddr(addr string){
	m.AddrFrom = addr
}


// Message is an interface that describes a message.
type Message interface {
	Command() string
	GetFromAddr() string
	SetFromAddr(addr string)
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

const(
	// ProtocolVersion is the latest protocol version this package supports.
	ProtocolVersion uint32 = 10000
)