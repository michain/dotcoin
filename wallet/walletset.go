package wallet

import (
	"fmt"
	"io/ioutil"
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"bytes"
	"github.com/michain/dotcoin/util"
	"github.com/pkg/errors"
	"sync"
)

var (
	ErrorNotExistsWalletFile = errors.New("not exists wallet file!")
)

const walletFileFormat = "dotwallet_%s.key"
// WalletMap stores a collection of wallets
type WalletSet struct {
	mutex *sync.RWMutex
	NodeID string
	MinerAddress string
	Wallets map[string]*Wallet
}

// LoadWallets load Wallets and fills it from a file
// if not exists file, auto create it
func LoadWallets(nodeID string) (*WalletSet, error) {
	wallets := WalletSet{NodeID:nodeID}
	wallets.Wallets = make(map[string]*Wallet)
	wallets.mutex = new(sync.RWMutex)

	err := wallets.LoadFromFile(nodeID)
	if err == ErrorNotExistsWalletFile{
		wallets.SaveToFile()
		err = nil
	}
	return &wallets, err
}


// SetMinerWallet set miner address into current WalletSet
func (ws *WalletSet) SetMinerAddress(addr string) error{
	if !ValidateAddress(addr){
		return errors.New("Validate address error " + addr)
	}

	ws.mutex.Lock()
	ws.mutex.Unlock()
	ws.MinerAddress = addr
	ws.SaveToFile()

	return nil
}

// GetMinerWallet get miner address from current WalletSet
func (ws *WalletSet) GetMinerAddress() string{
	return ws.MinerAddress
}


// CreateWallet adds a Wallet to Wallets
func (ws *WalletSet) CreateWallet() *Wallet {
	wallet := newWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	ws.Wallets[address] = wallet
	//save to file
	ws.SaveToFile()
	return wallet
}

// GetAddresses returns an array of addresses stored in the wallet file
func (ws *WalletSet) GetAddresses() []string {
	var addresses []string

	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet returns a Wallet by its address
// if not exists, return nil
func (ws *WalletSet) GetWallet(address string) *Wallet {
	return ws.Wallets[address]
}

// LoadFromFile loads wallets from the file
func (ws *WalletSet) LoadFromFile(nodeID string) error {
	walletFile := getWalletFileName(nodeID)
	if !util.ExitFile(walletFile){
		return ErrorNotExistsWalletFile
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets WalletSet
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	ws.MinerAddress = wallets.MinerAddress
	ws.mutex = new(sync.RWMutex)
	ws.NodeID = nodeID

	return nil
}

// SaveToFile saves wallets to a file
func (ws *WalletSet) SaveToFile() {
	var content bytes.Buffer
	walletFile := getWalletFileName(ws.NodeID)

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}


// getWalletFileName get wallet file's name with NodeID
// like dotwallet__XXXXXX.dat
func getWalletFileName(nodeID string) string{
	return fmt.Sprintf(walletFileFormat, nodeID)
}

