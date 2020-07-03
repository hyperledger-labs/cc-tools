package assets

import (
	"encoding/json"
	"fmt"
	"strings"

	eh "github.com/goledgerdev/template-cc/chaincode/src/errorhandler"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

/*
Key implements the json.Unmarshaler interface
*/
type Key map[string]interface{}

/*
UnmarshalJSON parses JSON-encoded data and returns
*/
func (k *Key) UnmarshalJSON(jsonData []byte) error {
	*k = make(Key)
	var err error

	aux := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &aux)
	if err != nil {
		return eh.NewCCError(400, err.Error())
	}

	newKey, err := NewKey(aux)
	if err != nil {
		return err
	}

	*k = newKey

	return nil
}

// NewKey constructs Key object
func NewKey(m map[string]interface{}) (k Key, err eh.ICCError) {
	if m == nil {
		err = eh.NewCCError(500, "cannot create key from nil map")
		return
	}

	k = Key{}
	for t, v := range m {
		k[t] = v
	}

	// Generate object key
	_, keyExists := k["@key"]
	if !keyExists {
		var keyStr string
		keyStr, err = GenerateKey(k)
		if err != nil {
			err = eh.WrapError(err, "error generating key for asset")
		}
		k["@key"] = keyStr
	}

	return
}

// ValidateProps checks if all props are compliant to format
func (k Key) ValidateProps() error {
	// Perform validation of the @assetType field
	assetType, exists := k["@assetType"]
	if !exists {
		return eh.NewCCError(400, "property @assetType is required")
	}
	assetTypeString, ok := assetType.(string)
	if !ok {
		return eh.NewCCError(400, "property @assetType must be a string")
	}

	// Fetch asset definition
	assetTypeDef := FetchAssetType(assetTypeString)
	if assetTypeDef == nil {
		return eh.NewCCError(400, fmt.Sprintf("assetType named '%s' does not exist", assetTypeString))
	}

	// Validate asset properties
	for _, prop := range assetTypeDef.Keys() {
		// Check if required property is included
		propInterface, propIncluded := k[prop.Tag]
		if !propIncluded {
			if prop.Required {
				return eh.NewCCError(400, fmt.Sprintf("property %s (%s) is required", prop.Tag, prop.Label))
			}
			if prop.IsKey {
				return eh.NewCCError(400, fmt.Sprintf("key property %s (%s) is required", prop.Tag, prop.Label))
			}
			continue
		}

		// Validate data types
		propInterface, err := validateProp(propInterface, prop)
		if err != nil {
			msg := fmt.Sprintf("error validating asset '%s' property", prop.Tag)
			return eh.WrapError(err, msg)
		}

		k[prop.Tag] = propInterface
	}

	for propTag := range k {
		if strings.HasPrefix(propTag, "@") {
			continue
		}
		if !assetTypeDef.HasProp(propTag) {
			return eh.NewCCError(400, fmt.Sprintf("property %s is not defined in type %s", propTag, assetTypeString))
		}
	}

	return nil
}

// GetBytes reads asset bytes from ledger
func (k *Key) GetBytes(stub shim.ChaincodeStubInterface) ([]byte, eh.ICCError) {
	var assetBytes []byte
	var err error
	if k.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(k.TypeTag(), k.Key())
	} else {
		assetBytes, err = stub.GetState(k.Key())
	}
	if err != nil {
		return nil, eh.WrapErrorWithStatus(err, "failed to get asset bytes", 400)
	}
	if assetBytes == nil {
		return nil, eh.NewCCError(404, "asset not found")
	}

	return assetBytes, nil
}

// Get reads asset from ledger
func (k *Key) Get(stub shim.ChaincodeStubInterface) (*Asset, eh.ICCError) {
	var assetBytes []byte
	var err error
	if k.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(k.TypeTag(), k.Key())
	} else {
		assetBytes, err = stub.GetState(k.Key())
	}
	if err != nil {
		return nil, eh.WrapErrorWithStatus(err, "unable to get asset", 400)
	}
	if assetBytes == nil {
		return nil, eh.NewCCError(404, "asset not found")
	}

	var response Asset
	err = json.Unmarshal(assetBytes, &response)
	if err != nil {
		return nil, eh.WrapErrorWithStatus(err, "failed to unmarshal asset from ledger", 500)
	}

	return &response, nil
}

// Type return the AssetType object for the asset
func (k Key) Type() *AssetType {
	// Fetch asset properties
	assetTypeTag := k.TypeTag()
	assetDef := FetchAssetType(assetTypeTag)
	return assetDef
}

// IsPrivate returns true if asset type belongs to private collection
func (k Key) IsPrivate() bool {
	// Fetch asset properties
	assetTypeDef := k.Type()
	if assetTypeDef == nil {
		return false
	}
	if len(assetTypeDef.Readers) > 0 {
		return true
	}

	return false
}

// TypeTag returns @assetType attribute
func (k Key) TypeTag() string {
	assetType, _ := k["@assetType"].(string)
	return assetType
}

// Key returns the asset's unique key
func (k Key) Key() string {
	assetKey := k["@key"].(string)
	return assetKey
}
