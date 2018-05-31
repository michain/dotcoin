package proof

import (
	"testing"
	"fmt"
)

func TestNewProofOfWork(t *testing.T) {
	pow := NewProofOfWorkT(16)
	fmt.Println(pow.target)
}
