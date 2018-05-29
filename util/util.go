package util

import (
	"encoding/binary"
	"bytes"
	"log"
	"os"
	"github.com/michain/dotcoin/util/uuid"
)

func GetRandData() string{
	return uuid.NewV4().String()
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}


func ExitFile(fileName string) bool{
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}

	return true
}