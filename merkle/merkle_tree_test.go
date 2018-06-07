package merkle

import (
	"encoding/hex"
	"testing"
	"fmt"
	"crypto/sha256"
)


func TestSha256(t *testing.T){
	fmt.Println(sha256.Sum256(nil))
	fmt.Println(sha256.Sum256([]byte("")))
	fmt.Println(sha256.Sum256([]byte("node3")))
}

func TestNewMerkleNode(t *testing.T) {
	data := [][]byte{
		[]byte("trans1"),
		[]byte("trans2"),
		[]byte("trans3"),
	}

	// Level 1

	n1 := NewMerkleNode(nil, nil, data[0])
	n2 := NewMerkleNode(nil, nil, data[1])
	n3 := NewMerkleNode(nil, nil, data[2])
	n4 := NewMerkleNode(nil, nil, nil)

	// Level 2
	n5 := NewMerkleNode(n1, n2, nil)
	n6 := NewMerkleNode(n3, n4, nil)

	// Level 3
	n7 := NewMerkleNode(n5, n6, nil)

	if "2e1422e69f6fc07be3611b9c92e2dae791302e960d4245930f899c5c491386b5" == hex.EncodeToString(n5.Data){
		t.Log(hex.EncodeToString(n5.Data))
	}else{
		t.Error("Level 1 hash 1 is correct",hex.EncodeToString(n5.Data))
	}


	if "ffe3a30dc8951df5dfc9934ed61226ad642acd2dfd98dd57eb7391cca8988f5c" == hex.EncodeToString(n6.Data){
		t.Log(hex.EncodeToString(n6.Data))
	}else{
		t.Error("Level 1 hash 2 is correct",hex.EncodeToString(n6.Data))
	}


	if "3c75e1082ff231c46db08d4b69f2819aec56796aaf660ce791b08ee10f6e4167" == hex.EncodeToString(n7.Data){
		t.Log(hex.EncodeToString(n7.Data))
	}else{
		t.Error("Root hash is correct",hex.EncodeToString(n7.Data))
	}
}

func TestNewMerkleTree(t *testing.T) {
	data := [][]byte{
		[]byte("trans1"),
		[]byte("trans2"),
		[]byte("trans3"),
		[]byte("trans4"),
		[]byte("trans5"),
	}
	// Level 1
	n1 := NewMerkleNode(nil, nil, data[0])
	n2 := NewMerkleNode(nil, nil, data[1])
	n3 := NewMerkleNode(nil, nil, data[2])
	n4 := NewMerkleNode(nil, nil, data[3])
	n5 := NewMerkleNode(nil, nil, data[4])
	n6 := NewMerkleNode(nil, nil, nil)


	// for loop i=0
	// Level 2
	n11 := NewMerkleNode(n1, n2, nil)
	n12 := NewMerkleNode(n3, n4, nil)
	n13 := NewMerkleNode(n5, n6, nil)
	n14 := NewMerkleNode(nil, nil, nil)

	// for loop i=1
	// Level 3
	n21 := NewMerkleNode(n11, n12, nil)
	n22 := NewMerkleNode(n13, n14, nil)

	n31 := NewMerkleNode(n21, n22, nil)


	rootHash := fmt.Sprintf("%x", n31.Data)
	mTree := NewMerkleTree(data)

	if rootHash == fmt.Sprintf("%x", mTree.RootNode.Data){
		t.Log(rootHash)
		t.Log(fmt.Sprintf("%x", mTree.RootNode.Data))
	}else{
		t.Error("Merkle tree root hash is correct",)
	}
}