package assets

import (
	"encoding/json"
	"net/http"

	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

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
	if GetEnabledDynamicAssetType() {
		l = append(l, GetListAssetType())
	}
	assetTypeList = l
}

// ReplaceAssetList replace assetTypeList to for a new one
func ReplaceAssetList(l []AssetType) {
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

// StoreAssetList stores the current assetList on the blockchain
func StoreAssetList(stub *sw.StubWrapper) errors.ICCError {
	assetList := AssetTypeList()
	l := ArrayFromAssetTypeList(assetList)

	listKey, err := NewKey(map[string]interface{}{
		"@assetType": "assetTypeListData",
		"id":         "primary",
	})
	if err != nil {
		return errors.NewCCError("error getting asset list key", http.StatusInternalServerError)
	}

	exists, err := listKey.ExistsInLedger(stub)
	if err != nil {
		return errors.NewCCError("error checking if asset list exists", http.StatusInternalServerError)
	}

	if exists {
		listAsset, err := listKey.Get(stub)
		if err != nil {
			return errors.NewCCError("error getting asset list", http.StatusInternalServerError)
		}
		listMap := (map[string]interface{})(*listAsset)

		listMap["list"] = l

		_, err = listAsset.Update(stub, listMap)
		if err != nil {
			return errors.NewCCError("error updating asset list", http.StatusInternalServerError)
		}
	} else {
		listMap := map[string]interface{}{
			"@assetType": "assetTypeListData",
			"id":         "primary",
			"list":       l,
		}

		listAsset, err := NewAsset(listMap)
		if err != nil {
			return errors.NewCCError("error creating asset list", http.StatusInternalServerError)
		}

		_, err = listAsset.PutNew(stub)
		if err != nil {
			return errors.NewCCError("error putting asset list", http.StatusInternalServerError)
		}
	}

	return nil
}

// RestoreAssetList restores the assetList from the blockchain
func RestoreAssetList(stub *sw.StubWrapper) errors.ICCError {
	listKey, err := NewKey(map[string]interface{}{
		"@assetType": "assetTypeListData",
		"id":         "primary",
	})
	if err != nil {
		return errors.NewCCError("error gettin asset list key", http.StatusInternalServerError)
	}

	exists, err := listKey.ExistsInLedger(stub)
	if err != nil {
		return errors.NewCCError("error checking if asset list exists", http.StatusInternalServerError)
	}

	if exists {
		listAsset, err := listKey.Get(stub)
		if err != nil {
			return errors.NewCCError("error getting asset list", http.StatusInternalServerError)
		}
		listMap := (map[string]interface{})(*listAsset)

		l := AssetTypeListFromArray(listMap["list"].([]map[string]interface{}))

		ReplaceAssetList(l)
	}

	return nil
}

func SetEventForList(stub *sw.StubWrapper) errors.ICCError {
	list := AssetTypeList()
	listJson, err := json.Marshal(list)
	if err != nil {
		return errors.NewCCError("error marshaling asset list", http.StatusInternalServerError)
	}

	err = stub.Stub.SetEvent("assetListChange", listJson)
	if err != nil {
		return errors.NewCCError("error setting event for asset list", http.StatusInternalServerError)
	}

	return nil
}
