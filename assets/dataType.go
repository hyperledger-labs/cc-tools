package assets

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/goledgerdev/cc-tools/errors"
)

// DataType is the interface required
type DataType struct {
	String   func(interface{}) (string, error)
	Validate func(interface{}) (interface{}, error)
}

var dataTypeMap = map[string]DataType{
	"string": {
		String: func(data interface{}) (string, error) {
			dataVal, ok := data.(string)
			if !ok {
				return "", errors.NewCCError("asset property should be a string", 400)
			}

			return dataVal, nil
		},

		Validate: func(data interface{}) (interface{}, error) {
			return nil, nil
		},
	},
	"number": {
		String: func(data interface{}) (string, error) {
			dataVal, ok := data.(float64)
			if !ok {
				propValStr, okStr := data.(string)
				if !okStr {
					return "", errors.NewCCError("asset property should be a number", 400)
				}
				var err error
				dataVal, err = strconv.ParseFloat(propValStr, 64)
				if err != nil {
					return "", errors.WrapErrorWithStatus(err, fmt.Sprintf("asset property should be a number"), 400)
				}
			}

			// Float IEEE 754 hexadecimal representation
			return strconv.FormatUint(math.Float64bits(dataVal), 16), nil
		},

		Validate: func(data interface{}) (interface{}, error) {
			return nil, nil
		},
	},
	"boolean": {
		String: func(data interface{}) (string, error) {
			dataVal, ok := data.(bool)
			if !ok {
				dataValStr, okStr := data.(string)
				if !okStr {
					return "", errors.NewCCError("asset property should be a boolean", 400)
				}
				if dataValStr != "true" && dataValStr != "false" {
					return "", errors.NewCCError("asset property should be a boolean", 400)
				}
				if dataValStr == "true" {
					dataVal = true
				}
			}

			if dataVal {
				return "t", nil
			}
			return "f", nil
		},

		Validate: func(data interface{}) (interface{}, error) {
			return nil, nil
		},
	},
	"datetime": {
		String: func(data interface{}) (string, error) {
			dataVal, ok := data.(string)
			if !ok {
				return "", errors.NewCCError("asset property should be a RFC3339 string", 400)
			}
			dataTime, err := time.Parse(time.RFC3339, dataVal)
			if err != nil {
				return "", errors.WrapErrorWithStatus(err, "invalid asset property RFC3339 format", 400)
			}

			return dataTime.Format(time.RFC3339), nil
		},

		Validate: func(data interface{}) (interface{}, error) {
			return nil, nil
		},
	},
}
