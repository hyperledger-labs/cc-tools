package assets

// dynamicAssetTypeConfig is the configuration data for the Dynamic assetTypes feature
var dynamicAssetTypeConfig = DynamicAssetType{}

// DynamicAssetTypeConfig returns a copy of the DynamicAssetType variable.
func DynamicAssetTypeConfig() DynamicAssetType {
	configCopy := DynamicAssetType{
		Enabled:     dynamicAssetTypeConfig.Enabled,
		AssetAdmins: dynamicAssetTypeConfig.AssetAdmins,
	}
	return configCopy
}

// InitDynamicAssetTypeConfig initilizes the dynamicAssetTypeConfig variable
func InitDynamicAssetTypeConfig(c DynamicAssetType) {
	dynamicAssetTypeConfig = c
}

// GetEnabledDynamicAssetType returns the value of the Enabled field
func GetEnabledDynamicAssetType() bool {
	return dynamicAssetTypeConfig.Enabled
}

// GetAssetAdminsDynamicAssetType returns the value of the AssetAdmins field
func GetAssetAdminsDynamicAssetType() []string {
	return dynamicAssetTypeConfig.AssetAdmins
}
