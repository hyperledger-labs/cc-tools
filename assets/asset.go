package assets

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	eh "github.com/goledgerdev/template-cc/chaincode/src/errorhandler"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
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
		return eh.NewCCError(400, err.Error())
	}

	newAsset, err := NewAsset(aux)
	if err != nil {
		return err
	}

	*a = newAsset

	return nil
}

// NewAsset constructs Asset object
func NewAsset(m map[string]interface{}) (a Asset, err eh.ICCError) {
	if m == nil {
		err = eh.NewCCError(500, "cannot create asset from nil map")
		return
	}

	a = Asset{}
	for k, v := range m {
		a[k] = v
	}

	// Generate object key
	key, err := GenerateKey(a)
	if err != nil {
		err = eh.WrapError(err, "error generating key for asset")
		return
	}
	(a)["@key"] = key

	// Filter, validate and convert props to proper format
	err = a.ValidateProps()
	if err != nil {
		err = eh.WrapError(err, "format error")
		return
	}

	return
}

// CheckWriters checks if tx creator is allowed to write asset
func (a Asset) CheckWriters(stub shim.ChaincodeStubInterface) error {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return eh.NewCCError(400, fmt.Sprintf("asset type named %s does not exist", a.TypeTag()))
	}

	// Get tx creator MSP ID
	txCreator, err := cid.GetMSPID(stub)
	if err != nil {
		return eh.WrapErrorWithStatus(err, "error getting tx creator", 500)
	}

	// Check full asset write permission
	if assetTypeDef.Writers != nil {
		writePermission := false
		for _, w := range assetTypeDef.Writers {
			match, err := regexp.MatchString(w, txCreator)
			if err != nil {
				return eh.NewCCError(500, "failed to check if writer matches regexp")
			}
			if match {
				writePermission = true
			}
		}
		if !writePermission {
			return eh.NewCCError(403, fmt.Sprintf("%s cannot write to this asset", txCreator))
		}
	}

	// Check attributes write permission
	for _, prop := range assetTypeDef.Props {
		if _, exists := a[prop.Tag]; exists && prop.Writers != nil {
			writePermission := false
			for _, w := range prop.Writers {
				match, err := regexp.MatchString(w, txCreator)
				if err != nil {
					return eh.NewCCError(500, "failed to check if writer matches regexp")
				}
				if match {
					writePermission = true
				}
			}
			if !writePermission {
				return eh.NewCCError(403, fmt.Sprintf("%s cannot write to this asset property", txCreator))
			}
		}
	}

	return nil
}

// Put inserts asset in blockchain
func (a *Asset) Put(stub shim.ChaincodeStubInterface) (map[string]interface{}, error) {
	// Write index of references this asset points to
	err := a.PutRefs(stub)
	if err != nil {
		return nil, eh.WrapError(err, "failed writing reference index")
	}

	// Marshal asset back to JSON format
	assetJSON, err := json.Marshal(a)
	if err != nil {
		return nil, eh.WrapError(err, "failed to encode asset to JSON format")
	}

	// Write asset to blockchain
	if a.IsPrivate() {
		err = stub.PutPrivateData(a.TypeTag(), a.Key(), assetJSON)
		if err != nil {
			return nil, eh.WrapError(err, "failed to write asset to ledger")
		}
		assetKeyOnly := map[string]interface{}{
			"@key":       a.Key(),
			"@assetType": a.TypeTag(),
		}
		return assetKeyOnly, nil
	}

	err = stub.PutState(a.Key(), assetJSON)
	if err != nil {
		return nil, eh.WrapError(err, "failed to write asset to ledger")
	}
	return *a, nil
}

// Get reads asset from ledger
func (a *Asset) Get(stub shim.ChaincodeStubInterface) (*Asset, eh.ICCError) {
	var assetBytes []byte
	var err error
	if a.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(a.TypeTag(), a.Key())
	} else {
		assetBytes, err = stub.GetState(a.Key())
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

// ExistsInLedger checks if asset already exists
func (a *Asset) ExistsInLedger(stub shim.ChaincodeStubInterface) (bool, eh.ICCError) {
	var assetBytes []byte
	var err error
	if a.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(a.TypeTag(), a.Key())
	} else {
		assetBytes, err = stub.GetState(a.Key())
	}
	if err != nil {
		return false, eh.WrapErrorWithStatus(err, "unable to check asset existence", 400)
	}
	if assetBytes != nil {
		return true, nil
	}

	return false, nil
}

