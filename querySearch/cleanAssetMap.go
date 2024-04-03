// This code contains functions related to cleaning up an asset map by removing certain tags.
package querysearch

// List of default tags that will be removed if enabled
func GetDefaultRemoveTags() *[]string {
	return &[]string{"@lastTouchBy", "@lastTx", "@lastUpdated", "@txId", "@txID"}
}

// Remove fields in map by default tags
func CleanAssetMapDefault(m map[string]interface{}) map[string]interface{} {
	return CleanAssetMap(m, GetDefaultRemoveTags())
}

// This is a recursive function that cleans an asset map by removing tags.
// It takes a map and a slice of tags to remove.
// It iterates over the map and removes any tags that are in the removeTag slice.
// If a value in the map is another map or a slice,
// it recursively calls CleanAssetMap on that value.
func CleanAssetMap(m map[string]interface{}, removeTag *[]string) map[string]interface{} {
	if removeTag != nil {
		for _, tag := range *removeTag {
			delete(m, tag)
		}
	}

	for k, v := range m {
		switch prop := v.(type) {
		case map[string]interface{}:
			m[k] = CleanAssetMap(prop, removeTag)

		case []interface{}:
			for idx, elem := range prop {
				if elemMap, ok := elem.(map[string]interface{}); ok {
					prop[idx] = CleanAssetMap(elemMap, removeTag)
				}
			}
		}
	}

	return m
}

// Internal function for remove tags
func (q *QuerySearch) removeTags(m *map[string]interface{}) {
	if q.config.NoRemoveTagsTransaction && len(q.config.RemoveTags) == 0 {
		return
	}

	remove := []string{}
	if !q.config.NoRemoveTagsTransaction {
		remove = append(remove, q.config.removeDefaultTags...)
	}

	remove = append(remove, q.config.RemoveTags...)
	*m = CleanAssetMap(*m, &remove)
}
