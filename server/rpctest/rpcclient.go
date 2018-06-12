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


func CallGetVersion(){

	client := connRpcServer()
	var reply packet.JsonResult
	err := client.Call("Rpc.GetVersion", "test", &reply)
	fmt.Println("GetVersion:", reply.Message)
	addresses :=  reply.Message
	fmt.Println(addresses, err)
}


func CallListBlocks(){

	client := connRpcServer()
	var reply packet.JsonResult
	err := client.Call("Rpc.ListBlocks", "test", &reply)
	fmt.Println(reply.Message)
	addresses :=  reply.Message
	fmt.Println(addresses, err)
}

func CallListMemPool(){
	client := connRpcServer()
	var reply packet.JsonResult
	err := client.Call("Rpc.ListMemPool", "", &reply)
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
	//service := "192.168.42.91:12398"
	service := ":12398"
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		fmt.Println("dial error:", err)
		os.Exit(1)
	}
	return client
}
