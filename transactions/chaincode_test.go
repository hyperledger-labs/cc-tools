package transactions

import (
	"fmt"
	"strings"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// testCC implements the shim.Chaincode interface
type testCC struct{}

var testTxList = []Transaction{}

var testAssetList = []assets.AssetType{
	{
		Tag:         "samplePerson",
		Label:       "Sample Person",
		Description: "",

		Props: []assets.AssetProp{
			{
				Tag:      "cpf",
				Label:    "CPF",
				DataType: "cpf",
				Writers:  []string{"org2MSP"},
			},
			{
				Tag:      "name",
				Label:    "Asset Name",
				Required: true,
				IsKey:    true,
				DataType: "string",
				Validate: func(name interface{}) error {
					nameStr := name.(string)
					if nameStr == "" {
						return fmt.Errorf("name must be non-empty")
					}
					return nil
				},
			},
			{
				Tag:          "readerScore",
				Label:        "Reader Score",
				DefaultValue: 0.0,
				DataType:     "number",
				Writers:      []string{`$org\dMSP`},
			},
			{
				Tag:      "secrets",
				Label:    "Secrets",
				DataType: "[]->sampleSecret",
			},
			{
				Tag:          "active",
				Label:        "Active",
				DefaultValue: false,
				DataType:     "boolean",
			},
		},
	},
	{
		Tag:         "author",
		Label:       "Author",
		Description: "",

		Props: []assets.AssetProp{
			{
				Tag:      "person",
				Label:    "Person",
				IsKey:    true,
				DataType: "->samplePerson",
			},
		},
	},
	{
		Tag:         "sampleBook",
		Label:       "Sample Book",
		Description: "",

		Props: []assets.AssetProp{
			{
				Tag:      "title",
				Label:    "Book Title",
				Required: true,
				IsKey:    true,
				DataType: "string",
			},
			{
				Tag:      "author",
				Label:    "Book Author",
				Required: true,
				IsKey:    true,
				DataType: "string",
			},
			{
				Tag:      "currentTenant",
				Label:    "Current Tenant",
				DataType: "->samplePerson",
			},
			{
				Tag:      "genres",
				Label:    "Genres",
				DataType: "[]string",
			},
			{
				Tag:      "published",
				Label:    "Publishment Date",
				DataType: "datetime",
			},
			{
				Tag:      "secret",
				Label:    "Book Secret",
				DataType: "->sampleSecret",
			},
		},
	},
	{
		Tag:         "sampleSecret",
		Label:       "Sample Secret",
		Description: "",

		Readers: []string{"org1MSP"},
		Props: []assets.AssetProp{
			{
				Tag:      "secretName",
				Label:    "Secret Name",
				IsKey:    true,
				DataType: "string",
			},
			{
				Tag:      "secret",
				Label:    "Secret",
				Required: true,
				DataType: "string",
			},
		},
	},
}

var testCustomDataTypes = map[string]assets.DataType{
	"cpf": {
		Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
			cpf, ok := data.(string)
			if !ok {
				return "", nil, errors.NewCCError("property must be a string", 400)
			}

			cpf = strings.ReplaceAll(cpf, ".", "")
			cpf = strings.ReplaceAll(cpf, "-", "")

			if len(cpf) != 11 {
				return "", nil, errors.NewCCError("CPF must have 11 digits", 400)
			}

			return cpf, cpf, nil
		},
	},
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *testCC) Init(stub shim.ChaincodeStubInterface) (response pb.Response) {
	err := StartupCheck()
	if err != nil {
		response = err.GetErrorResponse()
		return
	}

	// Get the args from the transaction proposal
	args := stub.GetStringArgs()

	// Test if argument list is empty
	if len(args) != 1 {
		response = shim.Error("the Init method expects 1 argument")
		response.Status = 400
		return
	}

	// Test if argument is "init" or "upgrade". Fails otherwise.
	if args[0] != "init" && args[0] != "upgrade" {
		response = shim.Error("the argument should be init or upgrade (as sent by Node.js SDK)")
		response.Status = 400
		return
	}

	response = shim.Success(nil)
	return
}

// Invoke is called per transaction on the chaincode.
func (t *testCC) Invoke(stub shim.ChaincodeStubInterface) (response pb.Response) {
	var result []byte

	result, err := Run(stub)

	if err != nil {
		response = err.GetErrorResponse()
		return
	}
	response = shim.Success([]byte(result))
	return
}
