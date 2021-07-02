package test

import (
	"log"
	"reflect"
	"testing"

	tx "github.com/goledgerdev/cc-tools/transactions"
)

func TestGetArgDef(t *testing.T) {
	arg := tx.CreateAsset.Args.GetArgDef("asset")
	if arg == nil {
		log.Println("GetArgDef didn't find 'asset' arg")
		t.FailNow()
	}
	if !reflect.DeepEqual(*arg, tx.CreateAsset.Args[0]) {
		log.Println("GetArgDef failed to fetch arg")
		log.Printf("expected: %#v\n", tx.CreateAsset.Args[0])
		log.Printf("     got: %#v\n", arg)
		t.FailNow()
	}
}

func TestGetArgDefNil(t *testing.T) {
	arg := tx.CreateAsset.Args.GetArgDef("assets")
	if arg != nil {
		log.Println("GetArgDef found inexistent arg")
		t.FailNow()
	}
}
