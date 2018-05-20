package protocol

type netMessage struct{
	AddrFrom string
	needBroadcast bool
}

func (m *netMessage) GetFromAddr() string{
	return m.AddrFrom
}

func (m *netMessage) NeedBroadcast() bool{
	return m.needBroadcast
}

func (m *netMessage) SetNeedBroadcast(flag bool){
	m.needBroadcast = flag
}

// Message is an interface that describes a message.
type Message interface {
	Command() string
	GetFromAddr() string
	NeedBroadcast() bool
	SetNeedBroadcast(flag bool)
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