package connx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"github.com/pkg/errors"
)

var ErrorNotMatchHeadFlag = errors.New("not match head flag")

// SendMessage send message with conn
func SendMessage(conn net.Conn, msgData []byte) error{
	bufFlag := new(bytes.Buffer)
	binary.Write(bufFlag, binary.LittleEndian, HeadFlag)

	head := &HeadInfo{}
	head.head_id = 34969
	head.data_type = 10101
	head.data_id = 1
	head.data_len = uint64(len(msgData))
	headb := head.GetBytes()
	send := make([]byte, 0)
	send = slice_merge(send, bufFlag.Bytes())
	send = slice_merge(send, headb)
	send = slice_merge(send, msgData)
	_, err:=conn.Write(send)
	return err
}

// ReadMessage read message from conn
func ReadMessage(conn net.Conn) ([]byte, error){
	//read HeadFlag
	flag, err:= readHeadFlag(conn)
	if err!= nil{
		if err == io.EOF {
		} else {
			fmt.Println("ReadDate:readHeadFlag error -> ", err)
		}
		return nil, err
	}
	if flag != HeadFlag{
		fmt.Println("ReadDate:readHeadFlag not match -> ", flag, HeadFlag)
		//not match, go to next read
		return nil, ErrorNotMatchHeadFlag
	}
	//read head
	head, err := readHead(conn)
	if err != nil {
		if err == io.EOF {
		} else {
			fmt.Println("ReadDate:readHead error -> ", err)
		}
		return nil, err
	}

	//read msg body
	var bufData []byte
	errRead := readSize(conn, int64(head.data_len), &bufData)
	if errRead == nil {
		return bufData, nil
	} else {
		return nil, err
	}
}
