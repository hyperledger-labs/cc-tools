# GoLedger CC Tools

[![Go Report Card](https://goreportcard.com/badge/github.com/goledgerdev/cc-tools)](https://goreportcard.com/report/github.com/goledgerdev/cc-tools)
[![GoDoc](https://godoc.org/github.com/goledgerdev/cc-tools?status.svg)](https://godoc.org/github.com/goledgerdev/cc-tools)

This project is a GoLedger open-source project aimed at providing tools for Hyperledger Fabric chaincode development in Golang. This might have breaking changes before we arrive at release v1.0.0. 

## Getting Started

Make sure you visit the repository [goledgerdev/cc-tools-demo](https://github.com/goledgerdev/cc-tools-demo), which is a template of a functional chaincode that uses cc-tools and provides ready-to-use scripts to deploy development networks. This is our preferred way of working, but you can feel free to import the package and assemble the chaincode as you choose. 

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

## Contributing
Feel free to fork it, create issues and PRs. We'll be happy to review them.

## Join our community

If you want to chat about Fabric, cc-tools and blockchain, you can reach GoLedger's technical team at our [Discord](https://discord.com/invite/GndkYHxNyQ)!