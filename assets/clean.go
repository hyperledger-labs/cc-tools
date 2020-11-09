package assets

func (a *Asset) clean() {
	for k, v := range *a {
		if v == nil {
			delete(*a, k)
		}
	}
}
