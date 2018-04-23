package main

import (
	"bytes"
    "crypto/rsa"
    "crypto/rand"
    "crypto/x509"
	"encoding/json"
	"fmt"
    "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

    "time"
)

var HONGKONG = time.FixedZone("Hong Kong", 8*3600)

// Define the Smart Contract structure
type PrivateSmartContract struct {
}

type Gender int
const (
    Male = iota
    Female
)

// Define the staff structure, with 7 properties.  Structure tags are used by encoding/json library
type Staff struct {
    HKID        string      `json:"id"`
    FirstName   string      `json:"firstname"`
    LastName    string      `json:"lastname"`
    Birth       time.Time   `json:"birthday"`
    Gender      Gender      `json:"gender"`
    Department  string      `json:"dpt"`
    Position    string      `json:"pos"`
    Start       time.Time   `json:"start"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *PrivateSmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *PrivateSmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
    switch function {
    case "initLedger": // i
        return s.initLedger(APIstub)
    case "put": // i
        return s.put(APIstub, args)
    case "getPriKey": // q
        return s.getPriKey(APIstub)
    case "genPriKey": // i
        return s.genPriKey(APIstub)
    case "query": // q
        return s.query(APIstub, args)
    case "get": // q
        return s.get(APIstub, args)
    case "getAll": // q
        return s.getAllStaffs(APIstub)
    }
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *PrivateSmartContract) get(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

	staffAsBytes, _ := APIstub.GetState("Staff."+args[0])
	return shim.Success(staffAsBytes)
}

/**
 *  args
 *  0.  HKID
 *  1.  Name
 *  2.  Birth.year
 *  3.  Birth.month
 *  4.  Birth.day
 *  5.  Gender
 *  6.  Department
 *  7.  Position
 *  8.  Start.year
 *  9.  Start.month
 *  10. Start.day
 **/
func (s *PrivateSmartContract) put(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12.")
	}

    // Birth
    birthYear, _ := strconv.Atoi(args[3])
    birthMonth, _ := strconv.Atoi(args[4])
    birthDay, _ := strconv.Atoi(args[5])
    birth := time.Date (
        birthYear,
        time.Month(birthMonth),
        birthDay,
        0,0,0,0,
        HONGKONG,
    )
    // Gender
    gender, _ := strconv.Atoi(args[6])
    // Start
    startYear, _ := strconv.Atoi(args[9])
    startMonth, _ := strconv.Atoi(args[10])
    startDay, _ := strconv.Atoi(args[11])
    start := time.Date (
        startYear,
        time.Month(startMonth),
        startDay,
        0,0,0,0,
        HONGKONG,
    )
    // Make new Staff struct
    staff := Staff{
        HKID: args[0],
        FirstName: args[1],
        LastName: args[2],
        Birth: birth,
        Gender: Gender(gender),
        Department: args[7],
        Position: args[8],
        Start: start,
    }
    // Make json and put into database
    staffAsBytes, _ := json.Marshal(staff)
    APIstub.PutState("Staff."+args[0], staffAsBytes)
	return shim.Success(staffAsBytes)
}

func (s *PrivateSmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	staffs := []Staff{
		Staff{
            HKID: "0",
            FirstName: "Kaihin",
            LastName: "Wong",
            Birth: time.Date(1996, time.Month(4), 25, 0,0,0,0, HONGKONG),
            Gender: Male,
            Position: "Docker", // not typo
            Start: time.Date(2018, time.Month(1), 1, 0,0,0,0, HONGKONG),
        },
		Staff{
            HKID: "1",
            FirstName: "Dennis",
            LastName: "Chau",
            Birth: time.Date(1996, time.Month(1), 1, 0, 0, 0, 0, HONGKONG),
            Gender: Male,
            Position: "Doctor",
            Start: time.Date(2018, time.Month(1), 1, 0,0,0,0, HONGKONG),
        },
	}

	var i int = 0
	for i < len(staffs) {
		fmt.Println("i is ", i)
		staffAsBytes, _ := json.Marshal(staffs[i])
		APIstub.PutState("Staff." + staffs[i].HKID, staffAsBytes)
		fmt.Println("Added", staffs[i])
		i = i + 1
	}
	return shim.Success(nil)
}

// Advaned: arbitrary query using query language
func (s *PrivateSmartContract) query(APIstub shim.ChaincodeStubInterface, query []string) sc.Response {

    fmt.Println("Querying...") 
    resultsIterator, err := APIstub.GetQueryResult(query[0])
    defer resultsIterator.Close()
    if err != nil {
        return shim.Error(err.Error())
    }
    var buffer bytes.Buffer
    buffer.WriteString("[")
    delimit := false
    for resultsIterator.HasNext() {
        queryResponse,
        err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        
        if delimit == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{")
        buffer.WriteString(queryResponse.Key)
        buffer.WriteString(": ")
        buffer.WriteString(string(queryResponse.Value))
        buffer.WriteString("}")
        delimit = true
    }
    buffer.WriteString("]")
    fmt.Println("Result:\n", buffer.String())
    return shim.Success(nil)
}

func (s *PrivateSmartContract) getAllStaffs(APIstub shim.ChaincodeStubInterface) sc.Response {

    resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("\n[\n")

	delimit := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if delimit {
			buffer.WriteString(",\n")
		}
		buffer.WriteString("    {\n        ")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(":\n            ")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("\n    }")
		delimit = true
	}
	buffer.WriteString("\n]")

	fmt.Println("- getAllStaffs:\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *PrivateSmartContract) getPriKey(APIstub shim.ChaincodeStubInterface) sc.Response {
    priKeyAsBytes, _ := APIstub.GetState("priKey")
    return shim.Success(priKeyAsBytes)
}

func (s *PrivateSmartContract) genPriKey(APIstub shim.ChaincodeStubInterface) sc.Response {

    priKeyAsBytes, err := APIstub.GetState("priKey")
    if err != nil {
        return shim.Error(err.Error())
    }
    if priKeyAsBytes != nil {
        return shim.Error("Key exists")
    }
    reader := rand.Reader
    bitSize := 2048
    priKey, err := rsa.GenerateKey(reader, bitSize)
    if err != nil {
        return shim.Error(err.Error())
    }

    priKeyAsBytes = x509.MarshalPKCS1PrivateKey(priKey)
    APIstub.PutState("priKey", priKeyAsBytes)

    return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(PrivateSmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}
