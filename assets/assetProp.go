package assets

// AssetProp describes properties of each asset attribute
type AssetProp struct {
	// Tag is the string used to reference the prop in the Asset map
	Tag string `json:"tag"`

	// Label is the pretty property name for front-end rendering
	Label string `json:"label"`

	// Description is a simple explanation describing the meaning of the property.
	Description string `json:"description"`

	// IsKey defines if the property is a Primary Key. At least one of the Asset's
	// properties must have this set as true.
	IsKey bool `json:"isKey"`

	// Required defines if the NewAsset function should fail if property is undefined.
	// If IsKey is set as true, this value is ignored. (primary key is always required)
	Required bool `json:"required"`

	// ReadOnly means property can only be set during Asset creation.
	ReadOnly bool `json:"readOnly"`

	// DefaultValue is the default property value in case it is not defined when parsing asset.
	DefaultValue interface{} `json:"defaultValue,omitempty"`

	// DataType can assume the following values:
	// Primary types: "string", "number", "integer", "boolean", "datetime"
	// Special types:
	//   -><assetType>: the specific asset type key (reference) as defined by <assetType> in the assets packages
	//   []<type>: an array of elements specified by <type> as any of the above valid types
	DataType string `json:"dataType"`

	// Writers is an array of orgs that specify who can write in the asset
	// and accepts either basic strings for exact matches
	// eg. []string{'org1MSP', 'org2MSP'}
	// or regular expressions
	// eg. []string{`$org\dMSP`} and cc-tools will
	// check for a match with regular expression `org\dMSP`
	Writers []string `json:"writers"`

	// Validate is a function called when validating property format.
	Validate func(interface{}) error `json:"-"`
}

// ToMap converts an AssetProp to a map[string]interface{}
func (p AssetProp) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"tag":          p.Tag,
		"label":        p.Label,
		"description":  p.Description,
		"isKey":        p.IsKey,
		"required":     p.Required,
		"readOnly":     p.ReadOnly,
		"defaultValue": p.DefaultValue,
		"dataType":     p.DataType,
		"writers":      p.Writers,
	}
}

// AssetPropFromMap converts a map[string]interface{} to an AssetProp
func AssetPropFromMap(m map[string]interface{}) AssetProp {
	description, ok := m["description"].(string)
	if !ok {
		description = ""
	}
	label, ok := m["label"].(string)
	if !ok {
		label = ""
	}
	isKey, ok := m["isKey"].(bool)
	if !ok {
		isKey = false
	}
	required, ok := m["required"].(bool)
	if !ok {
		required = false
	}
	readOnly, ok := m["readOnly"].(bool)
	if !ok {
		readOnly = false
	}

	res := AssetProp{
		Tag:          m["tag"].(string),
		Label:        label,
		Description:  description,
		IsKey:        isKey,
		Required:     required,
		ReadOnly:     readOnly,
		DefaultValue: m["defaultValue"],
		DataType:     m["dataType"].(string),
	}

	writers := make([]string, 0)
	writersArr, ok := m["writers"].([]interface{})
	if ok {
		for _, w := range writersArr {
			writers = append(writers, w.(string))
		}
	}
	if len(writers) > 0 {
		res.Writers = writers
	}

	return res
}

// ArrayFromAssetPropList converts an array of AssetProp to an array of map[string]interface
func ArrayFromAssetPropList(a []AssetProp) []map[string]interface{} {
	list := []map[string]interface{}{}
	for _, m := range a {
		list = append(list, m.ToMap())
	}
	return list
}

// AssetPropListFromArray converts an array of map[string]interface to an array of AssetProp
func AssetPropListFromArray(a []interface{}) []AssetProp {
	list := []AssetProp{}
	for _, m := range a {
		list = append(list, AssetPropFromMap(m.(map[string]interface{})))
	}
	return list
}
