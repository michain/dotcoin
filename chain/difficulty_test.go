package chain

import (
	"testing"
	"math/big"
)

func TestBigToCompact(t *testing.T) {
	tests := []struct {
		in  int64
		out uint32
	}{
		{-1, 25231360},
		{0, 0},
		{1, 16842752},
		{2, 16908288},
		{10, 17432576},
		{11, 17498112},
		{100, 23330816},
		{101, 23396352},
		{200, 33605632},
		{1111, 33838848},
	}

	for x, test := range tests {
		n := big.NewInt(test.in)
		r := BigToCompact(n)
		n1 := CompactToBig(r)
		r1 := BigToCompact(n1)
		n2 := CompactToBig(r1)

		t.Log(n, r, n1, r1, n2)
		if r != test.out {
			t.Errorf("TestBigToCompact test #%d failed: got %d want %d\n",
				x, r, test.out)
			return
		}
	}
}

