# GoLedger CC Tools

[![Go Report Card](https://goreportcard.com/badge/github.com/hyperledger-labs/cc-tools)](https://goreportcard.com/report/github.com/hyperledger-labs/cc-tools)
[![GoDoc](https://godoc.org/github.com/hyperledger-labs/cc-tools?status.svg)](https://godoc.org/github.com/hyperledger-labs/cc-tools)

This project is a GoLedger open-source project aimed at providing tools for Hyperledger Fabric chaincode development in Golang. This might have breaking changes before we arrive at release v1.0.0. 

## Getting Started

Make sure you visit the repository [hyperledger-labs/cc-tools-demo](https://github.com/hyperledger-labs/cc-tools-demo), which is a template of a functional chaincode that uses cc-tools and provides ready-to-use scripts to deploy development networks. This is our preferred way of working, but you can feel free to import the package and assemble the chaincode as you choose. 

CC Tools has been tested with Hyperledger Fabric 1.x and 2.x realeases.

## Features
- Standard asset data mapping (and their properties)
- Encapsulation of Hyperledger Fabric chaincode sdk interface functions
- Standard asset key management
- Basic types of asset properties (text, number, boolean, date) available
- Basic asset array type (text, number or date arrays) available
- New asset property types customization
- Asset within assets available as references
- Asset array available as references
- Management of asset details
- Write permissions by set of organizations for each asset's property
- Private data collections management by asset (read permissions)
- Create/Read/Update/Delete (CRUD) transactions embedded
- Custom transactions, with prior definition of arguments, webservice method (GET, POST etc)
- Management of transaction details
- Compatible web service

## Assets
In the chaincode, assets represent the information that will be used by the blockchain ledger

* Tag
* Label
* Description
* Props
* Readers
* Validate
* Dynamic

An asset has a set o properties (Props) to structure the asset
A property has the following fields:

* Tag: string field used to define the name of the property referenced internally by the code and by the Rest API endpoints. No space or special characters allowed.
* Label: string field to define the label to be used by external applications. Free text.
* Description: property description string field to be used by external applications. Free text.
* IsKey: identifies if the property is part of the asset's keys. Boolean field.
* Required: identifies if the property is mandatory. Boolean field.
* ReadOnly: identifies if the property can no longer be modified once created. Boolean field.
* DefaultValue: property default value.
* Writers: define the organizations that can create or change this property. If the property is key (isKey field: true) then the entire asset can only be created by the organization. List of strings.
* DataType: property type. CC-Tools has the following default types: string, number, datetime and boolean. Custom types can be defined in the chaincode/datatypes folder. Arrays and references to other asset types are also possible.
* Validate: property validation function. It is suggested only for simple validations, more complex functions should use custom datatypes.

The asset package implements validation for the information on the asset type defined on one's chaincode. All the defined asset types must be on the assetTypeList, where CC-tools will validate them.

Asset package implements funtions that can be used in a custom transactions to manage assets of the ledger, such as `NewAsset`, `PutNew`, `Delete`, `Update`, `ExistsInLedger`, `GenerateKey`, `GetMap`, `Search` and others.

Also has the `StartupCheck`, which is a function that verifies if all the assets definitions are correctly implemented in the chaincode.

### Dynamic Asset Types
The package also allows for dynamically definition and creation of assets during runtime, through invokes for the chaincode.

### Datatypes
The assets package contain pre defined data types (string, number, datetime and boolean) and also supports arrays and references to other assets.
CC-tools supports the injection of custom data types in a chaincode, with its own validations and rules with the following fields:

* AcceptedFormats: list of "core" types that can be accepted (string, number, integer, boolean, datetime).
* Description: text describing the data type.
* DropDownValues: set of predetermined values to be used in a dropdown menu on frontend rendering.
* Parse: function called to check if the input value is valid, make necessary conversions and return a string representation of the value.


## Transactions
In the chaincode, transactions represent the GoLang methods that can modify the assets within the Blockchain ledger

CC-tools has a range of pre-defined transactions, defined in the transactions package:
* CreateAsset - creation of a new asset
* UpdateAsset - update an existing asset
* DeleteAsset - removing an asset from the current ledger state
* ReadAsset - read the asset in its last ledger status
* ReadAssetHistory - history of an asset's ledger status
* Search - asset listing

It is possible to also create custom transactions in the chaincode, that must be added to the txList so CC-tools can validate them

The transactions package has the `StartupCheck` just like the assets package, which is a function that verifies if all the transactions are correctly implemented in the chaincode.

### StubWrapper
The main purpose of the StubWrapper is to provide additional functionalities and simplify the development of chaincodes. 

The StubWrapper maintains a WriteSet to ensure that modifications made during the execution of a chaincode are properly reflected when querying the ledger state. Even if these changes have not been confirmed on the ledger yet, the StubWrapper records the pending modifications in the WriteSet. This allows subsequent queries to utilize the WriteSet to return the updated data, ensuring consistency and accuracy of information during the execution of the chaincode. The same applies to private data.

## Events
Hyperledger Fabric allows client applications (such as the rest-server API) to receive block events while block are commited to the peer's ledger.

CC-Tools have a built-in event funcionality that allows for a quick register of the event, which will automatically be listened and handled by the standard CCAPI.

The events in CC-Tools can be of three types: log, registering the information sent in the event payload on the CCAPI logs; transaction, which invokes another chaincode transactions when the event is received; and custom, which executes a custom function previously defined in the event.

## Mock
MockStub is an implementation of ChaincodeStubInterface for unit testing chaincode.

## Contributing
Feel free to fork it, create issues and PRs. We'll be happy to review them.

## Join our community

If you want to chat about Fabric, cc-tools and blockchain, you can reach GoLedger's technical team at our [Discord](https://discord.com/invite/GndkYHxNyQ)!