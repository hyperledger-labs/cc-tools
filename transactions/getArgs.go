package transactions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// GetArgs validates the received arguments and assembles a map with the parsed key/values.
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
			return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal request args", 400)
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
				return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal transient args", 400)
			}
		}
	}

	cleanUp(req)
	cleanUp(transientReq)

	// Validate args format
	for _, argDef := range tx.Args {
		argKey := argDef.Tag
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

// Validates a given transaction argument
func validateTxArg(argType string, arg interface{}) (interface{}, errors.ICCError) {
	var argAsInterface interface{}
	var err error

	// If argType has "->" it means it's an asset reference
	if strings.HasPrefix(argType, "->") {
		argType = strings.TrimPrefix(argType, "->")

		// Handle @assetType
		var argMap map[string]interface{}
		argMap, ok := arg.(map[string]interface{})
		if !ok {
			return nil, errors.NewCCError("invalid argument format", 400)
		}
		assetTypeName, ok := argMap["@assetType"]
		if ok && assetTypeName != argType { // in case an @assetType is specified, check if it is correct
			return nil, errors.NewCCError(fmt.Sprintf("invalid @assetType '%s' (expecting '%s')", assetTypeName, argType), 400)
		}
		if !ok { // if @assetType is not specified, inject it
			argMap["@assetType"] = argType
		}
		key, err := assets.NewKey(argMap)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed constructing key", 400)
		}
		argAsInterface = key
	} else {
		switch argType {
		case "@asset":
			var argMap map[string]interface{}
			argMap, ok := arg.(map[string]interface{})
			if !ok {
				return nil, errors.NewCCError("invalid argument format", 400)
			}
			asset, err := assets.NewAsset(argMap)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed constructing asset", 400)
			}
			argAsInterface = asset
		case "@key":
			var argMap map[string]interface{}
			argMap, ok := arg.(map[string]interface{})
			if !ok {
				return nil, errors.NewCCError("invalid argument format", 400)
			}
			key, err := assets.NewKey(argMap)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed constructing key", 400)
			}
			argAsInterface = key
		case "@update":
			var argMap map[string]interface{}
			argMap, ok := arg.(map[string]interface{})
			if !ok {
				return nil, errors.NewCCError("invalid argument format", 400)
			}
			_, err := assets.NewKey(argMap)
			if err != nil {
				return nil, errors.WrapError(err, "argument of type '@update' must be a valid key")
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
		case "@object":
			var argMap map[string]interface{}
			argMap, ok := arg.(map[string]interface{})
			if !ok {
				return nil, errors.NewCCError("invalid argument format", 400)
			}
			argAsInterface = argMap
		default: // should be a specific datatype
			dataTypeMap := assets.DataTypeMap()
			dataType, dataTypeExists := dataTypeMap[argType]
			if !dataTypeExists {
				return nil, errors.NewCCError(fmt.Sprintf("invalid arg type '%s'", argType), 500)
			}

			_, argAsInterface, err = dataType.Parse(arg)

			if err != nil {
				return nil, errors.WrapError(err, "invalid argument format")
			}
		}
	}

	return argAsInterface, nil
}