// Update receives a map[string]interface{} with key/vals to update in asset
func (a *Asset) Update(stub shim.ChaincodeStubInterface, update map[string]interface{}) error {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return eh.NewCCError(400, fmt.Sprintf("asset type named %s does not exist", a.TypeTag()))
	}

	// Get tx creator MSP ID
	txCreator, err := cid.GetMSPID(stub)
	if err != nil {
		return eh.WrapErrorWithStatus(err, "error getting tx creator", 500)
	}

	// Check full asset write permission
	if assetTypeDef.Writers != nil {
		writePermission := false
		for _, w := range assetTypeDef.Writers {
			match, err := regexp.MatchString(w, txCreator)
			if err != nil {
				return eh.NewCCError(500, "failed to check if writer matches regexp")
			}
			if match {
				writePermission = true
			}
		}
		if !writePermission {
			return eh.NewCCError(403, fmt.Sprintf("%s cannot write to this asset", txCreator))
		}
	}

	// Validate new asset properties
	for _, prop := range assetTypeDef.Props {
		// If prop is key, it cannot be updated
		if prop.IsKey {
			continue
		}

		// Check if property is included in the update map
		propInterface, propIncluded := update[prop.Tag]
		if !propIncluded {
			continue
		}

		if prop.ReadOnly {
			return eh.NewCCError(403, fmt.Sprintf("cannot update asset property %s", prop.Label))
		}

		// Check if tx creator is allowed to update this attribute
		if prop.Writers != nil {
			writePermission := false
			for _, w := range prop.Writers {
				match, err := regexp.MatchString(w, txCreator)
				if err != nil {
					return eh.NewCCError(500, "failed to check if writer matches regexp")
				}
				if match {
					writePermission = true
				}
			}
			if !writePermission {
				return eh.NewCCError(403, fmt.Sprintf("%s cannot write to this asset property", txCreator))
			}
		}

		// Validate data types
		propInterface, err := validateProp(propInterface, prop)
		if err != nil {
			return eh.WrapError(err, "error validating asset property")
		}

		(*a)[prop.Tag] = propInterface
	}
	return nil
}

