package transactions

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
)

type Header struct {
	Name    string
	Version string
	Colors  map[string][]string
	Title   map[string]string
}

var header Header

func InitHeader(h Header) {
	header = h
}

// GetHeader returns data in CCHeader
var GetHeader = Transaction{
	Tag:         "getHeader",
	Label:       "Get Header",
	Description: "",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args:     []Argument{},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		colorMap := header.Colors
		nameMap := header.Title
		orgMSP, err := cid.GetMSPID(stub.Stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get MSP ID")
		}

		var colors []string
		colors, orgExists := colorMap[orgMSP]
		if !orgExists {
			colors = colorMap["@default"]
		}

		var orgTitle string
		orgTitle, orgExists = nameMap[orgMSP]
		if !orgExists {
			orgTitle = nameMap["@default"]
		}

		header := map[string]interface{}{
			"name":     header.Name,
			"version":  header.Version,
			"orgMSP":   orgMSP,
			"colors":   colors,
			"orgTitle": orgTitle,
		}
		headerBytes, err := json.Marshal(header)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal header")
		}

		return headerBytes, nil
	},
}
