package assets

// clean cleans an asset, erasing the data from the map
func (a *Asset) clean() {
	for k, v := range *a {
		if v == nil {
			delete(*a, k)
		}
	}
}
