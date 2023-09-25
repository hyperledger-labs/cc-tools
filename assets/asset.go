package assets

import (
	"encoding/json"
	"time"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// Asset implements the json.Unmarshaler interface and is the base object in cc-tools network.
// It is a generic map that stores information about a specific ledger asset. It is also used
// as the base interface to perform operations on the blockchain.
type Asset map[string]interface{}

// UnmarshalJSON parses JSON-encoded data and returns an Asset object pointer
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

// NewAsset constructs Asset object from a map of properties. It ensures every
// asset property is properly formatted and computes the asset's identifying
// key on the ledger, storing it in the property "@key".
func NewAsset(m map[string]interface{}) (a Asset, err errors.ICCError) {
	if m == nil {
		err = errors.NewCCError("cannot create asset from nil map", 500)
		return
	}

	a = Asset{}
	for k, v := range m {
		if v != nil {
			a[k] = v
		}
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

// InjectMetadata adds internal properties to the asset.
func (a *Asset) injectMetadata(stub *sw.StubWrapper) errors.ICCError {
	var err error

	// Generate object key
	key, err := GenerateKey(*a)
	if err != nil {
		return errors.WrapError(err, "error generating key for asset")
	}
	(*a)["@key"] = key

	lastTouchBy, err := stub.GetMSPID()
	if err != nil {
		return errors.WrapErrorWithStatus(err, "error getting tx creator", 500)
	}
	(*a)["@lastTouchBy"] = lastTouchBy

	lastTx, _ := stub.Stub.GetFunctionAndParameters()
	(*a)["@lastTx"] = lastTx

	lastUpdated, err := stub.Stub.GetTxTimestamp()
	if err != nil {
		return errors.WrapErrorWithStatus(err, "error getting tx timestamp", 500)
	}
	(*a)["@lastUpdated"] = lastUpdated.AsTime().Format(time.RFC3339)

	return nil
}

// IsPrivate returns true if the Asset's asset type belongs to a private collection.
func (a Asset) IsPrivate() bool {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return false
	}
	return assetTypeDef.IsPrivate()
}

// TypeTag returns the @assetType attribute.
func (a Asset) TypeTag() string {
	assetType, _ := a["@assetType"].(string)
	return assetType
}

// Key returns the asset's unique identifying key in the ledger.
func (a Asset) Key() string {
	assetKey, _ := a["@key"].(string)
	return assetKey
}

// Type return the AssetType object for the asset.
func (a Asset) Type() *AssetType {
	// Fetch asset properties
	assetTypeTag := a.TypeTag()
	assetTypeDef := FetchAssetType(assetTypeTag)
	return assetTypeDef
}

// String returns the Asset in stringified JSON format.
func (a Asset) String() string {
	ret, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(ret)
}

// JSON returns the Asset in JSON format.
func (a Asset) JSON() []byte {
	ret, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	return ret
}
