package main

import (
	"fmt"
	"github.com/michain/dotcoin/server"
)

const tcpPort = ":2398"

func main(){
	//nodeID := blockchain.CreateNodeID()
	nodeID := "3eb456d086f34118925793496cd20945"
	fmt.Println("[Important!!!] this node ID:", nodeID)
	server.StartServer(nodeID, "", tcpPort, tcpPort, true)

}
