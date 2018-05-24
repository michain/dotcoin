package peer

import (
	"encoding/gob"
	"net"
	"bytes"
	"github.com/michain/dotcoin/connx"
)

// WriteConnRequest write request data to conn
func WriteConnRequest(conn net.Conn, r interface{}) error{
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(r)
	if err != nil {
		return err
	}
	return connx.SendMessage(conn, encoded.Bytes())
}

// ReadConnRequest read request data from conn
func ReadConnRequest(conn net.Conn) (*Request, error){
	data, err := connx.ReadMessage(conn)
	if err !=nil{
		return nil, err
	}
	r := &Request{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(r)
	return r, err
}




