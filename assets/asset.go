package assets

import (
	"encoding/json"
	"fmt"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*
Asset implements the json.Unmarshaler interface
*/
type Asset map[string]interface{}

/*
UnmarshalJSON parses JSON-encoded data and returns
*/
func (a *Asset) UnmarshalJSON(jsonData []byte) error {
	var err error

	aux := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &aux)
	if err != nil {
		return errors.NewCCError(err.Error(), 400)
	}

	newAsset, err := NewAsset(aux)
	if err != nil {
		return err
	}

	*a = newAsset

	return nil
}

// NewAsset constructs Asset object
func NewAsset(m map[string]interface{}) (a Asset, err errors.ICCError) {
	if m == nil {
		err = errors.NewCCError("cannot create asset from nil map", 500)
		return
	}

	a = Asset{}
	for k, v := range m {
		a[k] = v
	}

	// Generate object key
	key, err := GenerateKey(a)
	if err != nil {
		err = errors.WrapError(err, "error generating key for asset")
		return
	}
	(a)["@key"] = key

	// Filter, validate and convert props to proper format
	err = a.ValidateProps()
	if err != nil {
		err = errors.WrapError(err, "format error")
		return
	}

	return
}

// InjectMetadata handles injection internal variables to the asset
func (a *Asset) InjectMetadata(stub shim.ChaincodeStubInterface) error {
	// Generate object key
	if _, keyExists := (*a)["@key"]; !keyExists {
		key, err := GenerateKey(*a)
		if err != nil {
			return fmt.Errorf("error generating key for asset: %s", err)
		}
		(*a)["@key"] = key
	}

	lastTouchBy, err := cid.GetMSPID(stub)
	if err != nil {
		return fmt.Errorf("error getting tx creator: %s", err)
	}
	(*a)["@lastTouchBy"] = lastTouchBy

	return nil
}

// IsPrivate returns true if asset type belongs to private collection
func (a Asset) IsPrivate() bool {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return false
	}
	if len(assetTypeDef.Readers) > 0 {
		return true
	}

	return false
}

// TypeTag returns @assetType attribute
func (a Asset) TypeTag() string {
	assetType, _ := a["@assetType"].(string)
	return assetType
}

// Key returns the asset's unique key
func (a Asset) Key() string {
	assetKey, _ := a["@key"].(string)
	return assetKey
}

// Type return the AssetType object for the asset
func (a Asset) Type() *AssetType {
	// Fetch asset properties
	assetTypeTag := a.TypeTag()
	assetTypeDef := FetchAssetType(assetTypeTag)
	return assetTypeDef
}
