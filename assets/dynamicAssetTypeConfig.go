package assets

// dynamicAssetTypeConfig is the configuration data for the Dynamic assetTypes feature
var dynamicAssetTypeConfig = DynamicAssetType{}

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

// GetListAssetType returns the Dynamic AssetType meta type
func GetListAssetType() AssetType {
	var AssetTypeListData = AssetType{
		Tag:         "assetTypeListData",
		Label:       "AssetTypeListData",
		Description: "AssetTypeListData",

		Props: []AssetProp{
			{
				Required: true,
				IsKey:    true,
				Tag:      "id",
				Label:    "ID",
				DataType: "string",
				Writers:  dynamicAssetTypeConfig.AssetAdmins,
			},
			{
				Required: true,
				Tag:      "list",
				Label:    "List",
				DataType: "[]@object",
				Writers:  dynamicAssetTypeConfig.AssetAdmins,
			},
			{
				Required: true,
				Tag:      "lastUpdated",
				Label:    "Last Updated",
				DataType: "datetime",
				Writers:  dynamicAssetTypeConfig.AssetAdmins,
			},
		},
	}
	return AssetTypeListData
}
