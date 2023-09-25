package assets

import (
	"fmt"
	"strings"

	"github.com/hyperledger-labs/cc-tools/errors"
)

// ValidateProps checks if all props are compliant to format
func (a Asset) ValidateProps() errors.ICCError {
	// Perform validation of the @assetType field
	assetType, exists := a["@assetType"]
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
	for _, prop := range assetTypeDef.Props {
		// Check if required property is included
		propInterface, propIncluded := a[prop.Tag]
		if !propIncluded {
			if prop.DefaultValue == nil {
				if prop.Required {
					return errors.NewCCError(fmt.Sprintf("property %s (%s) is required", prop.Tag, prop.Label), 400)
				}
				if prop.IsKey {
					return errors.NewCCError(fmt.Sprintf("key property %s (%s) is required", prop.Tag, prop.Label), 400)
				}
				continue
			}
			propInterface = prop.DefaultValue
		}

		// Validate data types
		propInterface, err := validateProp(propInterface, prop)
		if err != nil {
			msg := fmt.Sprintf("error validating asset '%s' property", prop.Tag)
			return errors.WrapError(err, msg)
		}

		a[prop.Tag] = propInterface
	}

	for propTag := range a {
		if strings.HasPrefix(propTag, "@") {
			continue
		}
		if !assetTypeDef.HasProp(propTag) {
			return errors.NewCCError(fmt.Sprintf("property %s is not defined in type %s", propTag, assetTypeString), 400)
		}
	}

	return nil
}
