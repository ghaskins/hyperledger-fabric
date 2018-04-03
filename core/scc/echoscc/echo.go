/*
Copyright State Street Bank. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package echocc

import (

	"hyperledger/cci/appinit"
	"hyperledger/cci/org/hyperledger/chaincode/system/echo"
	"hyperledger/ccs"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"fmt"
)

// EchoSysCC example simple Chaincode implementation
type EchoSysCC struct {
}

func (echoscc *EchoSysCC) Init(stub shim.ChaincodeStubInterface, param *appinit.Init) error {

	return nil
}

// Query callback representing the query of a chaincode
func (t *EchoSysCC) EchoRequest(stub shim.ChaincodeStubInterface, param *echo.Payload) (*echo.Payload, error) {

	fmt.Printf("Echo: %s\n", param.Data)
	return param, nil
}
