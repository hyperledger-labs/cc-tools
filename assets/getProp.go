package assets

// GetProp returns the prop value. It returns nil if it doesn't exist.
func (a Asset) GetProp(propTag string) interface{} {
	val, _ := a[propTag]
	return val
}

// GetProp returns the prop value. It returns nil if it doesn't exist.
func (k Key) GetProp(propTag string) interface{} {
	val, _ := k[propTag]
	return val
}
