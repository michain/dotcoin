package base58

import (
	"testing"
)

func TestEncode(t *testing.T) {
	message := "test58"
	base58String := "zxo1qVyR"
	if base58String == Encode([]byte(message)){
		t.Log(message, base58String)
	}else{
		t.Error("encode string not match!")
	}
}


func TestDecode(t *testing.T) {
	message := "test58"
	base58String := "zxo1qVyR"
	if message == string(Decode(base58String)){
		t.Log(message, base58String)
	}else{
		t.Error("decode string not match!")
	}
}

func TestEecodeAlphabet(t *testing.T) {
	message := "test58"
	base58String := "ZXN1QuYq"
	if base58String == EncodeAlphabet([]byte(message), FlickrAlphabet){
		t.Log(message, base58String)
	}else{
		t.Error("encode FlickrAlphabet string not match!")
	}
}

func TestDecodeAlphabet(t *testing.T) {
	message := "test58"
	base58String := "ZXN1QuYq"
	if message == string(DecodeAlphabet(base58String, FlickrAlphabet)){
		t.Log(message, base58String)
	}else{
		t.Error("decode FlickrAlphabet string not match!")
	}
}