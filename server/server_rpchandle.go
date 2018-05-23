package server

import (
	"github.com/michain/dotcoin/chain"
	"errors"
	"fmt"
	"github.com/michain/dotcoin/server/packet"
	"github.com/michain/dotcoin/wallet"
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
	tx := chain.NewUTXOTransaction(fromWallet, txPacket.To, txPacket.Money, h.server.BlockChain.GetUTXOSet())

	//add TX to mempool
	_, err := h.server.TXMemPool.MaybeAcceptTransaction(tx, true, true)
	if err != nil{
		return  err
	}

	//if nodeAddress == knownNodes[0] {
	//	for _, node := range knownNodes {
	//		if node != nodeAddress && node != txPacket.AddFrom {
				//TODO:send inventory to other server
				//sendInv(node, "tx", [][]byte{tx.ID})
	//		}
	//	}
	//}

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