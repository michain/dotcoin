package server

import "sync"

type AddrManager struct {
	mtx            *sync.Mutex
	addrIndex      map[string]string // address key to ka for all addrs.
}

func NewAddrManager() *AddrManager{
	return &AddrManager{
		mtx : new(sync.Mutex),
		addrIndex:make(map[string]string),
	}
}
