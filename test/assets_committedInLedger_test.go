package test

import (
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/mock"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

func TestCommittedInLedgerNoKey(t *testing.T) {
	key := assets.Key{
		"@assetType": "person",
		// "@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	stub := mock.NewMockStub("org1MSP", new(testCC))
	stub.MockTransactionStart("TestCommittedInLedgerNoKey")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err := key.CommittedInLedger(sw)
	if err.Status() != 500 || err.Message() != "key is empty" {
		t.FailNow()
	}
	stub.MockTransactionEnd("TestCommittedInLedgerNoAssetKey")
}

func TestCommittedInLedgerNoAssetKey(t *testing.T) {
	asset := assets.Asset{
		// "@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@assetType": "person",
		"name":       "Maria",
		"id":         "31820792048",
	}

	stub := mock.NewMockStub("org1MSP", new(testCC))
	stub.MockTransactionStart("TestCommittedInLedgerNoAssetKey")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err := asset.CommittedInLedger(sw)
	if err.Status() != 500 || err.Message() != "asset key is empty" {
		t.FailNow()
	}
	stub.MockTransactionEnd("TestCommittedInLedgerNoAssetKey")
}
