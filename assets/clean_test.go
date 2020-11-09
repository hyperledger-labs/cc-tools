package assets

import (
	"reflect"
	"testing"
)

func TestClean(t *testing.T) {
	a := Asset{
		"name": "Maria",
		"id":   0,
		"opt":  nil,
	}

	b := Asset{
		"name": "Maria",
		"id":   0,
	}

	a.clean()

	if !reflect.DeepEqual(a, b) {
		t.Fail()
	}
}
