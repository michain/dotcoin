package peer

import (
	"testing"
)

var seed = "127.0.0.1:2398"
var pnode1 = "127.0.0.1:2391"
var pnode1_1 = "127.0.0.1:2392"
var pnode2 = "127.0.0.1:2491"
var pnode2_1 = "127.0.0.1:2492"

func TestStartPeer(t *testing.T) {

	go func() {
		p := NewPeer(seed, "", nil)
		err := p.StartListen()
		if err != nil {
			t.Error("Seed Peer start error", err)
		} else {
			t.Log("Seed Peer start success")
		}
	}()

	go func() {
		p := NewPeer(pnode1, seed, nil)
		err := p.StartListen()
		if err != nil {
			t.Error("pnode1 Peer start error", err)
		} else {
			t.Log("pnode1 Peer start success")
		}
	}()

	go func() {
		p := NewPeer(pnode1_1, pnode1, nil)
		err := p.StartListen()
		if err != nil {
			t.Error("pnode1_1 Peer start error", err)
		} else {
			t.Log("pnode1_1 Peer start success")
		}
	}()


	go func() {
		p := NewPeer(pnode2, seed, nil)
		err := p.StartListen()
		if err != nil {
			t.Error("pnode2 Peer start error", err)
		} else {
			t.Log("pnode2 Peer start success")
		}
	}()

	var p_2_1 *Peer
	go func() {
		var err error
		p_2_1 = NewPeer(pnode2_1, pnode2, nil)
		err = p_2_1.StartListen()
		if err != nil {
			t.Error("pnode2_1 Peer start error", err)
		} else {
			t.Log("pnode2_1 Peer start success")
		}
	}()

	for{
		select{}
	}
}
