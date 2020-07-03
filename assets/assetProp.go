package assets

/*
AssetProp describes properties of each asset attribute
*/
type AssetProp struct {
	Tag      string `json:"tag"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	IsKey    bool   `json:"isKey"`
	ReadOnly bool   `json:"readOnly"`

	/* DataType can assume the following values:
	Primary types: "string", "number", "boolean", "datetime"
	Special types:
		<assetType>: the specific asset type key (reference) as defined by <assetType> in the assets packages
		[]<type>: an array of elements specified by <type> as any of the above valid types
	*/
	DataType string                  `json:"dataType"`
	Writers  []string                `json:"writers"`
	Validate func(interface{}) error `json:"-"`
}
