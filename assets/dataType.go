package assets

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/goledgerdev/cc-tools/errors"
)

// DataType is the struct defining a primitive data type
type DataType struct {
	String   func(interface{}) (string, error)
	Validate func(interface{}) (interface{}, error)
}

// CustomDataTypes allows cc developer to inject custom primitive data types
func CustomDataTypes(m map[string]DataType) error {
	for k, v := range m {
		if v.String == nil || v.Validate == nil {
			return errors.NewCCError("invalid custom data type", 500)
		}
		dataTypeMap[k] = v
	}
	return nil
}

// DataTypeMap returns a copy of the primitive data type map
func DataTypeMap() map[string]DataType {
	ret := map[string]DataType{}
	for k, v := range dataTypeMap {
		ret[k] = v
	}
	return ret
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
			parsedData, ok := data.(string)
			if !ok {
				return nil, errors.NewCCError("property must be a string", 400)
			}
			return parsedData, nil
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
			parsedData, ok := data.(float64)
			if !ok {
				return nil, errors.NewCCError("property must be a number", 400)
			}

			return parsedData, nil
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
			parsedData, ok := data.(bool)
			if !ok {
				return nil, errors.NewCCError("property must be a boolean", 400)
			}

			return parsedData, nil
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
			dataVal, ok := data.(string)
			if !ok {
				return nil, errors.NewCCError("asset property must be an RFC3339 string", 400)
			}
			parsedData, err := time.Parse(time.RFC3339, dataVal)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "invalid asset property RFC3339 format", 400)
			}

			return parsedData, nil
		},
	},
}
