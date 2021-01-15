package transactions

import (
	"log"
	"os"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
)

func TestMain(m *testing.M) {
	err := assets.CustomDataTypes(testCustomDataTypes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	assets.InitAssetList(testAssetList)

	InitTxList(testTxList)

	os.Exit(m.Run())
}

func TestStartUp(t *testing.T) {
	err := StartupCheck()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
}
