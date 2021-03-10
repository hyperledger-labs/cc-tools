package assets

//AssetProp describes properties of each asset attribute
// The tag is used to reference the prop
// The label is for frontend rendering
// The description is a simple explanation for the specific field
// IsKey defines if the prop is a Primary Key
// Readonly makes impossible to overwrite a property
// The DefaulValue is the assets default state
// DataType defines the asset's type. It can receive custom data types
// Writers is an array of orgs that specify who can write in the asset
// Validate receives a function to validate the asset input
type AssetProp struct {
	Tag         string `json:"tag"`
	Label       string `json:"label"`
	Description string `json:"description"`

	IsKey    bool `json:"isKey"`
	Required bool `json:"required"`
	ReadOnly bool `json:"readOnly"`

	DefaultValue interface{} `json:"defaultValue,omitempty"`

	/* DataType can assume the following values:
	Primary types: "string", "number", "integer", "boolean", "datetime"
	Special types:
		-><assetType>: the specific asset type key (reference) as defined by <assetType> in the assets packages
		[]<type>: an array of elements specified by <type> as any of the above valid types
	*/
	DataType string `json:"dataType"`

	// Writers accepts either []string{'org1MSP', 'org2MSP'} and cc-tools will
	// check for an exact match  or []string{`$org\dMSP`} and cc-tools will
	// check for a match with regular expression `org\dMSP`
	Writers  []string                `json:"writers"`
	Validate func(interface{}) error `json:"-"`
}
