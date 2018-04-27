package server

import (
	"fmt"
	"github.com/michain/dotcoin/wallet"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net"
	"log"
	"github.com/michain/dotcoin/util/uuid"
	"github.com/michain/dotcoin/chain"
	"github.com/michain/dotcoin/mempool"
	"github.com/pkg/errors"
)

var txPool *mempool.TxPool
var minerAddress string
var currentWallets *wallet.WalletSet
var currentBlockChain *chain.Blockchain
var currentAddrManager *AddrManager
var nodeAddress string
const(
	rpcPort = ":2398" //2398 = 1983+0415 my birthday!
	coinbaseReward = 10
	knowAddr = "localhost:3000"
)


func listenRPCServer(bc *chain.Blockchain) {
	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		return
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("Rpc", &RpcHandler{bc}); err != nil {
		return
	}

	fmt.Println("begin listen ", lis.Addr())

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
			continue
		}
		go srv.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

// initServer init server
func initServer(nodeID, minerAddr string, isGenesisNode bool) error{
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	var err error
	isFirstInit := false

	currentBlockChain, err = chain.LoadBlockChain(nodeID)
	if err == chain.ErrorBlockChainNotFount{
		isFirstInit = true
	}
	if !isGenesisNode && isFirstInit{
		//TODO:sync data from other node
		//if err, block myself
	}


	//load or create miner wallet
	currentWallets, err = wallet.LoadWallets(nodeID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//ignore config miner address when genesis and first init
	if isGenesisNode && isFirstInit{
		mw := currentWallets.CreateWallet()
		err = currentWallets.SetMinerAddress(mw.GetStringAddress())
		if err != nil{
			fmt.Println(err)
			return err
		}
		minerAddress = mw.GetStringAddress()
	}
	//if set minerAddr, validate it
	if minerAddr != ""{
		if !wallet.ValidateAddress(minerAddr) {
			msg := "Validate minerAddr error " + minerAddr
			fmt.Println(msg)
			return errors.New(msg)
		}else{
			minerAddress = minerAddr
			currentWallets.SetMinerAddress(minerAddress)
		}
	}else{
		minerAddress = currentWallets.GetMinerAddress()
		if minerAddress == "" {
			mw := currentWallets.CreateWallet()
			err = currentWallets.SetMinerAddress(mw.GetStringAddress())
			if err != nil{
				fmt.Println(err)
				return err
			}
			minerAddress = mw.GetStringAddress()
		}
		if minerAddress == ""{
			msg := "not set miner address"
			fmt.Println(msg)
			return errors.New(msg)
		}
	}

	fmt.Println("[Important!!!] this node miner wallet address:", minerAddress)

	//load and check blockchain
	if isGenesisNode && isFirstInit{
		currentBlockChain = chain.CreateBlockchain(minerAddress, nodeID)
	}

	if currentBlockChain == nil{
		msg := "Blockchain Load error "
		fmt.Println(msg)
		return errors.New(msg)
	}


	//init addr manager
	currentAddrManager = NewAddrManager()
	currentAddrManager.AddAddress(knowAddr)



	//TODO:save to db?
	//init mempool
	txPool = mempool.New(currentBlockChain)
	return nil
}


// StartServer starts a node
func StartServer(nodeID, minerAddr string, isGenesisNode bool) error{

	err := initServer(nodeID, minerAddr, isGenesisNode)
	if err != nil{
		return err
	}

	//TODO:check config
	go LoopMining(currentBlockChain)

	//TODO:sync this node version info
	//TODO:sync block data


	//TODO:check config
	listenRPCServer(currentBlockChain)

	return nil

}

// CreateNodeID create node id with uuid
func CreateNodeID() string{
	return uuid.NewV4().String32()
}