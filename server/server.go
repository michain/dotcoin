package server

import (
	"github.com/michain/dotcoin/sync"
	"net"
	"net/rpc"
	"fmt"
	"net/rpc/jsonrpc"
	"github.com/michain/dotcoin/peer"
	"github.com/michain/dotcoin/protocol"
	"github.com/michain/dotcoin/mempool"
	"github.com/michain/dotcoin/wallet"
	"github.com/michain/dotcoin/chain"
	"log"
	"github.com/michain/dotcoin/logx"
	"errors"
	"github.com/michain/dotcoin/util/uuid"
)

/*var curNodeID string
var curTXMemPool *mempool.TxPool
var curWallets *wallet.WalletSet
var curBlockChain *chain.Blockchain
var curAddrManager *AddrManager
var curSyncManager *sync.SyncManager
var curPeer *peer.Peer
var minerAddress string
var listenAddress string*/
var seedAddrs []string
var curServer *Server


const(
	rpcPort = ":12398" //2398 = 1983+0415 my birthday!

	coinbaseReward = 10
	knowAddr = "localhost:2398"
)

func init(){
	seedAddrs = []string{knowAddr}
	curServer = new(Server)
}

type Server struct{
	NodeID string
	ListenAddress string
	SeedAddress string
	TXMemPool *mempool.TxPool
	Wallets *wallet.WalletSet
	BlockChain *chain.Blockchain
	AddrManager *AddrManager
	SyncManager *sync.SyncManager
	Peer *peer.Peer
	IsGenesisNode bool
	minerAddress string

}


func (s *Server) listenRPCServer() {
	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		return
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("Rpc", &RpcHandler{server:s}); err != nil {
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

func (s *Server) listenPeer(){
	if !s.IsGenesisNode && s.SeedAddress == ""{
		log.Fatalf("listenPeer error: SeedAddress is nil")
	}
	logx.Debugf("listenPeer begin listen:%v seed:%v", s.ListenAddress, s.SeedAddress)
	s.Peer.StartListen()
}

// initServer init server
func initServer(nodeID, minerAddr string, listenAddr, seedAddr string, isGenesisNode bool) (*Server, error){
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("[InitServer] Begin node server:", nodeID, minerAddr, isGenesisNode)
	serv := new(Server)
	serv.ListenAddress = listenAddr
	serv.SeedAddress = seedAddr
	serv.NodeID = nodeID
	serv.IsGenesisNode = isGenesisNode
	var err error
	isFirstInit := false

	serv.BlockChain, err = chain.LoadBlockChain(nodeID)
	if err == chain.ErrorBlockChainNotFount{
		isFirstInit = true
	}
	if !isGenesisNode && isFirstInit{
		//TODO:sync data from other node
		//if err, block myself
	}


	//load or create miner wallet
	serv.Wallets, err = wallet.LoadWallets(nodeID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//ignore config miner address when genesis and first init
	if isGenesisNode && isFirstInit{
		mw := serv.Wallets.CreateWallet()
		err = serv.Wallets.SetMinerAddress(mw.GetStringAddress())
		if err != nil{
			fmt.Println(err)
			return nil, err
		}
		serv.minerAddress = mw.GetStringAddress()
	}
	//if set minerAddr, validate it
	if minerAddr != ""{
		if !wallet.ValidateAddress(minerAddr) {
			msg := "Validate minerAddr error " + minerAddr
			fmt.Println(msg)
			return nil, errors.New(msg)
		}else{
			serv.minerAddress = minerAddr
			serv.Wallets.SetMinerAddress(serv.minerAddress)
		}
	}else{
		serv.minerAddress = serv.Wallets.GetMinerAddress()
		if serv.minerAddress == "" {
			mw := serv.Wallets.CreateWallet()
			err = serv.Wallets.SetMinerAddress(mw.GetStringAddress())
			if err != nil{
				fmt.Println(err)
				return nil, err
			}
			serv.minerAddress = mw.GetStringAddress()
		}
		if serv.minerAddress == ""{
			msg := "not set miner address"
			fmt.Println(msg)
			return nil, errors.New(msg)
		}
	}

	fmt.Println("[Important!!!] this node miner wallet address:", serv.minerAddress, "IsFirstInit", isFirstInit)

	//load and check blockchain
	if isFirstInit{
		serv.BlockChain = chain.CreateBlockchain(isGenesisNode, serv.minerAddress, nodeID)
	}

	if serv.BlockChain == nil{
		msg := serv.NodeID + " Blockchain Load error "
		return nil, errors.New(msg)
	}


	//init addr manager
	serv.AddrManager = NewAddrManager()
	serv.AddrManager.AddAddress(knowAddr)

	//init peer
	serv.Peer = peer.NewPeer(serv.ListenAddress, serv.SeedAddress, NewMessageHandler(serv))

	//init sync manager
	serv.SyncManager, err = sync.New(&sync.Config{
		Chain : serv.BlockChain,
		TxMemPool:serv.TXMemPool,
		MaxPeers:MaxPeerNum,
		Peer:serv.Peer,
	})
	if err!= nil{
		return nil, err
	}

	//TODO:save to db?
	//init mempool
	serv.TXMemPool = mempool.New(serv.BlockChain)
	return serv, nil
}


// StartServer starts a node
func StartServer(nodeID, minerAddr string, listenAddr, seedAddr string, isGenesisNode bool) error{

	serv, err := initServer(nodeID, minerAddr, listenAddr, seedAddr, isGenesisNode)
	if err != nil{
		return err
	}
	curServer = serv
	//start peer
	go serv.listenPeer()

	//start sync loop
	go serv.SyncManager.StartSync()

	//TODO:check config
	go serv.LoopMining()

	if !isGenesisNode {
		//send this node version info
		msg := protocol.NewMsgVersion(curServer.BlockChain.GetBestHeight())
		msg.AddrFrom = serv.Peer.GetSeedAddr()
		curServer.Peer.SendSingleMessage(msg)
	}

	//TODO:sync block data


	//TODO:check config
	serv.listenRPCServer()

	return nil

}

// CreateNodeID create node id with uuid
func CreateNodeID() string{
	return uuid.NewV4().String32()
}