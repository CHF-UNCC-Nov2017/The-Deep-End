package main

/* This code is based off the carTrace demo file */

import (
	/*
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"math/big"
	"strings"
	"time"
	*/

	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


// AutoTraceChaincode example simple Chaincode implementation
type LoanTraceChaincode struct {
}

// @MODIFY_HERE add recall fields to vehicle JSON object
type loan struct {
	SerialNumber string `json:"SerialNumber"`
	SenderAccount string `json:"SenderAccount"`
	DestinationAccount string `json:"DestinationAccount"`
	Value int `json:"Value"` // TODO float
	InitiationDate int `json:"InitiationDate"`
	CompletionDate int `json:"CompletionDate"` // to be added at workshop
	Satisfied string `json:"Satisfied"`
}

func main() {
	err := shim.Start(new(LoanTraceChaincode))
	if err != nil {
		fmt.Printf("Error starting Loan Trace chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *LoanTraceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *LoanTraceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initLoan" { //create a new vehicle
		return t.initLoan(stub, args)
	} else if function == "updateLoan" { //change owner of a specific vehicle
		return t.updateLoan(stub, args)
	} else if function == "readLoan" { //read a vehicle
		return t.readLoan(stub, args)
	}

	// @MODIFY_HERE
	// ==== Write a sub-routine to mark a vehicle part as recalled by ".Name"

	// @MODIFY_HERE
	// ==== Write a sub-routine to mark a vehicle as recalled by ".Manufacturer" & ".Model"

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}


// ============================================================
// initVehicle - create a new vehicle , store into chaincode state
// ============================================================
func (t *LoanTraceChaincode) initLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error


	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}


	// ==== Input sanitization ====
	fmt.Println("- start init loan")
	serialNumber := args[0]
	senderAccount := args[1]
	destinationAccount := args[2]
	value, err := strconv.Atoi(args[3])
	initiationDate, err := strconv.Atoi(args[4])
	completionDate, err := strconv.Atoi(args[5])
	satisfied := args[6]


	// @MODIFY_HERE parts recall fields
	// ==== Create vehicle object and marshal to JSON ====
	// objectType := "loan"
	//vehicle := &vehicle{objectType, chassisNumber, manufacturer, model, assemblyDate, airbagSerialNumber, owner}
	//vehicle := &vehicle{objectType, chassisNumber, manufacturer, model, assemblyDate, airbagSerialNumber, owner, recall, recallDate}
	myLoan := &loan{serialNumber, senderAccount, destinationAccount, value, initiationDate, completionDate, satisfied}

	loanJSONasBytes, err := json.Marshal(myLoan)
	if err != nil {
		return shim.Error(err.Error())
	}


	// === Save vehicle to state ===
	err = stub.PutState(serialNumber, loanJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	/*
	//  ==== Index the vehicle parts to enable assember & owner-based range queries, e.g. return all tata parts ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~assember~chassisNumber.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "manufacturer~chassisNumber"
	ownersIndex := "owner~identifier"
	err = t.createIndex(stub, indexName, []string{vehicle.Manufacturer, vehicle.ChassisNumber})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = t.createIndex(stub, ownersIndex, []string{vehicle.Owner, vehicle.ChassisNumber})
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Vehicle part saved and indexed. Return success ====
	fmt.Println("- end init vehicle")
	*/
	return shim.Success(nil)
}

// ===============================================
// readLoan - read a loan from chaincode state
// ===============================================
func (t *LoanTraceChaincode) readLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var serialNumber, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting serial number of the loan to query")
	}

	serialNumber = args[0]
	valAsbytes, err := stub.GetState(serialNumber) //get the vehicle from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + serialNumber + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Loan does not exist: " + serialNumber + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *LoanTraceChaincode) updateLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}


