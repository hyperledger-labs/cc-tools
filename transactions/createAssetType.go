package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// CreateAssetType is the transaction which creates a dynamic Asset Type
var CreateAssetType = Transaction{
	Tag:         "createAssetType",
	Label:       "Create Asset Type",
	Description: "",
	Method:      "POST",
	Callers:     assets.GetAssetAdminsDynamicAssetType(),

	MetaTx: true,
	Args: ArgList{
		{
			Tag:         "assetTypes",
			Description: "Asset Types to be created.",
			DataType:    "[]@object",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		assetTypes := req["assetTypes"].([]interface{})
		list := make([]assets.AssetType, 0)

		for _, assetType := range assetTypes {
			assetTypeMap := assetType.(map[string]interface{})

			newAssetType, err := BuildAssetType(assetTypeMap)
			if err != nil {
				return nil, errors.WrapError(err, "failed to build asset type")
			}

			// Verify Asset Type existance
			assetTypeCheck := assets.FetchAssetType(newAssetType.Tag)
			if assetTypeCheck == nil {
				list = append(list, newAssetType)
			}
		}

		assets.UpdateAssetList(list)

		resBytes, err := json.Marshal(list)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}

func BuildAssetType(typeMap map[string]interface{}) (assets.AssetType, errors.ICCError) {
	// *********** Build Props Array ***********
	propsArr, ok := typeMap["props"].([]interface{})
	if !ok {
		return assets.AssetType{}, errors.NewCCError("invalid props array", http.StatusBadRequest)
	}
	props := make([]assets.AssetProp, len(propsArr))
	for i, prop := range propsArr {
		propMap := prop.(map[string]interface{})
		assetProp, err := BuildAssetProp(propMap)
		if err != nil {
			return assets.AssetType{}, errors.WrapError(err, "failed to build asset prop")
		}
		props[i] = assetProp
	}

	// *********** Check Type Values ***********
	// Tag
	tagValue, err := CheckValue(typeMap["tag"], true, "string", "tag")
	if err != nil {
		return assets.AssetType{}, errors.WrapError(err, "invalid tag value")
	}

	// Label
	labelValue, err := CheckValue(typeMap["label"], true, "string", "label")
	if err != nil {
		return assets.AssetType{}, errors.WrapError(err, "invalid label value")
	}

	// Description
	descriptionValue, err := CheckValue(typeMap["description"], false, "string", "description")
	if err != nil {
		return assets.AssetType{}, errors.WrapError(err, "invalid description value")
	}

	assetType := assets.AssetType{
		Tag:         tagValue.(string),
		Label:       labelValue.(string),
		Description: descriptionValue.(string),
		Props:       props,
		// Validate
	}

	// *********** Build Readers Array ***********
	readers := make([]string, 0)
	readersArr, ok := typeMap["readers"].([]interface{})
	if ok {
		for _, reader := range readersArr {
			readerValue, err := CheckValue(reader, false, "string", "reader")
			if err != nil {
				return assets.AssetType{}, errors.WrapError(err, "invalid reader value")
			}

			readers = append(readers, readerValue.(string))
		}
	}

	if len(readers) > 0 {
		assetType.Readers = readers
	}

	return assetType, nil
}

func BuildAssetProp(propMap map[string]interface{}) (assets.AssetProp, errors.ICCError) {
	// *********** Check Prop Values ***********
	// Tag
	tagValue, err := CheckValue(propMap["tag"], true, "string", "tag")
	if err != nil {
		return assets.AssetProp{}, errors.WrapError(err, "invalid tag value")
	}

	// Label
	labelValue, err := CheckValue(propMap["label"], true, "string", "label")
	if err != nil {
		return assets.AssetProp{}, errors.WrapError(err, "invalid label value")
	}

	// Description
	descriptionValue, err := CheckValue(propMap["description"], false, "string", "description")
	if err != nil {
		return assets.AssetProp{}, errors.WrapError(err, "invalid description value")
	}

	// Required
	requiredValue, err := CheckValue(propMap["required"], false, "boolean", "required")
	if err != nil {
		return assets.AssetProp{}, errors.WrapError(err, "invalid required value")
	}

	// IsKey
	isKeyValue, err := CheckValue(propMap["isKey"], false, "boolean", "isKey")
	if err != nil {
		return assets.AssetProp{}, errors.WrapError(err, "invalid isKey value")
	}

	// ReadOnly
	readOnlyValue, err := CheckValue(propMap["readOnly"], false, "boolean", "readOnly")
	if err != nil {
		return assets.AssetProp{}, errors.WrapError(err, "invalid readOnly value")
	}

	// DataType
	dataTypeValue, err := CheckValue(propMap["dataType"], true, "string", "dataType")
	if err != nil {
		return assets.AssetProp{}, errors.WrapError(err, "invalid dataType value")
	}

	err = CheckDataType(dataTypeValue.(string))
	if err != nil {
		return assets.AssetProp{}, errors.WrapError(err, "failed checking data type")
	}

	assetProp := assets.AssetProp{
		Tag:         tagValue.(string),
		Label:       labelValue.(string),
		Description: descriptionValue.(string),
		Required:    requiredValue.(bool),
		IsKey:       isKeyValue.(bool),
		ReadOnly:    readOnlyValue.(bool),
		DataType:    dataTypeValue.(string),
		// Validate
	}

	// *********** Build Writers Array ***********
	writers := make([]string, 0)
	writersArr, ok := propMap["writers"].([]interface{})
	if ok {
		for _, writer := range writersArr {
			writerValue, err := CheckValue(writer, false, "string", "writer")
			if err != nil {
				return assets.AssetProp{}, errors.WrapError(err, "invalid writer value")
			}

			writers = append(writers, writerValue.(string))
		}
	}
	if len(writers) > 0 {
		assetProp.Writers = writers
	}

	// ********* Validate Default Value *********
	if propMap["defaultValue"] != nil {
		// TODO: Reorganize utils functions and return ValidateProp to protected (validateProp)
		defaultValue, err := assets.ValidateProp(propMap["defaultValue"], assetProp)
		if err != nil {
			return assets.AssetProp{}, errors.WrapError(err, "invalid Default Value")
		}

		assetProp.DefaultValue = defaultValue
	}

	return assetProp, nil
}

func CheckDataType(dataType string) errors.ICCError {
	trimDataType := strings.TrimPrefix(dataType, "[]")

	if strings.HasPrefix(trimDataType, "->") {
		trimDataType = strings.TrimPrefix(trimDataType, "->")

		assetType := assets.FetchAssetType(trimDataType)
		if assetType == nil {
			return errors.NewCCError(fmt.Sprintf("invalid dataType value %s", dataType), http.StatusBadRequest)
		}
	} else {
		dataTypeObj := assets.FetchDataType(trimDataType)
		if dataTypeObj == nil {
			return errors.NewCCError(fmt.Sprintf("invalid dataType value %s", dataType), http.StatusBadRequest)
		}
	}

	return nil
}

func CheckValue(value interface{}, required bool, expectedType, fieldName string) (interface{}, errors.ICCError) {
	if value == nil {
		if required {
			return nil, errors.NewCCError(fmt.Sprintf("required value %s missing", fieldName), http.StatusBadRequest)
		}
		switch expectedType {
		case "string":
			return "", nil
		case "number":
			return 0, nil
		case "boolean":
			return false, nil
		}
	}

	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return nil, errors.NewCCError(fmt.Sprintf("value %s is not a string", fieldName), http.StatusBadRequest)
		}
	case "number":
		if _, ok := value.(float64); !ok {
			return nil, errors.NewCCError(fmt.Sprintf("value %s is not a number", fieldName), http.StatusBadRequest)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return nil, errors.NewCCError(fmt.Sprintf("value %s is not a boolean", fieldName), http.StatusBadRequest)
		}
	}

	return value, nil
}
