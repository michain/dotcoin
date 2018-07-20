package wallet

import (
	"testing"
	"fmt"
	"github.com/michain/dotcoin/base58"
	"os"
	"syscall"
	"path/filepath"
	"crypto/elliptic"
	"encoding/gob"
	"bytes"
	"log"
)

func Test_newKeyPair(t *testing.T){
	private, pubKey := newKeyPair()
	fmt.Println(private)
	fmt.Println(base58.Encode(pubKey))
}


func Test_newWallet(t *testing.T){
	//create 30 wallets
	for i:=0;i<30;i++ {
		w := newWallet()
		address := string(w.GetAddress())
		data := getEncodeBytes(w)
		fmt.Println(save2File(address, data))
	}
}


func save2File(address string, data []byte) (string, error){
	pathDir := filepath.Dir("d://dotechnology/genesis-coin/")
	logFile := "d://dotechnology/genesis-coin/" + address + ".key"
	if !existFile(pathDir) {
		//create path
		err := os.MkdirAll(pathDir, 0777)
		if err != nil {
			fmt.Println("save2File create path error ", err)
			return logFile, err
		}
	}

	var mode os.FileMode
	flag := syscall.O_RDWR | syscall.O_APPEND | syscall.O_CREAT
	mode = 0666
	file, err := os.OpenFile(logFile, flag, mode)
	defer file.Close()
	if err != nil {
		fmt.Println(logFile, err)
		return logFile, err
	}
	file.Write(data)
	return logFile, nil
}

func getEncodeBytes(w  *Wallet) []byte{
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	if err != nil {
		log.Panic(err)
	}
	return content.Bytes()
}

//check filename is exist
func existFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
