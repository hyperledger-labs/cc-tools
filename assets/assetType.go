package assets

import "strings"

// AssetType is a list of all asset properties
// The label is for frontend rendering
// The description is a simple explanation for the specific field
// Props receives an array of assetProps, definig the assets properties
// Readers is an array that specifies which organizations can read the asset (used for private data)
// Validates is a function that performs the asset input validation
type AssetType struct {
	Tag         string `json:"tag"`
	Label       string `json:"label"`
	Description string `json:"description"`

	Props    []AssetProp       `json:"props"`
	Readers  []string          `json:"readers,omitempty"`
	Validate func(Asset) error `json:"-"`
}

// Keys returns a list of asset properties which are defined as primary keys
func (t AssetType) Keys() (keys []AssetProp) {
	for _, prop := range t.Props {
		if prop.IsKey {
			keys = append(keys, prop)
		}
	}
	return
}

// SubAssets returns a list of asset properties which are subAssets
func (t AssetType) SubAssets() (subAssets []AssetProp) {
	for _, prop := range t.Props {
		dataType := prop.DataType
		if strings.HasPrefix(dataType, "[]") {
			dataType = strings.TrimPrefix(dataType, "[]")
		}
		if strings.HasPrefix(dataType, "->") {
			dataType = strings.TrimPrefix(dataType, "->")
		}
		subAssetType := FetchAssetType(dataType)
		if subAssetType != nil {
			subAssets = append(subAssets, prop)
		}
	}
	return
}

// HasProp returns true if asset type has a property with the given tag
func (t AssetType) HasProp(propTag string) bool {
	for _, prop := range t.Props {
		if prop.Tag == propTag {
			return true
		}
	}
	return false
}

// GetPropDef fetches the propDef with tag propTag
func (t AssetType) GetPropDef(propTag string) *AssetProp {
	for _, prop := range t.Props {
		if prop.Tag == propTag {
			return &prop
		}
	}
	return nil
}

// IsPrivate returns true if asset is in a private collection
func (t AssetType) IsPrivate() bool {
	return len(t.Readers) > 0
}
