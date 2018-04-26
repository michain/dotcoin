package main

import (
	"fmt"
	"github.com/michain/dotcoin/server"
)

func main(){
	//nodeID := blockchain.CreateNodeID()
	nodeID := "3eb456d086f34118925793496cd20945"
	fmt.Println("[Important!!!] this node ID:", nodeID)
	server.StartServer(nodeID, "", true)

}
