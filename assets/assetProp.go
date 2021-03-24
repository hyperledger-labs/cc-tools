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

	// Validate receives a function to be called when validating property format
	Validate func(interface{}) error `json:"-"`
}
