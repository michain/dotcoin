package cli

import (
	"fmt"
	"os"
	"flag"
	"log"
	"github.com/michain/dotcoin/const"
)

const tcpPort = ":2398"
const(
	cmdVersion ="version"
	cmdCreateWallet = "createwallet"
	cmdGetBalance = "getbalance"
	cmdListAddresses = "listaddresses"
	cmdPrintChain = "printchain"
	cmdStartNode = "startnode"
	cmdPrintLogo = "printlogo"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  startnode -nodeid NodeID -miner ADDRESS -isgenesis IsGenesis  - - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 3 {
		cli.printUsage()
		os.Exit(1)
	}
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	versionCmd := flag.NewFlagSet(cmdVersion, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(cmdGetBalance, flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet(cmdCreateWallet, flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet(cmdListAddresses, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(cmdPrintChain, flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet(cmdStartNode, flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
	startNodeIsGenesis := startNodeCmd.Bool("isgenesis", false, "Set is isGenesis Mode")
	startNodeListen := startNodeCmd.String("listen", "", "Set listen addr")
	startNodeSeed := startNodeCmd.String("seed", "", "Set seed addr")


	nodeID:=os.Args[1]
	command:=os.Args[2]

	if nodeID == ""{
		fmt.Println("NodeID not set")
		os.Exit(1)
	}

	args := os.Args[3:]
	switch command {
	case cmdVersion:

	case cmdGetBalance:
		err := versionCmd.Parse(args)
		if err != nil {
			log.Panic(err)
		}
	case cmdCreateWallet:
		err := createWalletCmd.Parse(args)
		if err != nil {
			log.Panic(err)
		}
	case cmdListAddresses:
		err := listAddressesCmd.Parse(args)
		if err != nil {
			log.Panic(err)
		}
	case cmdPrintChain:
		err := printChainCmd.Parse(args)
		if err != nil {
			log.Panic(err)
		}
	case cmdStartNode:
		err := startNodeCmd.Parse(args)
		if err != nil {
			log.Panic(err)
		}
	case cmdPrintLogo:
		printLogo()
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if versionCmd.Parsed(){
		fmt.Printf("Version is %s\n", _const.Version)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress, nodeID)
	}


	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses(nodeID)
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)
	}

	if startNodeCmd.Parsed() {
		cli.startNode(nodeID, *startNodeMiner, *startNodeIsGenesis, *startNodeListen, *startNodeSeed)
	}
}
