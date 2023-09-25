package assets

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/hyperledger-labs/cc-tools/errors"
)

// DataType is the struct defining a primitive data type.
type DataType struct {
	// AcceptedFormats is a list of "core" types that can be accepted (string, number, integer, boolean, datetime)
	AcceptedFormats []string `json:"acceptedFormats"`

	// Description is a simple text describing the data type
	Description string `json:"description,omitempty"`

	// DropDownValues is a set of predetermined values to be used in a dropdown menu on frontend rendering
	DropDownValues map[string]interface{} `json:"DropDownValues"`

	// Parse is called to check if the input value is valid, make necessary
	// conversions and returns a string representation of the value
	Parse func(interface{}) (string, interface{}, errors.ICCError) `json:"-"`
}

// CustomDataTypes allows cc developer to inject custom primitive data types
func CustomDataTypes(m map[string]DataType) error {
	for k, v := range m {
		if v.Parse == nil {
			return errors.NewCCError(fmt.Sprintf("invalid custom data type '%s': nil Parse function", k), 500)
		}

		dataType := v
		dataTypeMap[k] = &dataType
	}
	return nil
}

// DataTypeMap returns a copy of the primitive data type map
func DataTypeMap() map[string]DataType {
	ret := map[string]DataType{}
	for k, v := range dataTypeMap {
		ret[k] = *v
	}
	return ret
}

// FetchDataType returns a pointer to the DataType object or nil if asset type is not found.
func FetchDataType(dataTypeTag string) *DataType {
	return dataTypeMap[dataTypeTag]
}

// dataTypeMap contains the "standard" primitive data types
var dataTypeMap = map[string]*DataType{
	"string": {
		AcceptedFormats: []string{"string"},
		Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
			parsedData, ok := data.(string)
			if !ok {
				return parsedData, nil, errors.NewCCError("property must be a string", 400)
			}
			return parsedData, parsedData, nil
		},
	},
	"number": {
		AcceptedFormats: []string{"number"},
		Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
			var dataVal float64
			switch v := data.(type) {
			case float64:
				dataVal = v
			case int:
				dataVal = (float64)(v)
			case string:
				var err error
				dataVal, err = strconv.ParseFloat(v, 64)
				if err != nil {
					return "", nil, errors.WrapErrorWithStatus(err, "asset property must be a number", 400)
				}
			default:
				return "", nil, errors.NewCCError("asset property must be a number", 400)
			}

			// Float IEEE 754 hexadecimal representation
			return strconv.FormatUint(math.Float64bits(dataVal), 16), dataVal, nil
		},
	},
	"integer": {
		AcceptedFormats: []string{"number"},
		Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
			var dataVal float64
			switch v := data.(type) {
			case float64:
				dataVal = v
			case int:
				dataVal = (float64)(v)
			case string:
				var err error
				dataVal, err = strconv.ParseFloat(v, 64)
				if err != nil {
					return "", nil, errors.WrapErrorWithStatus(err, "asset property must be an integer", 400)
				}
			default:
				return "", nil, errors.NewCCError("asset property must be an integer", 400)
			}

			retVal := math.Trunc(dataVal)

			if dataVal != retVal {
				return "", nil, errors.NewCCError("asset property must be an integer", 400)
			}

			// Float IEEE 754 hexadecimal representation
			return fmt.Sprintf("%d", int64(retVal)), int64(retVal), nil
		},
	},
	"boolean": {
		AcceptedFormats: []string{"boolean"},
		Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
			dataVal, ok := data.(bool)
			if !ok {
				dataValStr, okStr := data.(string)
				if !okStr {
					return "", nil, errors.NewCCError("asset property must be a boolean", 400)
				}
				if dataValStr != "true" && dataValStr != "false" {
					return "", nil, errors.NewCCError("asset property must be a boolean", 400)
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
		Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
			dataTime, ok := data.(time.Time)
			if !ok {
				dataVal, ok := data.(string)
				if !ok {
					return "", nil, errors.NewCCError("asset property must be a RFC3339 string", 400)
				}
				var err error
				dataTime, err = time.Parse(time.RFC3339, dataVal)
				if err != nil {
					return "", nil, errors.WrapErrorWithStatus(err, "invalid asset property RFC3339 format", 400)
				}
			}

			return dataTime.Format(time.RFC3339), dataTime, nil
		},
	},
	"@object": {
		AcceptedFormats: []string{"@object"},
		Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
			dataVal, ok := data.(map[string]interface{})
			if !ok {
				switch v := data.(type) {
				case []byte:
					err := json.Unmarshal(v, &dataVal)
					if err != nil {
						return "", nil, errors.WrapErrorWithStatus(err, "failed to unmarshal []byte into map[string]interface{}", http.StatusBadRequest)
					}
				case string:
					err := json.Unmarshal([]byte(v), &dataVal)
					if err != nil {
						return "", nil, errors.WrapErrorWithStatus(err, "failed to unmarshal string into map[string]interface{}", http.StatusBadRequest)
					}
				default:
					return "", nil, errors.NewCCError(fmt.Sprintf("asset property must be either a byte array or a string, but received type is: %T", data), http.StatusBadRequest)
				}
			}

			dataVal["@assetType"] = "@object"

			retVal, err := json.Marshal(dataVal)
			if err != nil {
				return "", nil, errors.WrapErrorWithStatus(err, "failed to marshal return value", http.StatusInternalServerError)
			}

			return string(retVal), dataVal, nil
		},
	},
}
