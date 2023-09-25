package assets

import (
	"fmt"

	"github.com/hyperledger-labs/cc-tools/errors"
)

// SetProp sets the prop value with proper validation. It does not update the asset in the ledger.
func (a *Asset) SetProp(propTag string, value interface{}) errors.ICCError {
	if len(propTag) == 0 {
		return errors.NewCCError("propTag cannot be empty", 500)
	}
	if propTag[0] == '@' {
		return errors.NewCCError("cannot modify internal properties", 500)
	}
	assetType := a.Type()
	if assetType == nil {
		return errors.NewCCError("asset type does not exist", 500)
	}
	propDef := assetType.GetPropDef(propTag)
	if propDef == nil {
		return errors.NewCCError(fmt.Sprintf("asset type '%s' does not have prop named '%s'", assetType.Tag, propTag), 500)
	}

	if propDef.IsKey {
		return errors.NewCCError("SetProp on key asset property is not yet implemented", 501) // TODO
	}

	propType := DataTypeMap()[propDef.DataType]

	_, parsedVal, err := propType.Parse(value)
	if err != nil {
		return errors.WrapError(err, fmt.Sprintf("invalid '%s' value", propTag))
	}
	(*a)[propTag] = parsedVal

	return nil
}
