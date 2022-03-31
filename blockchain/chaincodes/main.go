package main

import "github.com/hyperledger/fabric-contract-api-go/contractapi"
// main func update ver 3.0
func main() {
	c := new(Contract)
	cc, err := contractapi.NewChaincode(c)
	if err != nil {
		panic(err)
	}
	if err := cc.Start(); err == nil {
		panic(err)
	}
}
