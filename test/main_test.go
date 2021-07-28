package test

import (
	"log"
	"os"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile)

	tx.InitHeader(tx.Header{
		Name:    "CC Tools Test",
		Version: "v0.7.1",
		Colors: map[string][]string{
			"@default": {"#4267B2", "#34495E", "#ECF0F1"},
		},
		Title: map[string]string{
			"@default": "CC Tools Demo",
		},
	})

	tx.InitTxList(testTxList)

	err := assets.CustomDataTypes(testCustomDataTypes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	assets.InitAssetList(testAssetList)

	err = assets.StartupCheck()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = tx.StartupCheck()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
