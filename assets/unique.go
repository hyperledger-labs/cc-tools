package assets

import (
	"fmt"
	"strings"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// UniqueProps returns a list of asset properties which are defined as unique values
func (t AssetType) UniqueProps() (uniqueProps []AssetProp) {
	for _, prop := range t.Props {
		if prop.Unique {
			uniqueProps = append(uniqueProps, prop)
		}
	}
	return
}

// calculateUniqueMarkers returns map with keys=propDef.Tag values=uniqueMarker
func (a *Asset) calculateUniqueMarkers(stub shim.ChaincodeStubInterface) (map[string]string, errors.ICCError) {
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return nil, errors.NewCCError(fmt.Sprintf("asset type named '%s' does not exist", a.TypeTag()), 400)
	}

	uniqueProps := assetTypeDef.UniqueProps()

	markers := map[string]string{}
	for _, propDef := range uniqueProps {
		markerSeed := assetTypeDef.Tag + ":" + propDef.Tag

		// Check if unique property is included
		propInterface, propIncluded := (*a)[propDef.Tag]
		if !propIncluded {
			continue
		}

		// The [] check is for possible future implementation of yet-undefined array-like uniqueness
		dataTypeName := propDef.DataType
		if strings.HasPrefix(dataTypeName, "[]") {
			dataTypeName = strings.TrimPrefix(dataTypeName, "[]")
		}

		dataType, dataTypeExists := dataTypeMap[dataTypeName]
		if dataTypeExists {
			// Generate marker for primitive data type prop
			seed, _, err := dataType.Parse(propInterface)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "error generating unique marker", 400)
			}
			markerSeed += ":" + seed
		} else {
			// If not a primitive type, check if type is defined in assetMap
			subAssetType := FetchAssetType(dataTypeName)
			if subAssetType == nil {
				return nil, errors.NewCCError(fmt.Sprintf("invalid data type named '%s'", dataTypeName), 400)
			}
			// Generate marker for subasset
			propAsMap, ok := propInterface.(map[string]interface{})
			if !ok {
				return nil, errors.NewCCError(fmt.Sprintf("failed to generate unique marker for prop '%s'", propDef.Label), 400)
			}

			key, err := NewKey(propAsMap)
			if err != nil {
				return nil, errors.WrapError(err, fmt.Sprintf("failed to generate unique marker for prop '%s'", propDef.Label))
			}

			markerSeed += ":" + key.Key()
		}

		marker := "unique:" + uuid.NewSHA1(uuid.NameSpaceOID, []byte(markerSeed)).String()

		markers[propDef.Tag] = marker
	}

	return markers, nil
}

func (a *Asset) putUniqueMarkers(stub shim.ChaincodeStubInterface) errors.ICCError {
	var err error

	markers, err := a.calculateUniqueMarkers(stub)
	if err != nil {
		return errors.WrapError(err, fmt.Sprintf("failed to generate unique markers for asset of type '%s'", a.TypeTag()))
	}

	for propTag, marker := range markers {
		err = stub.PutState(marker, []byte{0x00})
		if err != nil {
			return errors.WrapError(err, fmt.Sprintf("failed to put unique marker for prop '%s'", propTag))
		}
	}

	return nil
}

func (a *Asset) checkUniqueMarkers(stub shim.ChaincodeStubInterface) errors.ICCError {
	var err error

	markers, err := a.calculateUniqueMarkers(stub)
	if err != nil {
		return errors.WrapError(err, fmt.Sprintf("failed to generate unique markers for asset of type '%s'", a.TypeTag()))
	}

	for propTag, marker := range markers {
		data, err := stub.GetState(marker)
		if err != nil {
			return errors.WrapError(err, fmt.Sprintf("failed to check uniqueness for prop '%s'", propTag))
		}
		if data != nil {
			return errors.WrapError(err, fmt.Sprintf("prop '%s' must be unique", propTag))
		}
	}

	return nil
}

func (a *Asset) delUniqueMarkers(stub shim.ChaincodeStubInterface) errors.ICCError {
	var err error

	markers, err := a.calculateUniqueMarkers(stub)
	if err != nil {
		return errors.WrapError(err, fmt.Sprintf("failed to generate unique markers for asset of type '%s'", a.TypeTag()))
	}

	for propTag, marker := range markers {
		err := stub.DelState(marker)
		if err != nil {
			return errors.WrapError(err, fmt.Sprintf("failed to delete unique marker for prop '%s'", propTag))
		}
	}

	return nil
}
