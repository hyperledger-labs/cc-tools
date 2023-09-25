package transactions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
)

// StartupCheck verifies if tx definitions are properly coded, returning an error if they're not.
func StartupCheck() errors.ICCError {
	// Checks if there are references to undefined types
	for _, tx := range txList {
		txName := tx.Tag
		for _, c := range tx.Callers {
			if len(c) <= 1 {
				continue
			}
			if c[0] == '$' {
				_, err := regexp.Compile(c[1:])
				if err != nil {
					return errors.WrapErrorWithStatus(err, fmt.Sprintf("invalid caller regular expression %s for tx %s", c, txName), 500)
				}
			}
		}

		argSet := map[string]interface{}{}
		for _, arg := range tx.Args {
			if _, duplicate := argSet[arg.Tag]; duplicate {
				return errors.NewCCError(fmt.Sprintf("duplicate arg tag %s in tx %s", arg.Tag, txName), 500)
			}
			argSet[arg.Tag] = struct{}{}

			dtype := strings.TrimPrefix(arg.DataType, "[]")
			if dtype != "@asset" &&
				dtype != "@key" &&
				dtype != "@update" &&
				dtype != "@query" &&
				dtype != "@object" {
				if strings.HasPrefix(dtype, "->") {
					dtype = strings.TrimPrefix(dtype, "->")
					if assets.FetchAssetType(dtype) == nil {
						return errors.NewCCError(fmt.Sprintf("invalid arg type %s in tx %s", arg.DataType, txName), 500)
					}
				} else {
					if assets.FetchDataType(dtype) == nil {
						return errors.NewCCError(fmt.Sprintf("invalid arg type %s in tx %s", arg.DataType, txName), 500)
					}
				}
			}
		}
	}
	return nil
}
