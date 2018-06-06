package cli

import (
	"github.com/michain/dotcoin/chain"
	"log"
)

func (cli *CLI) printChain(nodeID string) {
	bc, err := chain.LoadBlockChain(nodeID)
	if err != nil{
		log.Panic("ERROR: Load blockchain failed", err)
	}
	bc.ListBlockHashs()
}

