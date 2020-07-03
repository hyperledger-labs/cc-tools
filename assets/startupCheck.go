package assets

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goledgerdev/cc-tools/errors"
)

// StartupCheck verifies if asset definitions are properly coded, panicking if they're not
func StartupCheck() errors.ICCError {
	// Checks if there are references to undefined types
	for _, assetType := range assetTypeList {
		tag := assetType.Tag
		if tag == "" {
			return errors.NewCCError("asset has empty tag", 500)
		}
		if assetType.Label == "" {
			return errors.NewCCError(fmt.Sprintf("asset with tag %s has no label", tag), 500)
		}
		for _, w := range assetType.Writers {
			_, err := regexp.Compile(w)
			if err != nil {
				return errors.NewCCError(fmt.Sprintf("invalid writer regular expression %s for asset %s: %s", w, tag, err), 500)
			}
		}
		hasKey := false
		for _, prop := range assetType.Props {
			dataType := prop.DataType
			if strings.HasPrefix(dataType, "[]") {
				dataType = strings.TrimPrefix(dataType, "[]")
			}
			switch dataType {
			case "string":
			case "number":
			case "boolean":
			case "datetime":
			default:
				// Checks if the prop's datatype exists on assetMap
				propTypeDef := FetchAssetType(dataType)
				if propTypeDef == nil {
					return errors.NewCCError(fmt.Sprintf("reference for undefined type '%s'", prop.DataType), 500)
				}
			}

			for _, w := range prop.Writers {
				_, err := regexp.Compile(w)
				if err != nil {
					return errors.NewCCError(
						fmt.Sprintf("invalid writer regular expression %s for property %s of asset %s: %s", w, prop.Label, tag, err),
						500,
					)
				}
			}

			if prop.IsKey {
				hasKey = true
			}
		}
		if !hasKey {
			return errors.NewCCError(fmt.Sprintf("asset '%s' has no key properties", tag), 500)
		}
	}
	return nil
}
