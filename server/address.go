package server

import (
	"sync"
	"syscall"
)

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

// AddAddress add addr into manager
// if exists, convert it
func (a *AddrManager) AddAddress(addr string){
	a.mtx.Lock()
	a.addrIndex[addr] = addr
	a.mtx.Unlock()
}

// GetAddresses get all address in manager
func (a *AddrManager) GetAddresses() []string{
	addrs := []string{}
	for k, _:=range a.addrIndex{
		addrs = append(addrs, k)
	}
	return addrs
}
