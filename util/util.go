package util

import (
	"encoding/binary"
	"bytes"
	"log"
	"os"
	"fmt"
	"math/rand"
)

func GetRandData() string{
	randData := make([]byte, 20)
	_, err := rand.Read(randData)
	if err != nil {
		log.Panic(err)
	}

	return fmt.Sprintf("%x", randData)
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