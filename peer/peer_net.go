package peer

import (
	"encoding/gob"
	"net"
)

// WriteConnRequest write request data to conn
func WriteConnRequest(conn net.Conn, r interface{}) error{
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(r)
	return err
}

// ReadConnRequest read request data from conn
func ReadConnRequest(conn net.Conn) (*Request, error){
	decoder := gob.NewDecoder(conn)
	r := &Request{}
	err := decoder.Decode(r)
	return r, err
}

