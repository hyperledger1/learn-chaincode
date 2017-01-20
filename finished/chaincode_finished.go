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
         "encoding/json"


	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

	

type laptop struct{
	Name string `json:"name"`	 					//attributes of laptop
	RAM string  `json:"ram"`	 

	ROM string   `json:"rom"`	
	User string  `json:"user"`	 
}







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

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}else if function == "init_laptop"{

                return t.init_laptop(stub,args)
       }         
         
       



	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}


func (t *SimpleChaincode) init_laptop(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
var err error    //to capture any errors

//1x22 4GB 1TB thrinath
//name RAM ROM User

if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
}

fmt.Println("- start init laptop")

Name := args[0]
RAM := args[1]
ROM := args[2]
User := args[3]

//check if laptop already exists
laptopAsBytes, err := stub.GetState(Name)
	if err != nil {
		return nil, errors.New("Failed to get marble name")
	}
	res :=laptop{}
	json.Unmarshal(laptopAsBytes, &res)
	if res.Name == Name{


		fmt.Println("This laptop` arleady exists: " +Name)
		fmt.Println(res);
		return nil, errors.New("This laptop arleady exists")				//all stop a marble by this name exists
}

str := `{"name": "` + Name + `", "ram": "` + RAM + `", "rom": ` + ROM + `, "user": "` + User + `"}`
	err = stub.PutState(Name, []byte(str))							//store marble with id as key
	if err != nil {
		return nil, err
}
fmt.Println("- end init marble")
return nil, nil
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

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
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
