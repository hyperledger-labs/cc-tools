package test

import (
	"fmt"
	"strings"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// testCC implements the shim.Chaincode interface
type testCC struct{}

var testTxList = []tx.Transaction{
	tx.CreateAsset,
	tx.UpdateAsset,
	tx.DeleteAsset,
	tx.CreateAssetType,
	tx.UpdateAssetType,
	tx.DeleteAssetType,
	tx.LoadAssetTypeList,
}

var testAssetList = []assets.AssetType{
	{
		Tag:         "person",
		Label:       "Person",
		Description: "Personal data of someone",

		Props: []assets.AssetProp{
			{
				// Primary key
				Required: true,
				IsKey:    true,
				Tag:      "id",
				Label:    "CPF (Brazilian ID)",
				DataType: "cpf", // Datatypes are identified at datatypes folder
				Writers:  []string{`org1MSP`},
			},
			{
				// Mandatory property
				Required: true,
				Tag:      "name",
				Label:    "Name of the person",
				DataType: "string",
				// Validate funcion
				Validate: func(name interface{}) error {
					nameStr := name.(string)
					if nameStr == "" {
						return fmt.Errorf("name must be non-empty")
					}
					return nil
				},
			},
			{
				// Optional property
				Tag:      "dateOfBirth",
				Label:    "Date of Birth",
				DataType: "datetime",
				Writers:  []string{`org1MSP`},
			},
			{
				// Property with default value
				Tag:          "height",
				Label:        "Person's height",
				DefaultValue: 0,
				DataType:     "number",
			},
			{
				// Generic JSON object
				Tag:      "info",
				Label:    "Other Info",
				DataType: "@object",
			},
		},
	},
	{
		Tag:         "library",
		Label:       "Library",
		Description: "Library as a collection of books",

		Props: []assets.AssetProp{
			{
				// Primary Key
				Required: true,
				IsKey:    true,
				Tag:      "name",
				Label:    "Library Name",
				DataType: "string",
				Writers:  []string{`org3MSP`}, // This means only org3 can create the asset (others can edit)
			},
			{
				// Asset reference list
				Tag:      "books",
				Label:    "Book Collection",
				DataType: "[]->book",
			},
			{
				// Asset reference list
				Tag:      "entranceCode",
				Label:    "Entrance Code for the Library",
				DataType: "->secret",
			},
		},
	},
	{
		Tag:         "book",
		Label:       "Book",
		Description: "Book",

		Props: []assets.AssetProp{
			{
				// Composite Key
				Required: true,
				IsKey:    true,
				Tag:      "title",
				Label:    "Book Title",
				DataType: "string",
				Writers:  []string{`$org\dMSP`},
			},
			{
				// Composite Key
				Required: true,
				IsKey:    true,
				Tag:      "author",
				Label:    "Book Author",
				DataType: "string",
				Writers:  []string{`$org\dMSP`},
			},
			{
				/// Reference to another asset
				Tag:      "currentTenant",
				Label:    "Current Tenant",
				DataType: "->person",
			},
			{
				// String list
				Tag:      "genres",
				Label:    "Genres",
				DataType: "[]string",
			},
			{
				// Date property
				Tag:      "published",
				Label:    "Publishment Date",
				DataType: "datetime",
			},
		},
	},
	{
		Tag:         "secret",
		Label:       "Secret",
		Description: "Secret between Org2 and Org3",

		Readers: []string{"org2MSP", "org3MSP"},
		Props: []assets.AssetProp{
			{
				// Primary Key
				IsKey:    true,
				Tag:      "secretName",
				Label:    "Secret Name",
				DataType: "string",
				Writers:  []string{`org2MSP`}, // This means only org2 can create the asset (org3 can edit)
			},
			{
				// Mandatory Property
				Required: true,
				Tag:      "secret",
				Label:    "Secret",
				DataType: "string",
			},
		},
	},
	{
		Tag:         "assetTypeListData",
		Label:       "AssetTypeListData",
		Description: "AssetTypeListData",

		Props: []assets.AssetProp{
			{
				Required: true,
				IsKey:    true,
				Tag:      "id",
				Label:    "ID",
				DataType: "string",
				Writers:  []string{`org1MSP`},
			},
			{
				Required: true,
				Tag:      "list",
				Label:    "List",
				DataType: "[]@object",
				Writers:  []string{`org1MSP`},
			},
			{
				Required: true,
				Tag:      "lastUpdated",
				Label:    "Last Updated",
				DataType: "datetime",
				Writers:  []string{`org1MSP`},
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
	err := tx.StartupCheck()
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

	result, err := tx.Run(stub)

	if err != nil {
		response = err.GetErrorResponse()
		return
	}
	response = shim.Success([]byte(result))
	return
}
