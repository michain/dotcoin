package cli

import (
	"fmt"
	"github.com/michain/dotcoin/wallet"
	"log"
)

func (cli *CLI) listAddresses(nodeID string) {
	wallets, err := wallet.LoadWallets(nodeID)
	if err!=nil{
		log.Panic("ERROR: Load wallets failed", err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

