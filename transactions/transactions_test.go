package transactions

import (
	"log"
	"os"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile)

	InitHeader(Header{
		Name:    "CC Tools Test",
		Version: "v0.7.0",
		Colors: map[string][]string{
			"@default": {"#4267B2", "#34495E", "#ECF0F1"},
		},
		Title: map[string]string{
			"@default": "CC Tools Demo",
		},
	})

	InitTxList(testTxList)

	err := assets.CustomDataTypes(testCustomDataTypes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	assets.InitAssetList(testAssetList)

	os.Exit(m.Run())
}

func TestStartUp(t *testing.T) {
	err := StartupCheck()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
}
