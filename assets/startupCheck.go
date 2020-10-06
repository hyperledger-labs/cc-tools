package assets

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goledgerdev/cc-tools/errors"
)

// StartupCheck verifies if asset definitions are properly coded, panicking if they're not
func StartupCheck() errors.ICCError {
	assetTagSet := map[string]struct{}{}
	assetLabelSet := map[string]struct{}{}
	for _, assetType := range assetTypeList {
		// Check if asset tag and label are empty
		tag := assetType.Tag
		if tag == "" {
			return errors.NewCCError("asset has empty tag", 500)
		}
		label := assetType.Label
		if label == "" {
			return errors.NewCCError(fmt.Sprintf("asset with tag '%s' has no label", tag), 500)
		}

		// Check if asset tag or label is duplicate
		if _, duplicate := assetTagSet[tag]; duplicate {
			return errors.NewCCError(fmt.Sprintf("duplicate asset tag '%s'", tag), 500)
		}
		assetTagSet[tag] = struct{}{}
		if _, duplicate := assetLabelSet[label]; duplicate {
			return errors.NewCCError(fmt.Sprintf("duplicate asset label '%s'", label), 500)
		}
		assetLabelSet[label] = struct{}{}

		propTagSet := map[string]struct{}{}
		propLabelSet := map[string]struct{}{}
		hasKey := false
		for _, prop := range assetType.Props {
			// Check if prop tag or label is empty
			tag := prop.Tag
			if tag == "" {
				return errors.NewCCError(fmt.Sprintf("asset '%s' prop has empty tag", assetType.Tag), 500)
			}
			label := prop.Label
			if label == "" {
				return errors.NewCCError(fmt.Sprintf("asset '%s' prop with tag '%s' has no label", assetType.Tag, tag), 500)
			}

			// Check if prop tag or label is duplicate
			if _, duplicate := propTagSet[tag]; duplicate {
				return errors.NewCCError(fmt.Sprintf("duplicate asset prop tag '%s' in asset type '%s'", tag, assetType.Tag), 500)
			}
			propTagSet[tag] = struct{}{}
			if _, duplicate := propLabelSet[label]; duplicate {
				return errors.NewCCError(fmt.Sprintf("duplicate asset prop label '%s' in asset type '%s'", label, assetType.Tag), 500)
			}
			propLabelSet[label] = struct{}{}

			dataTypeName := prop.DataType
			if strings.HasPrefix(dataTypeName, "[]") {
				dataTypeName = strings.TrimPrefix(dataTypeName, "[]")
			}

			// Check if there are references to undefined types
			// Check if prop's datatype is a primitive type
			dataType, dataTypeExists := dataTypeMap[dataTypeName]
			if !dataTypeExists {
				// Checks if the prop's datatype exists on assetMap
				propTypeDef := FetchAssetType(dataTypeName)
				if propTypeDef == nil {
					return errors.NewCCError(fmt.Sprintf("reference for undefined type '%s'", prop.DataType), 500)
				}
				if prop.DefaultValue != nil {
					return errors.NewCCError(fmt.Sprintf("default value cannot be used in reference in prop '%s' of asset '%s'", prop.Label, assetType.Label), 500)
				}
			} else {
				if prop.DefaultValue != nil {
					_, _, err := dataType.Parse(prop.DefaultValue)
					if err != nil {
						return errors.WrapErrorWithStatus(err, fmt.Sprintf("invalid default value in prop '%s' of asset '%s'", prop.Label, assetType.Label), 500)
					}
				}
			}

			// Check if writers in regex mode compile
			for _, w := range prop.Writers {
				if len(w) <= 1 {
					continue
				}
				if w[0] == '$' {
					_, err := regexp.Compile(w[1:])
					if err != nil {
						return errors.WrapErrorWithStatus(err, fmt.Sprintf("invalid writer regular expression %s for property %s of asset %s", w, prop.Label, tag), 500)
					}
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
