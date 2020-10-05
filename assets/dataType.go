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
	// AcceptedFormats is a list of "core" types that can be accepted (string, number or boolean)
	AcceptedFormats []string `json:"acceptedFormats"`
	Description     string   `json:"description,omitempty"`

	Parse func(interface{}) (string, interface{}, error) `json:"-"`

	KeyString func(interface{}) (string, error)      `json:"-"` // DEPRECATED. Use Parse instead.
	Validate  func(interface{}) (interface{}, error) `json:"-"` // DEPRECATED. Use Parse instead.
}

// CustomDataTypes allows cc developer to inject custom primitive data types
func CustomDataTypes(m map[string]DataType) error {
	for k, v := range m {
		if v.Parse == nil {
			// These function signatures are deprecated and this is here for backwards compatibility only.
			if v.KeyString == nil || v.Validate == nil {
				return errors.NewCCError("invalid custom data type", 500)
			}
			v.Parse = func(d interface{}) (key string, ret interface{}, err error) {
				key, err = v.KeyString(d)
				if err != nil {
					return
				}

				ret, err = v.Validate(d)
				return
			}
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
		AcceptedFormats: []string{"string"},
		Parse: func(data interface{}) (string, interface{}, error) {
			parsedData, ok := data.(string)
			if !ok {
				return parsedData, nil, errors.NewCCError("property must be a string", 400)
			}
			return parsedData, parsedData, nil
		},
	},
	"number": {
		AcceptedFormats: []string{"number"},
		Parse: func(data interface{}) (string, interface{}, error) {
			dataVal, ok := data.(float64)
			if !ok {
				propValStr, okStr := data.(string)
				if !okStr {
					return "", nil, errors.NewCCError("asset property should be a number", 400)
				}
				var err error
				dataVal, err = strconv.ParseFloat(propValStr, 64)
				if err != nil {
					return "", nil, errors.WrapErrorWithStatus(err, fmt.Sprintf("asset property should be a number"), 400)
				}
			}

			// Float IEEE 754 hexadecimal representation
			return strconv.FormatUint(math.Float64bits(dataVal), 16), dataVal, nil
		},
	},
	"integer": {
		AcceptedFormats: []string{"number"},
		Parse: func(data interface{}) (string, interface{}, error) {
			dataVal, ok := data.(float64)
			if !ok {
				propValStr, okStr := data.(string)
				if !okStr {
					return "", nil, errors.NewCCError("asset property should be an integer", 400)
				}
				var err error
				dataVal, err = strconv.ParseFloat(propValStr, 64)
				if err != nil {
					return "", nil, errors.WrapErrorWithStatus(err, fmt.Sprintf("asset property should be an integer"), 400)
				}
			}

			retVal := math.Trunc(dataVal)

			if dataVal != retVal {
				return "", nil, errors.NewCCError("asset property should be an integer", 400)
			}

			// Float IEEE 754 hexadecimal representation
			return fmt.Sprintf("%d", int64(retVal)), int64(retVal), nil
		},
	},
	"boolean": {
		AcceptedFormats: []string{"boolean"},
		Parse: func(data interface{}) (string, interface{}, error) {
			dataVal, ok := data.(bool)
			if !ok {
				dataValStr, okStr := data.(string)
				if !okStr {
					return "", nil, errors.NewCCError("asset property should be a boolean", 400)
				}
				if dataValStr != "true" && dataValStr != "false" {
					return "", nil, errors.NewCCError("asset property should be a boolean", 400)
				}
				if dataValStr == "true" {
					dataVal = true
				}
			}

			if dataVal {
				return "t", dataVal, nil
			}
			return "f", dataVal, nil
		},
	},
	"datetime": {
		AcceptedFormats: []string{"string"},
		Parse: func(data interface{}) (string, interface{}, error) {
			dataVal, ok := data.(string)
			if !ok {
				return "", nil, errors.NewCCError("asset property should be a RFC3339 string", 400)
			}
			dataTime, err := time.Parse(time.RFC3339, dataVal)
			if err != nil {
				return "", nil, errors.WrapErrorWithStatus(err, "invalid asset property RFC3339 format", 400)
			}

			return dataTime.Format(time.RFC3339), dataTime, nil
		},
	},
}
