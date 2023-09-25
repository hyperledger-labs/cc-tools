package transactions

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

type Header struct {
	Name           string
	Version        string
	Colors         map[string][]string
	Title          map[string]string
	CCToolsVersion string
}

var header Header

func InitHeader(h Header) {
	header = h
	header.CCToolsVersion = "v0.8.1"
}

// GetHeader returns data in CCHeader
var GetHeader = Transaction{
	Tag:         "getHeader",
	Label:       "Get Header",
	Description: "",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args:     ArgList{},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error
		colorMap := header.Colors
		nameMap := header.Title
		orgMSP, err := stub.GetMSPID()
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
			"name":           header.Name,
			"version":        header.Version,
			"orgMSP":         orgMSP,
			"colors":         colors,
			"orgTitle":       orgTitle,
			"ccToolsVersion": header.CCToolsVersion,
		}
		headerBytes, err := json.Marshal(header)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal header")
		}

		return headerBytes, nil
	},
}
