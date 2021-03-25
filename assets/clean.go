package assets

// clean erases nil data from the asset map
func (a *Asset) clean() {
	for k, v := range *a {
		if v == nil {
			delete(*a, k)
		}
	}
}
