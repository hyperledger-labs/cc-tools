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

// RemoveAssetType removes an asset type from an AssetType List and returns a copy of the new list
func RemoveAssetType(tag string, l []AssetType) []AssetType {
	for i, assetType := range l {
		if assetType.Tag == tag {
			l = append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// ReplaceAssetType replaces an asset type from an AssetType List with an updated version and returns a copy of the new list
// This function does not automatically update the assetTypeList variable
func ReplaceAssetType(assetType AssetType, l []AssetType) []AssetType {
	for i, v := range l {
		if v.Tag == assetType.Tag {
			l = append(append(l[:i], assetType), l[i+1:]...)
		}
	}
	return l
}
