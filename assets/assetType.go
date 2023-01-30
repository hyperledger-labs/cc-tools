package assets

import "strings"

// AssetType is a list of all asset properties
type AssetType struct {
	// Tag is how the asset type will be referenced in the "@assetType" metaproperty.
	Tag string `json:"tag"`

	// Label is the pretty asset type name for front-end rendering
	Label string `json:"label"`

	// Description is a simple explanation describing the meaning of the asset type.
	Description string `json:"description"`

	// Props receives an array of assetProps, defining the asset's properties.
	Props []AssetProp `json:"props"`

	// Readers is an array that specifies which organizations can read the asset.
	// Must be coherent with private data collections configuration.
	// Accepts either basic strings for exact matches
	// eg. []string{'org1MSP', 'org2MSP'}
	// or regular expressions
	// eg. []string{`$org\dMSP`} and cc-tools will
	// check for a match with regular expression `org\dMSP`
	Readers []string `json:"readers,omitempty"`

	// Validate is a function called when validating asset as a whole.
	Validate func(Asset) error `json:"-"`

	// Dynamic is a flag that indicates if the asset type is dynamic.
	Dynamic bool `json:"dynamic,omitempty"`
}

// Keys returns a list of asset properties which are defined as primary keys. (IsKey == true)
func (t AssetType) Keys() (keys []AssetProp) {
	for _, prop := range t.Props {
		if prop.IsKey {
			keys = append(keys, prop)
		}
	}
	return
}

// SubAssets returns a list of asset properties which are subAssets (DataType is `->someAssetType`)
func (t AssetType) SubAssets() (subAssets []AssetProp) {
	for _, prop := range t.Props {
		dataType := prop.DataType
		dataType = strings.TrimPrefix(dataType, "[]")
		dataType = strings.TrimPrefix(dataType, "->")
		subAssetType := FetchAssetType(dataType)
		if subAssetType != nil {
			subAssets = append(subAssets, prop)
		}
	}
	return
}

// HasProp returns true if asset type has a property with the given tag.
func (t AssetType) HasProp(propTag string) bool {
	for _, prop := range t.Props {
		if prop.Tag == propTag {
			return true
		}
	}
	return false
}

// GetPropDef fetches the propDef with tag propTag.
func (t AssetType) GetPropDef(propTag string) *AssetProp {
	for _, prop := range t.Props {
		if prop.Tag == propTag {
			return &prop
		}
	}
	return nil
}

// IsPrivate returns true if asset is in a private collection.
func (t AssetType) IsPrivate() bool {
	return len(t.Readers) > 0
}

// ToMap returns a map representation of the asset type.
func (t AssetType) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"tag":         t.Tag,
		"label":       t.Label,
		"description": t.Description,
		"props":       ArrayFromAssetPropList(t.Props),
		"readers":     t.Readers,
		"dynamic":     t.Dynamic,
	}
}

// AssetTypeFromMap returns an asset type from a map representation.
func AssetTypeFromMap(m map[string]interface{}) AssetType {
	label, ok := m["label"].(string)
	if !ok {
		label = ""
	}
	description, ok := m["description"].(string)
	if !ok {
		description = ""
	}
	dynamic, ok := m["dynamic"].(bool)
	if !ok {
		dynamic = false
	}

	res := AssetType{
		Tag:         m["tag"].(string),
		Label:       label,
		Description: description,
		Props:       AssetPropListFromArray(m["props"].([]interface{})),
		Dynamic:     dynamic,
	}

	readers := make([]string, 0)
	readersArr, ok := m["readers"].([]interface{})
	if ok {
		for _, r := range readersArr {
			readers = append(readers, r.(string))
		}
	}
	if len(readers) > 0 {
		res.Readers = readers
	}

	return res
}

// ArrayFromAssetTypeList converts an array of AssetType to an array of map[string]interface
func ArrayFromAssetTypeList(assetTypes []AssetType) (array []map[string]interface{}) {
	for _, assetType := range assetTypes {
		array = append(array, assetType.ToMap())
	}
	return
}

// AssetTypeListFromArray converts an array of map[string]interface to an array of AssetType
func AssetTypeListFromArray(array []interface{}) (assetTypes []AssetType) {
	for _, v := range array {
		assetTypes = append(assetTypes, AssetTypeFromMap(v.(map[string]interface{})))
	}
	return
}
