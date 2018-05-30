package protocol

import "encoding/gob"

func init(){
	gob.Register(NetMessage{})
	gob.Register(MsgAddr{})
	gob.Register(MsgGetAddr{})
	gob.Register(MsgInv{})
	gob.Register(MsgVersion{})
	gob.Register(MsgGetBlocks{})
	gob.Register(MsgGetData{})
	gob.Register(MsgBlock{})
	gob.Register(MsgTx{})
}
