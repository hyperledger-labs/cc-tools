package transactions

type txArgument struct {
	Label       string `json:"label"`
	Description string `json:"description"`

	/* DataType can assume the following values:
	Primary types: "string", "number", "boolean", "datetime"
	Special types:
		@asset: any asset type defined in the assets package
		@key: key properties for any asset type defined in the assets package
		@update: update request for any asset type defined in the assets package
		@query: query string according to CouchDB specification: https://docs.couchdb.org/en/2.2.0/api/database/find.html
		<assetType>: the specific asset type as defined by <assetType> in the assets packages
		[]<type>: an array of elements specified by <type> as any of the above valid types
	*/
	DataType string `json:"dataType"`

	Required bool `json:"required"`
	Private  bool `json:"private"`
}
