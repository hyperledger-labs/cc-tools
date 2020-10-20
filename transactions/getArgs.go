package transactions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GetArgs validates and merges tx arguments from args list and transient map
func (tx Transaction) GetArgs(stub shim.ChaincodeStubInterface) (map[string]interface{}, errors.ICCError) {
	var err error

	// Extract the function and args from the transaction proposal
	_, args := stub.GetFunctionAndParameters()

	reqMap := make(map[string]interface{})

	// Prepare request arguments
	var req map[string]interface{}
	var transientReq map[string]interface{}

	// Unmarshal public args
	if len(args) > 0 {
		err = json.Unmarshal([]byte(args[0]), &req)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, fmt.Sprintf("failed to unmarshal request args"), 400)
		}
	}

	// Unmarshal private args
	transientMap, _ := stub.GetTransient()
	var transientArgs []byte
	if transientMap != nil {
		var transientExists bool
		transientArgs, transientExists = transientMap["@request"]
		if transientExists && transientArgs != nil {
			err = json.Unmarshal(transientArgs, &transientReq)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, fmt.Sprintf("failed to unmarshal transient args"), 400)
			}
		}
	}

	// Validate args format
	for argKey, argDef := range tx.Args {
		var arg interface{}
		var argExists bool
		if argDef.Private {
			arg, argExists = transientReq[argKey]
		} else {
			// If argument is not private, it can either come as transient arg or public arg
			// TODO: In case it comes as both, we could try to merge them if they were an array
			arg, argExists = req[argKey]
			if !argExists {
				arg, argExists = transientReq[argKey]
			}
		}
		if !argExists {
			if argDef.Required {
				return nil, errors.NewCCError(fmt.Sprintf("missing argument '%s'", argKey), 400)
			}
			continue
		}

		argType := argDef.DataType
		var isArray bool
		if strings.HasPrefix(argType, "[]") {
			argType = strings.TrimPrefix(argType, "[]")
			isArray = true
		}

		if isArray {
			argAsSlice, ok := arg.([]interface{})
			if !ok {
				return nil, errors.NewCCError(fmt.Sprintf("argument '%s' must be an array", argKey), 400)
			}
			if argDef.Required && len(argAsSlice) == 0 {
				return nil, errors.NewCCError(fmt.Sprintf("required argument '%s' must be non-empty", argKey), 400)
			}

			for argIdx, arg := range argAsSlice {
				validArgElem, err := validateTxArg(argType, arg)
				if err != nil {
					return nil, errors.WrapError(err, fmt.Sprintf("invalid argument '%s'", argKey))
				}
				argAsSlice[argIdx] = validArgElem
			}
			reqMap[argKey] = argAsSlice
		} else {
			validArg, err := validateTxArg(argType, arg)
			if err != nil {
				return nil, errors.WrapError(err, fmt.Sprintf("invalid argument '%s'", argKey))
			}
			reqMap[argKey] = validArg
		}

	}

	return reqMap, nil
}

func validateTxArg(argType string, arg interface{}) (interface{}, errors.ICCError) {
	var argAsInterface interface{}
	var err error

	dataTypeMap := assets.DataTypeMap()
	dataType, dataTypeExists := dataTypeMap[argType]
	if dataTypeExists { // if argument is a primitive data type
		if !dataType.IsLegacy() {
			_, argAsInterface, err = dataType.Parse(arg)
		} else {
			argAsInterface, err = dataType.Validate(arg)
		}
		if err != nil {
			return nil, errors.WrapError(err, "invalid argument format")
		}
	} else {
		switch argType {
		case "@asset":
			var asset assets.Asset
			argBytes, err := json.Marshal(arg)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed marshaling arg", 400)
			}
			err = json.Unmarshal(argBytes, &asset)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed unmarshaling arg", 400)
			}
			argAsInterface = asset
		case "@key":
			var key assets.Key
			argBytes, err := json.Marshal(arg)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed marshaling arg", 400)

			}
			err = json.Unmarshal(argBytes, &key)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed unmarshaling arg", 400)
			}
			argAsInterface = key
		case "@update":
			var argMap map[string]interface{}
			argMap, ok := arg.(map[string]interface{})
			if !ok {
				return nil, errors.NewCCError("invalid argument format", 400)
			}
			_, ok = argMap["@assetType"]
			if !ok {
				return nil, errors.NewCCError("missing @assetType", 400)
			}
			argAsInterface = argMap
		case "@query":
			var argMap map[string]interface{}
			argMap, ok := arg.(map[string]interface{})
			if !ok {
				return nil, errors.NewCCError("invalid argument format", 400)
			}
			_, ok = argMap["selector"]
			if !ok {
				return nil, errors.NewCCError("missing selector", 400)
			}
			argAsInterface = argMap
		default:
			var asset assets.Asset
			argBytes, err := json.Marshal(arg)
			if err != nil {
				return nil, errors.NewCCError("failed to marshal arg", 400)
			}
			err = json.Unmarshal(argBytes, &asset)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed unmarshaling arg", 400)
			}
			if asset.TypeTag() != argType {
				return nil, errors.NewCCError(fmt.Sprintf("arg must be of type %s", argType), 400)
			}
			argAsInterface = asset
		}
	}

	return argAsInterface, nil
}
