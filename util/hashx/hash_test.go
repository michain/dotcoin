package hashx

import (
	"testing"
	"fmt"
)

func TestHash_IsEqual(t *testing.T) {
	hOne := *ZeroHash()

	fmt.Println(hOne.IsEqual(ZeroHash()))

	fmt.Println(hOne)
	fmt.Println(ZeroHash())
}
