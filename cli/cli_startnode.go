package cli

import (
	"fmt"
	"github.com/michain/dotcoin/wallet"
	"log"
	"github.com/michain/dotcoin/server"
)

func (cli *CLI) startNode(nodeID, minerAddress string, isGenesis bool) {
	fmt.Printf("Starting node %s\n", nodeID)
	nodeID = "3eb456d086f34118925793496cd20945"
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	server.StartServer(nodeID, minerAddress, tcpPort, tcpPort, isGenesis)
}

