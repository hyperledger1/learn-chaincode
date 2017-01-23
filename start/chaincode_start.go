/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"strings"

         "github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


type Laptop struct{
	ID string `json:"id"`					//the fieldtags are needed to keep case from bouncing around
	User string `json:"user"`
        RAM  string `json:"ram"`
        ROM string `json:"rom"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	 else if function == "init_laptop" {									//create a new marble
		return t.init_laptop(stub, args)
}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

return nil, errors.New("Received unknown function query: " + function)
}


func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}


func (t *SimpleChaincode) init_laptop(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "1x22354", "sarath", "4GB", "1TB"
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}


	fmt.Println("Starting init laptop")
	id := args[0]
	user := args[1]
	ram := args[2]
	rom := args[3]

	//check if laptop already exists
	laptopAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil, errors.New("Failed to get marble name")
	}
	del := Laptop{}
	json.Unmarshal(marbleAsBytes, &del)
	if del.ID == id{
		fmt.Println("This laptop arleady exists in the network: " + id)
		fmt.Println(del);
		return nil, errors.New("This laptop arleady exists")		
	}
	
	//build the laptop  json string manually
	str := `{"id": "` + id + `", "user": "` + user + `", "ram": ` + ram + `, "rom": "` + rom + `"}`
	err = stub.PutState(id, []byte(str))									//store marble with id as key
	if err != nil {
		return nil, err
	}
		
	fmt.Println("- end init marble")
	return nil, nil
}
