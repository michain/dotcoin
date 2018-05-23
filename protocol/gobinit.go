package protocol

import "encoding/gob"

func init(){
	gob.Register(MsgAddr{})
	gob.Register(MsgGetAddr{})
	gob.Register(MsgInv{})
	gob.Register(MsgVersion{})
	gob.Register(MsgGetBlocks{})

}
