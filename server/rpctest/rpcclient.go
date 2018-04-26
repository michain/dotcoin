package rpctest

import (
	"net/rpc/jsonrpc"
	"fmt"
	"os"
	"github.com/michain/dotcoin/server/packet"
	"net/rpc"
)

func CallCreateWallet(){
	client := connRpcServer()

	var reply packet.JsonResult
	err := client.Call("Rpc.CreateWallet", "test", &reply)
	fmt.Println(reply.Message, err)
}

func CallListAddress(){

	client := connRpcServer()
	var reply packet.JsonResult
	err := client.Call("Rpc.ListAddress", "test", &reply)
	fmt.Println(reply.Message)
	addresses :=  reply.Message
	fmt.Println(addresses, err)
}


func CallSendTX(from, to string) error{
	client := connRpcServer()
	var reply packet.JsonResult
	packet := packet.TXPacket{From:from, To:to, Money:1}
	err := client.Call("Rpc.SendTX", packet, &reply)
	fmt.Println(err, packet, reply)
	return err
}


func connRpcServer() *rpc.Client{
	service := "127.0.0.1:2398"
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		fmt.Println("dial error:", err)
		os.Exit(1)
	}
	return client
}