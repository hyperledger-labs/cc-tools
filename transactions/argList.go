package transactions

// ArgList defines the type for argument list in Transaction object
type ArgList []Argument

// GetArgDef fetches arg definition for arg with given tag
func (l ArgList) GetArgDef(tag string) *Argument {
	for _, arg := range l {
		if arg.Tag == tag {
			return &arg
		}
	}
	return nil
}
