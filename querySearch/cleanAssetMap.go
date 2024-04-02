package querysearch

func GetDefaultRemoveTags() *[]string {
	return &[]string{"@lastTouchBy", "@lastTx", "@lastUpdated", "@txId", "@txID"}
}

func CleanAssetMapDefault(m map[string]interface{}) map[string]interface{} {
	tags := GetDefaultRemoveTags()
	return CleanAssetMap(m, tags)
}

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
