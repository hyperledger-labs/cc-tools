package accesscontrol

import (
	"regexp"

	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func AllowCaller(stub shim.ChaincodeStubInterface, allowedCallers []Caller) (bool, error) {
	if allowedCallers == nil {
		return true, nil
	}

	callerMSP, err := cid.GetMSPID(stub)
	if err != nil {
		return false, errors.WrapError(err, "could not get MSP id")
	}

	var grantedPermission bool
	for i := 0; i < len(allowedCallers) && !grantedPermission; i++ {
		allowed := allowedCallers[i]
		isAllowedMSP, err := checkMSP(callerMSP, allowed.MSP)
		if err != nil {
			return false, errors.WrapError(err, "failed to check MSP")
		}

		isAllowedOU, err := checkOU(stub, allowed.OU)
		if err != nil {
			return false, errors.WrapError(err, "failed to check OU")
		}

		isAllowedAttrs, err := checkAttributes(stub, allowed.Attributes)
		if err != nil {
			return false, errors.WrapError(err, "failed to check attributes")
		}

		grantedPermission = isAllowedMSP && isAllowedOU && isAllowedAttrs
	}

	return grantedPermission, nil
}

func checkMSP(callerMsp, allowedMSP string) (bool, error) {
	if len(allowedMSP) <= 1 {
		return true, nil
	}

	// if caller is regexp
	if allowedMSP[0] == '$' {
		match, err := regexp.MatchString(allowedMSP[1:], callerMsp)
		if err != nil {
			return false, errors.NewCCError("failed to check if caller matches regexp", 500)
		}

		return match, nil
	}

	// if caller is not regexss
	return callerMsp == allowedMSP, nil
}

func checkOU(stub shim.ChaincodeStubInterface, allowedOU string) (bool, error) {
	if allowedOU == "" {
		return true, nil
	}

	return cid.HasOUValue(stub, allowedOU)
}

func checkAttributes(stub shim.ChaincodeStubInterface, allowedAttrs map[string]string) (bool, error) {
	if allowedAttrs == nil {
		return true, nil
	}

	for key, value := range allowedAttrs {
		callerValue, _, err := cid.GetAttributeValue(stub, key)
		if err != nil {
			return false, err
		}

		if callerValue != value {
			return false, nil
		}
	}

	return true, nil
}
