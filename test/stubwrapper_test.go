package test

import (
	"encoding/hex"
	"testing"

	"github.com/hyperledger-labs/cc-tools/v2/mock"
	sw "github.com/hyperledger-labs/cc-tools/v2/stubwrapper"
)

func TestGetPrivateDataHashFromWriteSet(t *testing.T) {
	stub := &sw.StubWrapper{
		Stub: mock.NewMockStub("org1MSP", new(testCC)),
		PvtWriteSet: map[string]map[string][]byte{
			"collection": {"key": []byte("abc"), "deleted": nil},
		},
	}

	hash, err := stub.GetPrivateDataHash("collection", "key")
	if err != nil {
		t.Fatal(err)
	}
	// SHA-256 of "abc", as computed by Fabric for private data hashes
	want := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if got := hex.EncodeToString(hash); got != want {
		t.Fatalf("got hash %s, want %s", got, want)
	}

	hash, err = stub.GetPrivateDataHash("collection", "deleted")
	if err != nil || hash != nil {
		t.Fatalf("got %x, %v for a deleted key; want nil, nil", hash, err)
	}
}
