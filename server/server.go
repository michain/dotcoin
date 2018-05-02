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
	"github.com/michain/dotcoin/sync"
	"github.com/michain/dotcoin/logx"
	"github.com/michain/dotcoin/peer"
)

var curTXMemPool *mempool.TxPool
var curWallets *wallet.WalletSet
var curBlockChain *chain.Blockchain
var curAddrManager *AddrManager
var curSyncManager *sync.SyncManager
var curPeer *peer.Peer
var minerAddress string
var listenAddress string
var seedAddrs []string

const(
	rpcPort = ":12398" //2398 = 1983+0415 my birthday!
	tcpPort = ":2398"
	coinbaseReward = 10
	knowAddr = "localhost:2398"
)

func init(){
	seedAddrs = []string{knowAddr}
	listenAddress = tcpPort
}


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

func listenPeer(){
	if len(seedAddrs) <= 0{
		log.Fatalf("listenPeer error: seedAddrs is nil")
	}
	logx.Debugf("listenPeer begin listen:%v seed:%v", tcpPort, seedAddrs[0])
	curPeer = peer.NewPeer(tcpPort, seedAddrs[0], NewMessageHandler(curPeer))
	curPeer.StartListen()
}

// initServer init server
func initServer(nodeID, minerAddr string, isGenesisNode bool) error{
	listenAddress = fmt.Sprintf("localhost:%s", nodeID)
	var err error
	isFirstInit := false

	curBlockChain, err = chain.LoadBlockChain(nodeID)
	if err == chain.ErrorBlockChainNotFount{
		isFirstInit = true
	}
	if !isGenesisNode && isFirstInit{
		//TODO:sync data from other node
		//if err, block myself
	}


	//load or create miner wallet
	curWallets, err = wallet.LoadWallets(nodeID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//ignore config miner address when genesis and first init
	if isGenesisNode && isFirstInit{
		mw := curWallets.CreateWallet()
		err = curWallets.SetMinerAddress(mw.GetStringAddress())
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
			curWallets.SetMinerAddress(minerAddress)
		}
	}else{
		minerAddress = curWallets.GetMinerAddress()
		if minerAddress == "" {
			mw := curWallets.CreateWallet()
			err = curWallets.SetMinerAddress(mw.GetStringAddress())
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
		curBlockChain = chain.CreateBlockchain(minerAddress, nodeID)
	}

	if curBlockChain == nil{
		msg := "Blockchain Load error "
		return errors.New(msg)
	}


	//init addr manager
	curAddrManager = NewAddrManager()
	curAddrManager.AddAddress(knowAddr)

	//init sync manager
	curSyncManager, err = sync.New(&sync.Config{
		Chain : curBlockChain,
		TxMemPool:curTXMemPool,
		MaxPeers:MaxPeerNum,
		Peer:curPeer,
	})
	if err!= nil{
		return err
	}
	go curSyncManager.Start()


	//TODO:save to db?
	//init mempool
	curTXMemPool = mempool.New(curBlockChain)
	return nil
}


// StartServer starts a node
func StartServer(nodeID, minerAddr string, isGenesisNode bool) error{

	err := initServer(nodeID, minerAddr, isGenesisNode)
	if err != nil{
		return err
	}

	//TODO:check config
	go LoopMining(curBlockChain)

	//TODO:sync this node version info
	//TODO:sync block data

	//start peer
	go listenPeer()

	//TODO:check config
	listenRPCServer(curBlockChain)

	return nil

}

// CreateNodeID create node id with uuid
func CreateNodeID() string{
	return uuid.NewV4().String32()
}