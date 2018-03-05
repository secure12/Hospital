/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"encoding/binary"
	"fmt"
    "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

    "time"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Gender int
const (
    Male = iota
    Female
)

type Record struct {
    DateFrom    time.Time   `json:"from"`
    DateTo      time.Time   `json:"to"`
    Department  string      `json:"dpt"`
    RoomNumber  string      `json:"room"`
    BedNumber   string      `json:"bed"`
}

type Report struct {
    Date    time.Time   `json:"date"`
    Type    string      `json:"type"`
    Yesno   bool        `json:"yesno"`
}

// Define the patient structure, with 4 properties.  Structure tags are used by encoding/json library
type Patient struct {
    HKID    string      `json:"id"`
    Name    string      `json:"name"`
    Birth   time.Time   `json:"birthday"`
    Gender  Gender      `json:"gender"`
    Records []Record    `json:"records"`
    Reports []Report    `json:"reports"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
    if function == "initLedger" {
        return s.initLedger(APIstub)
    } else
	if function == "query" {
        return s.query(APIstub, args)
    } else
    if function == "get" {
        return s.get(APIstub, args)
    } else
    if function == "putPatient" {
        return s.putPatient(APIstub, args)
    } else
    if function == "putRecord" {
        return s.putRecord(APIstub, args)
    } else
    if function == "putReport" {
        return s.putReport(APIstub, args)
    }

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) get(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(patientAsBytes)
}

func (s *SmartContract) putPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

    year, _ := strconv.Atoi(args[2])
    month, _ := strconv.Atoi(args[3])
    day, _ := strconv.Atoi(args[4])
    gender, _ := strconv.Atoi(args[5])
    patient := Patient{
        HKID: args[0],
        Name: args[1],
        Birth: time.Date(
            year,
            time.Month(month),
            day,
            0,0,0,0,
            time.UTC,
        ),
        Gender: Gender(gender),
        Records: []Record{},
        Reports: []Report{},
    }

    patientAsBytes, _ := json.Marshal(patient)
    APIstub.PutState(args[0], patientAsBytes)
	return shim.Success(patientAsBytes)
}

/** HKID
 *  Datefrom
 *  DateTo
 *  Department Name
 *  RoomNumber
 *  BedNumber
 */
func (s *SmartContract) putRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    yearFrom, _ := strconv.Atoi(args[1])
    monthFrom, _ := strconv.Atoi(args[2])
    dayFrom, _ := strconv.Atoi(args[3])

    yearTo, _ := strconv.Atoi(args[4])
    monthTo, _ := strconv.Atoi(args[5])
    dayTo, _ := strconv.Atoi(args[6])

    record := Record{
        DateFrom: time.Date(
            yearFrom,
            time.Month(monthFrom),
            dayFrom,
            0,0,0,0,
            time.UTC,
        ),
        DateTo: time.Date(
            yearTo,
            time.Month(monthTo),
            dayTo,
            0,0,0,0,
            time.UTC,
        ),
        Department: "ABC",
        RoomNumber: args[7],
        BedNumber: args[8],
    }

    patientAsBytes, _ := APIstub.GetState(args[0])
    patient := Patient{}
    json.Unmarshal(patientAsBytes, &patient)
    patient.Records = append(patient.Records, record)
    patientAsBytes, _ = json.Marshal(patient)
    APIstub.PutState(args[0], patientAsBytes)

    nRecords, _ := APIstub.GetState("numRecords")
    numRecords := binary.BigEndian.Uint64(nRecords)
    recordsAsBytes, _ := APIstub.GetState("Records"+numRecords)
    recordAsBytes, _ = json.Marshal(record)
    recordsAsBytes = append(recordsAsBytes, recordAsBytes)
    APIstub.PutState("Records"+numRecords, recordsAsBytes)
	return shim.Success(recordAsBytes)
}

func (s *SmartContract) putReport(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    year, _ := strconv.Atoi(args[1])
    month, _ := strconv.Atoi(args[2])
    day, _ := strconv.Atoi(args[3])

    report := Report{
        Date: time.Date(
            year,
            month,
            day,
            0,0,0,0,
            time.UTC,
        ),
        Type: args[4],
        Yesno: args[5],
    }

    patientAsBytes, _ := APIstub.GetState(args[0])
    patient := Patient{}
    json.Unmarshal(patientAsBytes, &patient)
    patient.Reports = append(patient.Reports, report)
    patientAsBytes, _ := json.Marshal(patient)
    APIstub.PutState(args[0], patientAsBytes)

    nReports := APIsutb.GetState("numReports")
    numReports := binary.BigEndian.Uint64(nReports)
    reportsAsBytes, _ := APIStub.GetState("Reports"+numReports)
    reportAsBytes, _ := json.Marshal(report)
    reportsAsBytes = append(reportsAsBytes, reportAsBytes)
    APIstub.PutState("Reports"+numReports, reportsAsBytes)
	return shim.Success(reportAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	patients := []Patient{
		Patient{
            HKID: "Y000000(1)",
            Name: "Aa",
            Birth: time.Date(1996, 7, 3, 0, 0, 0, 0, time.UTC+8),
            Gender: Male,
            Records: []Record{},
            Reports: []Report{},
        },
		Patient{
            HKID: "Y000001(2)",
            Name: "Bb",
            Birth: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC+8),
            Gender: Male,
            Records: []Record{},
            Reports: []Report{},
        },
		Patient{
            HKID: "Y000002(3)",
            Name: "Cc",
            Birth: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC+8),
            Gender: Female,
            Records: []Record{},
            Reports: []Report{},
        },
		Patient{
            HKID: "Y000003(4)",
            Name: "Dd",
            Birth: time.Date(20, 1, 1, 0, 0, 0, 0, time.UTC+8),
            Gender: Male,
            Records: []Record{},
            Reports: []Report{},
        },
	}

	i := 0
	for i < len(patients) {
		fmt.Println("i is ", i)
		patientAsBytes, _ := json.Marshal(patients[i])
		APIstub.PutState("Patient."+patients[i][HKID], PatientAsBytes)
		fmt.Println("Added", patients[i])
		i = i + 1
	}

    APIstub.PutState("numRecords", '0')
    APIstub.PutState("numReports", '0')
	return shim.Success(nil)
}

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, query []string) sc.Response {

    fmt.Println("Querying...") 
    resultsIterator, err := APIstub.getQueryResult(query)
    defer resultsIterator.Close()
    if err != nil {
        return nil, err
    }

    var buffer bytes.Buffer
    buffer.WriteString("[")
    delimit := false
    for resultsIterator.HasNext() {
        queryResponse,
        err := resultsIterator.Next()
        if err != nil {
            return nil, err
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

func (s *SmartContract) queryAllPatients(APIstub shim.ChaincodeStubInterface) sc.Response {

    resultsIterator, err := APIstub.GetStateByRange()
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	delimit := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if delimit == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(": ")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		delimit = true
	}
	buffer.WriteString("]")

	fmt.Println("- queryAllPatients:\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
