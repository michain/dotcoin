package storage

import (
	"testing"
)

func Test_GetDBFileName(t *testing.T){
	nodeID := "test"
	wantName := "dotchain_test.db"

	genFileName := GetDBFileName(nodeID)

	if genFileName != wantName{
		t.Error("dbfile name is not correct, want", wantName, ", gen", genFileName)
	}else{
		t.Log(genFileName)
	}
}
