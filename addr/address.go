package addr

import (
	"sync"
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

// HasAddress check is exists address
func (m *AddrManager) HasAddress(addr string) bool{
	_, exists:=m.addrIndex[addr]
	return exists
}

// AddAddress add addr into manager
// if exists, convert it
func (m *AddrManager) AddAddress(addr string){
	m.mtx.Lock()
	m.addrIndex[addr] = addr
	m.mtx.Unlock()
}

// GetAddresses get all address in manager
func (addr *AddrManager) GetAddresses() []string{
	addrs := []string{}
	for k, _:=range addr.addrIndex{
		addrs = append(addrs, k)
	}
	return addrs
}
