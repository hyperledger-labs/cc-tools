package assets

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
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
		return errors.NewCCError(err.Error(), 400)
	}

	newKey, err := NewKey(aux)
	if err != nil {
		return err
	}

	*k = newKey

	return nil
}

// NewKey constructs Key object
func NewKey(m map[string]interface{}) (k Key, err errors.ICCError) {
	if m == nil {
		err = errors.NewCCError("cannot create key from nil map", 500)
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
			err = errors.WrapError(err, "error generating key for asset")
		}
		k["@key"] = keyStr
	}

	for t := range k {
		if t == "@key" || t == "@assetType" {
			continue
		}
		delete(k, t)
	}

	return
}

// ValidateProps checks if all props are compliant to format
func (k Key) ValidateProps() error {
	// Perform validation of the @assetType field
	assetType, exists := k["@assetType"]
	if !exists {
		return errors.NewCCError("property @assetType is required", 400)
	}
	assetTypeString, ok := assetType.(string)
	if !ok {
		return errors.NewCCError("property @assetType must be a string", 400)
	}

	// Fetch asset definition
	assetTypeDef := FetchAssetType(assetTypeString)
	if assetTypeDef == nil {
		return errors.NewCCError(fmt.Sprintf("assetType named '%s' does not exist", assetTypeString), 400)
	}

	// Validate asset properties
	for _, prop := range assetTypeDef.Keys() {
		// Check if required property is included
		propInterface, propIncluded := k[prop.Tag]
		if !propIncluded {
			if prop.Required {
				return errors.NewCCError(fmt.Sprintf("property %s (%s) is required", prop.Tag, prop.Label), 400)
			}
			if prop.IsKey {
				return errors.NewCCError(fmt.Sprintf("key property %s (%s) is required", prop.Tag, prop.Label), 400)
			}
			continue
		}

		// Validate data types
		propInterface, err := validateProp(propInterface, prop)
		if err != nil {
			msg := fmt.Sprintf("error validating asset '%s' property", prop.Tag)
			return errors.WrapError(err, msg)
		}

		k[prop.Tag] = propInterface
	}

	for propTag := range k {
		if strings.HasPrefix(propTag, "@") {
			continue
		}
		if !assetTypeDef.HasProp(propTag) {
			return errors.NewCCError(fmt.Sprintf("property %s is not defined in type %s", propTag, assetTypeString), 400)
		}
	}

	return nil
}

// GetBytes reads asset bytes from ledger
func (k *Key) GetBytes(stub shim.ChaincodeStubInterface) ([]byte, errors.ICCError) {
	var assetBytes []byte
	var err error
	if k.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(k.TypeTag(), k.Key())
	} else {
		assetBytes, err = stub.GetState(k.Key())
	}
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to get asset bytes", 400)
	}
	if assetBytes == nil {
		return nil, errors.NewCCError("asset not found", 404)
	}

	return assetBytes, nil
}

// Get reads asset from ledger
func (k *Key) Get(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var assetBytes []byte
	var err error
	if k.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(k.TypeTag(), k.Key())
	} else {
		assetBytes, err = stub.GetState(k.Key())
	}
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "unable to get asset", 400)
	}
	if assetBytes == nil {
		return nil, errors.NewCCError("asset not found", 404)
	}

	var response Asset
	err = json.Unmarshal(assetBytes, &response)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal asset from ledger", 500)
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
	return assetTypeDef.IsPrivate()
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

func (k Key) String() string {
	ret, _ := json.Marshal(k)
	return string(ret)
}
