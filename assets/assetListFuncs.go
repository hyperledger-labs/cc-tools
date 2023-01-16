package assets

// AssetTypeList returns a copy of the assetTypeList variable.
func AssetTypeList() []AssetType {
	listCopy := make([]AssetType, len(assetTypeList))
	copy(listCopy, assetTypeList)
	return listCopy
}

// FetchAssetType returns a pointer to the AssetType object or nil if asset type is not found.
func FetchAssetType(assetTypeTag string) *AssetType {
	for _, assetType := range assetTypeList {
		if assetType.Tag == assetTypeTag {
			return &assetType
		}
	}
	return nil
}

// InitAssetList appends custom assets to assetTypeList to avoid initialization loop.
func InitAssetList(l []AssetType) {
	assetTypeList = l
}

// UpdateAssetList updates the assetTypeList variable on runtime
func UpdateAssetList(l []AssetType) {
	assetTypeList = append(assetTypeList, l...)
}

// RemoveAssetType removes an asset type from an assetTypeList and returns a copy of the new list
func RemoveAssetType(tag string, assetTypeList []AssetType) []AssetType {
	for i, assetType := range assetTypeList {
		if assetType.Tag == tag {
			assetTypeList = append(assetTypeList[:i], assetTypeList[i+1:]...)
		}
	}
	return assetTypeList
}

// ReplaceAssetType replaces an asset type from an assetTypeList with an updated version and returns a copy of the new list
func ReplaceAssetType(assetType AssetType, assetTypeList []AssetType) []AssetType {
	for i, v := range assetTypeList {
		if v.Tag == assetType.Tag {
			assetTypeList = append(append(assetTypeList[:i], assetType), assetTypeList[i+1:]...)
		}
	}
	return assetTypeList
}
