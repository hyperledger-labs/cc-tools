package assets

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := CustomDataTypes(testCustomDataTypes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	InitAssetList(testAssetList)

	os.Exit(m.Run())
}

func TestStartUp(t *testing.T) {
	err := StartupCheck()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
}

func TestAssetList(t *testing.T) {
	l := AssetTypeList()
	if len(l) != 4 {
		fmt.Println("expected only 3 asset types in asset type list")
		t.FailNow()
	}
}
