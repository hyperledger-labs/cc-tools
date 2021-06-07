package transactions

// Argument struct stores the transaction argument info
// describing this specific input
type Argument struct {
	// Tag is the key of the value on the input map
	Tag string `json:"tag"`

	// Label is the name used in frontend
	Label string `json:"label"`

	// Description is a simple explanation of the argument
	Description string `json:"description"`

	// DataType can assume the following values:
	// Primary types: "string", "number", "integer", "boolean", "datetime"
	// Special types:
	//	  @asset: any asset type defined in the assets package
	//	  @key: key properties for any asset type defined in the assets package
	//	  @update: update request for any asset type defined in the assets package
	//	  @query: query string according to CouchDB specification: https://docs.couchdb.org/en/2.2.0/api/database/find.html
	//	  @object: arbitrary object
	//	  ->assetType: the specific asset type as defined by <assetType> in the assets packages
	//    dataType: any specific data type format defined by the chaincode
	//	  []type: an array of elements specified by <type> as any of the above valid types
	//
	DataType string `json:"dataType"`

	// Tells if the argument is required
	Required bool `json:"required"`

	// Tells if the argument will be used for private data
	Private bool `json:"private"`
}
