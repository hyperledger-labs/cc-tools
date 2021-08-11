package test

import (
	"log"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/mock"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

func TestCheckWriters(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	sw := &sw.StubWrapper{
		Stub: stub,
	}

	a := assets.Asset{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"id":         "31820792048",
		"name":       "Maria",
		"height":     0.0,
	}

	err := a.CheckWriters(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	stub.Name = "org2MSP"

	err = a.CheckWriters(sw)
	if err == nil {
		log.Println("expected asset.CheckWriters to fail")
		t.FailNow()
	}
	if err.Status() != 403 {
		log.Println("expected err: 403")
		log.Printf("got err: %d", err.Status())
		t.FailNow()
	}
	if err.Message() != `org2MSP cannot write to the 'id' (CPF (Brazilian ID)) asset property` {
		log.Println(`expected err: org2MSP cannot write to the 'id' (CPF (Brazilian ID)) asset property`)
		log.Printf("got err: %s", err.Message())
		t.FailNow()
	}
}
