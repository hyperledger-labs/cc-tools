package assets

import (
	"net/http"
	"time"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
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

// RemoveAssetType removes an asset type from an AssetType List and returns the new list
func RemoveAssetType(tag string, l []AssetType) []AssetType {
	for i, assetType := range l {
		if assetType.Tag == tag {
			l = append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// ReplaceAssetType replaces an asset type from an AssetType List with an updated version and returns the new list
// This function does not automatically update the assetTypeList variable
func ReplaceAssetType(assetType AssetType, l []AssetType) []AssetType {
	for i, v := range l {
		if v.Tag == assetType.Tag {
			l[i] = assetType
		}
	}
	return l
}

// StoreAssetList stores the current assetList on the blockchain
func StoreAssetList(stub *sw.StubWrapper) errors.ICCError {
	assetList := AssetTypeList()
	l := ArrayFromAssetTypeList(assetList)

	txTimestamp, err := stub.Stub.GetTxTimestamp()
	if err != nil {
		return errors.WrapError(err, "failed to get tx timestamp")
	}
	txTime := txTimestamp.AsTime()
	txTimeStr := txTime.Format(time.RFC3339)

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
			return errors.WrapError(err, "error getting asset list")
		}
		listMap := (map[string]interface{})(*listAsset)

		listMap["list"] = l
		listMap["lastUpdated"] = txTimeStr

		_, err = listAsset.Update(stub, listMap)
		if err != nil {
			return errors.WrapError(err, "error updating asset list")
		}
	} else {
		listMap := map[string]interface{}{
			"@assetType":  "assetTypeListData",
			"id":          "primary",
			"list":        l,
			"lastUpdated": txTimeStr,
		}

		listAsset, err := NewAsset(listMap)
		if err != nil {
			return errors.WrapError(err, "error creating asset list")
		}

		_, err = listAsset.PutNew(stub)
		if err != nil {
			return errors.WrapError(err, "error putting asset list")
		}
	}

	SetAssetListUpdateTime(txTime)

	return nil
}

// RestoreAssetList restores the assetList from the blockchain
func RestoreAssetList(stub *sw.StubWrapper, init bool) errors.ICCError {
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

		txTime := listMap["lastUpdated"].(time.Time)

		if GetAssetListUpdateTime().After(txTime) {
			return nil
		}

		l := AssetTypeListFromArray(listMap["list"].([]interface{}))

		l = getRestoredList(l, init)

		ReplaceAssetList(l)

		SetAssetListUpdateTime(txTime)
	}

	return nil
}

// getRestoredList reconstructs the assetTypeList from the stored list comparing to the current list
func getRestoredList(storedList []AssetType, init bool) []AssetType {
	assetList := AssetTypeList()

	deleteds := AssetTypeList()

	for _, assetTypeStored := range storedList {
		if !assetTypeStored.Dynamic {
			continue
		}

		exists := false
		for i, assetType := range assetList {
			if assetType.Tag == assetTypeStored.Tag {
				exists = true

				assetTypeStored.Validate = assetType.Validate
				assetList[i] = assetTypeStored

				deleteds = RemoveAssetType(assetType.Tag, deleteds)

				break
			}
		}
		if !exists {
			assetList = append(assetList, assetTypeStored)
		}
	}

	if !init {
		for _, deleted := range deleteds {
			if !deleted.Dynamic {
				continue
			}

			assetList = RemoveAssetType(deleted.Tag, assetList)
		}
	}

	return assetList
}

// GetAssetListUpdateTime returns the last time the asset list was updated
func GetAssetListUpdateTime() time.Time {
	return listUpdateTime
}

// SetAssetListUpdateTime sets the last time the asset list was updated
func SetAssetListUpdateTime(t time.Time) {
	listUpdateTime = t
}