// Delete erases asset from world state
func (a *Asset) Delete(stub shim.ChaincodeStubInterface) ([]byte, error) {
	isReferenced, err := a.IsReferenced(stub)
	if err != nil {
		return nil, eh.WrapError(err, "failed to check if asset if being referenced")
	}
	if isReferenced {
		return nil, eh.NewCCError(400, "another asset holds a reference to this one")
	}

	err = a.DelRefs(stub)
	if err != nil {
		return nil, eh.WrapError(err, "failed cleaning reference index")
	}

	var assetJSON []byte
	if !a.IsPrivate() {
		err = stub.DelState(a.Key())
		if err != nil {
			return nil, eh.WrapError(err, "failed to delete state from ledger")
		}
		assetJSON, err = json.Marshal(a)
		if err != nil {
			return nil, eh.WrapError(err, "failed to marshal asset")
		}
	} else {
		err = stub.DelPrivateData(a.TypeTag(), a.Key())
		if err != nil {
			return nil, eh.WrapError(err, "failed to delete state from private collection")
		}
		assetKeyOnly := map[string]interface{}{
			"@key":       a.Key(),
			"@assetType": a.TypeTag(),
		}
		assetJSON, err = json.Marshal(assetKeyOnly)
		if err != nil {
			return nil, eh.WrapError(err, "failed to marshal private asset key")
		}
	}

	return assetJSON, nil
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

// Refs returns all subAsset reference keys
func (a Asset) Refs(stub shim.ChaincodeStubInterface) ([]Key, eh.ICCError) {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return nil, eh.NewCCError(400, fmt.Sprintf("asset type named %s does not exist", a.TypeTag()))
	}
	assetSubAssets := assetTypeDef.SubAssets()
	var keys []Key
	for _, subAsset := range assetSubAssets {
		subAssetRefInterface, subAssetIncluded := a[subAsset.Tag]
		if !subAssetIncluded {
			// If subAsset is not included, no need to append
			continue
		}

		var isArray bool
		subAssetDataType := subAsset.DataType
		if strings.HasPrefix(subAssetDataType, "[]") {
			subAssetDataType = strings.TrimPrefix(subAssetDataType, "[]")
			isArray = true
		}

		// Handle array-like sub-asset property types
		var subAssetAsArray []interface{}
		if !isArray {
			subAssetAsArray = []interface{}{subAssetRefInterface}
		} else {
			var ok bool
			subAssetAsArray, ok = subAssetRefInterface.([]interface{})
			if !ok {
				return nil, eh.NewCCError(400, fmt.Sprintf("asset property %s must and array of type %s", subAsset.Label, subAsset.DataType))
			}
		}

		for _, subAssetRefInterface := range subAssetAsArray {
			subAssetRefMap, ok := subAssetRefInterface.(map[string]interface{})
			if !ok {
				// If subAsset is badly formatted, this method shouldn't have been called
				return nil, eh.NewCCError(400, "sub-asset reference badly formatted")
			}
			subAssetRefMap["@assetType"] = subAsset.DataType

			// Generate key for subAsset
			key, err := NewKey(subAssetRefMap)
			if err != nil {
				return nil, eh.WrapError(err, "failed to generate unique identifier for asset")
			}

			// Append key to response
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// ValidateProps checks if all props are compliant to format
func (a Asset) ValidateProps() eh.ICCError {
	// Perform validation of the @assetType field
	assetType, exists := a["@assetType"]
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
	for _, prop := range assetTypeDef.Props {
		// Check if required property is included
		propInterface, propIncluded := a[prop.Tag]
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

		a[prop.Tag] = propInterface
	}

	for propTag := range a {
		if strings.HasPrefix(propTag, "@") {
			continue
		}
		if !assetTypeDef.HasProp(propTag) {
			return eh.NewCCError(400, fmt.Sprintf("property %s is not defined in type %s", propTag, assetTypeString))
		}
	}

	return nil
}

// ValidateRefs checks if subAsset refs exists in blockchain
func (a Asset) ValidateRefs(stub shim.ChaincodeStubInterface) eh.ICCError {
	// Fetch references contained in asset
	refKeys, err := a.Refs(stub)
	if err != nil {
		return eh.WrapError(err, "failed to fetch references")
	}

	// Check if references exist
	for _, referencedKey := range refKeys {
		// Check if asset exists in blockchain
		assetJSON, err := referencedKey.Get(stub)
		if err != nil {
			return eh.WrapErrorWithStatus(err, "failed to read asset from blockchain", 400)
		}
		if assetJSON == nil {
			return eh.NewCCError(404, "referenced asset not found")
		}
	}
	return nil
}

// DelRefs deletes all the reference index for this asset from blockchain
func (a Asset) DelRefs(stub shim.ChaincodeStubInterface) error {
	// Fetch references contained in asset
	refKeys, err := a.Refs(stub)
	if err != nil {
		return eh.WrapErrorWithStatus(err, "failed to fetch references", 400)
	}

	assetKey := a.Key()

	// Delete reference indexes
	for _, referencedKey := range refKeys {
		// Construct reference key
		indexKey, err := stub.CreateCompositeKey(referencedKey.Key(), []string{assetKey})
		// Check if asset exists in blockchain
		err = stub.DelState(indexKey)
		if err != nil {
			return eh.WrapErrorWithStatus(err, "failed to read asset from blockchain", 400)
		}
	}

	return nil
}

// PutRefs writes to the blockchain the references
func (a Asset) PutRefs(stub shim.ChaincodeStubInterface) error {
	// Fetch references contained in asset
	refKeys, err := a.Refs(stub)
	if err != nil {
		return fmt.Errorf("failed to fetch references: %s", err)
	}

	assetKey := a.Key()

	// Delete reference indexes
	for _, referencedKey := range refKeys {
		// Construct reference key
		refKey, err := stub.CreateCompositeKey(referencedKey.Key(), []string{assetKey})
		if err != nil {
			return fmt.Errorf("failed generating composite key for reference: %s", err)
		}
		err = stub.PutState(refKey, []byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to put sub asset reference into blockchain: %s", err)
		}
	}

	return nil
}

// IsReferenced checks if asset is referenced by other asset
func (a Asset) IsReferenced(stub shim.ChaincodeStubInterface) (bool, error) {
	// Get asset key
	assetKey := a.Key()
	queryIt, err := stub.GetStateByPartialCompositeKey(assetKey, []string{})
	if err != nil {
		return false, fmt.Errorf("failed to check reference index: %s", err)
	}
	defer queryIt.Close()

	if queryIt.HasNext() {
		return true, nil
	}
	return false, nil
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
	assetKey := a["@key"].(string)
	return assetKey
}

// Type return the AssetType object for the asset
func (a Asset) Type() *AssetType {
	// Fetch asset properties
	assetTypeTag := a.TypeTag()
	assetTypeDef := FetchAssetType(assetTypeTag)
	return assetTypeDef
}
