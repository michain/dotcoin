package proof

import (
	"testing"
	"fmt"
)

func TestNewProofOfWork(t *testing.T) {
	pow := NewProofOfWork()
	fmt.Println(pow.target)
}
