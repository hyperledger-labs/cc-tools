package test

import (
	"log"
	"reflect"
	"testing"
)

func TestGetPropDef(t *testing.T) {
	propDef := *testAssetList[0].GetPropDef("id")
	expectedPropDef := testAssetList[0].Props[0]

	if !reflect.DeepEqual(propDef, expectedPropDef) {
		log.Println("these should be deeply equal")
		log.Println(propDef)
		log.Println(expectedPropDef)
		t.FailNow()
	}
}
