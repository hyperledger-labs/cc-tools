package assets

//AssetProp describes properties of each asset attribute
type AssetProp struct {
	// The tag is used to reference the prop
	Tag string `json:"tag"`

	// The label is for frontend rendering
	Label string `json:"label"`

	// The description is a simple explanation for the specific field
	Description string `json:"description"`

	// IsKey defines if the prop is a Primary Key
	IsKey bool `json:"isKey"`

	// Tells if the prop is required
	Required bool `json:"required"`

	// Readonly makes impossible to overwrite a property
	ReadOnly bool `json:"readOnly"`

	// The DefaulValue is the assets default state
	DefaultValue interface{} `json:"defaultValue,omitempty"`

	/* DataType can assume the following values:
	Primary types: "string", "number", "integer", "boolean", "datetime"
	Special types:
		-><assetType>: the specific asset type key (reference) as defined by <assetType> in the assets packages
		[]<type>: an array of elements specified by <type> as any of the above valid types
	*/
	DataType string `json:"dataType"`

	// Writers is an array of orgs that specify who can write in the asset
	// and accepts either basic strings for exact matches
	// eg. []string{'org1MSP', 'org2MSP'}
	// or regular expressions
	// eg. []string{`$org\dMSP`} and cc-tools will
	// check for a match with regular expression `org\dMSP`
	Writers []string `json:"writers"`

	// Validate receives a function to validate the asset input
	Validate func(interface{}) error `json:"-"`
}
