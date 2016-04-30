/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package sample_syscc

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/system_chaincode/sample_syscc/build/ccs"
	"github.com/hyperledger/fabric/core/system_chaincode/sample_syscc/build/cci/appinit"
	"github.com/hyperledger/fabric/core/system_chaincode/sample_syscc/build/cci/org/hyperledger/chaincode/sample_syscc"
)

// SampleSysCC example simple Chaincode implementation
type SampleSysCC struct {
}

func init() {
	self := &SampleSysCC{}
	interfaces := ccs.Interfaces {
		"org.hyperledger.chaincode.sample_syscc": self,
		"appinit": self,
	}

	err := ccs.Start(interfaces) // Our one instance implements both Transactions and Queries interfaces
	if err != nil {
		fmt.Printf("Error starting example chaincode: %s", err)
	}
}

func (t *SampleSysCC) Init(stub *shim.ChaincodeStub, params *appinit.Init) error {

	// Initialize the chaincode
	return stub.PutState(params.Key, []byte(params.Value))
}

// Transaction makes payment of X units from A to B
func (t *SampleSysCC) PutVal(stub *shim.ChaincodeStub, params *sample_syscc.KeyValue) error {

	_, err := stub.GetState(params.Key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get val for " + params.Key + "\"}"
		return errors.New(jsonResp)
	}

	// Write the state to the ledger
	err = stub.PutState(params.Key, []byte(params.Value))
	if err != nil {
		return err
	}

	return nil
}

// Query callback representing the query of a chaincode
func (t *SampleSysCC) GetVal(stub *shim.ChaincodeStub, params *sample_syscc.Item) (*sample_syscc.Item, error) {

	// Get the state from the ledger
	valbytes, err := stub.GetState(params.Data)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + params.Data + "\"}"
		return nil, errors.New(jsonResp)
	}

	if valbytes == nil {
		jsonResp := "{\"Error\":\"Nil val for " + params.Data + "\"}"
		return nil, errors.New(jsonResp)
	}

	return &sample_syscc.Item{Data: string(valbytes)}, nil
}
