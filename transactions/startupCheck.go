package transactions

import (
	"fmt"
	"regexp"

	"github.com/goledgerdev/cc-tools/errors"
)

// StartupCheck verifies if asset definitions are properly coded, panicking if they're not
func StartupCheck() errors.ICCError {
	// Checks if there are references to undefined types
	for _, tx := range txList {
		txName := tx.Tag
		for _, w := range tx.Callers {
			_, err := regexp.Compile(w)
			if err != nil {
				return errors.NewCCError(fmt.Sprintf("invalid caller regular expression %s for tx %s: %s", w, txName, err), 500)
			}
		}
	}
	return nil
}
