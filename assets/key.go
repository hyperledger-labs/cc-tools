package assets

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Key implements the json.Unmarshaler interface
//It stores the information for retrieving assets from the ledger
//Instead of having every field, it only has the ones needes for querying
type Key map[string]interface{}

//UnmarshalJSON parses JSON-encoded data and returns
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

// NewKey constructs Key object from a map
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

// Type return the AssetType object configuration for the asset
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
