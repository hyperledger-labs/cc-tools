package assets

// DynamicAssetType is the configuration for the Dynamic AssetTypes
type DynamicAssetType struct {
	// Enabled defines whether the Dynamic AssetTypes feature is active
	Enabled bool `json:"enabled"`

	// AssetAdmins is an array that specifies which organizations can operate the Dynamic AssetTyper feature.
	// Accepts either basic strings for exact matches
	// eg. []string{'org1MSP', 'org2MSP'}
	// or regular expressions
	// eg. []string{`$org\dMSP`} and cc-tools will
	// check for a match with regular expression `org\dMSP`
	AssetAdmins []string `json:"assetAdmins"`
}
