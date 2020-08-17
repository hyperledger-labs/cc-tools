package assets

import "strings"

// AssetType is a list of all asset properties
type AssetType struct {
	Tag         string `json:"tag"`
	Label       string `json:"label"`
	Description string `json:"description"`

	Props    []AssetProp       `json:"props"`
	Writers  []string          `json:"writers,omitempty"`
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

// IsPrivate returns true if asset is in a private collection
func (t AssetType) IsPrivate() bool {
	return len(t.Readers) > 0
}
