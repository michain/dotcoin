package server

import (
	"github.com/michain/dotcoin/chain"
	"errors"
	"fmt"
	"github.com/michain/dotcoin/server/packet"
	"github.com/michain/dotcoin/wallet"
	"github.com/michain/dotcoin/protocol"
)

var (
	ErrorPacketDeserialize = errors.New("packet data deserialize err")
)


type RpcHandler struct{
	server *Server
}

// SendTX send TX on rpc server
func (h *RpcHandler) SendTX(txPacket packet.TXPacket, result *packet.JsonResult) error {
	if !wallet.ValidateAddress(txPacket.From){
		return errors.New("validate [from] address error")
	}
	if !wallet.ValidateAddress(txPacket.To){
		return errors.New("validate [to] address error")
	}
	fromWallet := h.server.Wallets.GetWallet(txPacket.From)
	if fromWallet  == nil{
		return errors.New("not exists [from] address")
	}
	tx, err := chain.NewUTXOTransaction(fromWallet, txPacket.To, txPacket.Money, h.server.BlockChain.GetUTXOSet(), h.server.TXMemPool)
	if err != nil{
		return  err
	}


	//send inv message
	inv := protocol.NewInvInfo(protocol.InvTypeTx, *tx.GetHash())
	msgInv := protocol.NewMsgInv()
	msgInv.AddrFrom = h.server.ListenAddress
	msgInv.AddInvInfo(inv)
	h.server.Peer.BroadcastMessage(msgInv)

	*result = packet.JsonResult{0, "ok", tx.StringID()}
	return nil
}

// CreateWallet 创建账户
func (h *RpcHandler) CreateWallet(name string, result *packet.JsonResult) error {
	newWallet := h.server.Wallets.CreateWallet()

	*result = packet.JsonResult{RetCode:0, RetMsg:"ok", Message:newWallet.GetStringAddress()}
	fmt.Println(result)
	return nil
}

func (h *RpcHandler) ListAddress(tag string, result *packet.JsonResult) error {
	*result = packet.JsonResult{RetCode:0, RetMsg:"ok", Message:packet.WalletListPacket{h.server.Wallets.GetAddresses()}}
	fmt.Println(result)
	return nil
}

// ListMemPool list tx in mempool
func (h *RpcHandler) ListMemPool(name string, result *packet.JsonResult) error {
	txs := h.server.TXMemPool.TxDescs()

	*result = packet.JsonResult{RetCode:0, RetMsg:"ok", Message:txs}
	fmt.Println(result)
	return nil
}

// ListBlocks list blocks
func (h *RpcHandler) ListBlocks(name string, result *packet.JsonResult) error {
	h.server.BlockChain.ListBlockHashs()
	last, err := h.server.BlockChain.GetLastBlock()
	if err != nil{
		return err
	}

	*result = packet.JsonResult{RetCode:0, RetMsg:"ok", Message:last}
	fmt.Println(result)
	return nil
}
