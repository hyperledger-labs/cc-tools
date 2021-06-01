package assets

// GetProp returns the prop value. It returns nil if it doesn't exist.
func (a Asset) GetProp(propTag string) interface{} {
	return a[propTag]
}
