package assets

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// Key stores the information for retrieving an Asset from the ledger.
// Instead of having every asset property mapped such as the Asset type,
// Key only has the properties needed for fetching the full Asset.
type Key map[string]interface{}

// UnmarshalJSON parses JSON-encoded data and returns a Key object pointer
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

// NewKey constructs Key object from a map of properties.
// The map must contain the `@assetType` entry and either
// all the key properties of the asset (`IsKey == true`) or
// the `@key` property.
// Either way, the Key object returned contains only the
// `@assetType` and `@key` entries.
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
		if t != "@key" && t != "@assetType" {
			delete(k, t)
		}
	}

	return
}

// GetBytes reads the asset as bytes from ledger
func (k *Key) GetBytes(stub *sw.StubWrapper) ([]byte, errors.ICCError) {
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

// GetMap reads the asset as bytes from ledger
func (k *Key) GetMap(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {
	var err error
	assetBytes, err := k.GetBytes(stub)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to get asset bytes", 400)
	}

	var ret map[string]interface{}
	err = json.Unmarshal(assetBytes, &ret)
	if err != nil {
		return nil, errors.WrapError(err, "failed to unmarshal asset")
	}

	return ret, nil
}

// Type returns the AssetType configuration object for the asset
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

// Key returns the asset's unique identifying key in the ledger.
func (k Key) Key() string {
	assetKey := k["@key"].(string)
	return assetKey
}

// String returns the Key in stringified JSON format.
func (k Key) String() string {
	ret, _ := json.Marshal(k)
	return string(ret)
}

// JSON returns the Asset in JSON format.
func (k Key) JSON() []byte {
	ret, err := json.Marshal(k)
	if err != nil {
		return nil
	}
	return ret
}
