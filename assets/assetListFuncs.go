package assets

// AssetTypeList returns a copy of the assetTypeList variable
func AssetTypeList() []AssetType {
	listCopy := make([]AssetType, len(assetTypeList))
	copy(listCopy, assetTypeList)
	return listCopy
}

// FetchAssetType returns a pointer to the AssetType object or nil if asset type is not found
func FetchAssetType(assetTypeTag string) *AssetType {
	for _, assetType := range assetTypeList {
		if assetType.Tag == assetTypeTag {
			return &assetType
		}
	}
	return nil
}

// InitAssetList appends custom assets to assetTypeList to avoid initialization loop
func InitAssetList(l []AssetType) {
	assetTypeList = l
}
