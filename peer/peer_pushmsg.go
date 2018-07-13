package peer

import (
	"github.com/michain/dotcoin/protocol"
	"math/rand"
	"github.com/michain/dotcoin/logx"
)

func (p *Peer) PushAddrMsg(msg *protocol.MsgAddr) error {
	addressCount := len(msg.AddrList)


	// Nothing to send.
	if addressCount == 0 {
		return nil
	}

	addresses := make([]string, addressCount)
	copy(addresses, msg.AddrList)

	// Randomize the addresses sent if there are more than the maximum allowed.
	if addressCount > protocol.MaxAddrPerMsg {
		// Shuffle the address list.
		for i := 0; i < protocol.MaxAddrPerMsg; i++ {
			j := i + rand.Intn(addressCount-i)
			addresses[i], addresses[j] = addresses[j], addresses[i]
		}

		// Truncate it to the maximum size.
		addresses = addresses[:protocol.MaxAddrPerMsg]
	}
	msg.AddrList = addresses


	return nil
}

func (p *Peer) PushVersion(msg *protocol.MsgVersion) error{
	//logx.DevPrintf("Peer.PushVersion %v", msg)
	if msg.AddrFrom == ""{
		msg.AddrFrom = p.GetSeedAddr()
	}
	p.SendSingleMessage(msg)
	return nil
}

func (p *Peer) PushGetBlocks(msg *protocol.MsgGetBlocks) error{
	//logx.DevPrintf("Peer.PushGetBlocks peer:%v msg:%v", p.GetListenAddr(), msg)
	p.SendSingleMessage(msg)
	return nil
}

func (p *Peer) PushBlock(msg *protocol.MsgBlock) error{
	logx.DevPrintf("Peer.PushBlock peer:%v remote:%v msg:%v trans:%d", p.GetListenAddr(), msg.GetFromAddr(), msg.Block.GetHash(), len(msg.Block.Transactions))
	p.SendSingleMessage(msg)
	return nil
}

func (p *Peer) PushTx(msg *protocol.MsgTx) error{
	logx.DevPrintf("Peer.PushTx peer:%v remote:%v msg:%v", p.GetListenAddr(), msg.GetFromAddr(), msg.Tx.GetHash())
	p.SendSingleMessage(msg)
	return nil
}