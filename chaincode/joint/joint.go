package main

import (
	"bytes"
    "crypto/rand"
    "crypto/rsa"
    //"crypto/sha256"
    "crypto/x509"
	"encoding/json"
	"encoding/binary"
	"fmt"
    "math/big"
    "strconv"
    "strings"

    "github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
	sc "github.com/hyperledger/fabric/protos/peer"
    "github.com/pkg/errors"

    "time"
)

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
type Patient struct {
    HKID        string      `json:"id"`
    FirstName   string      `json:"firstname"`
    LastName    string      `json:"lastname"`
    Birth       time.Time   `json:"birthday"`
    Gender      Gender      `json:"gender"`
    Records     []Record    `json:"records"`
    Reports     []Report    `json:"reports"`
}
type PKCS1PublicKey struct {
    N   *big.Int    `json:"N"`
    E   int         `json:"E"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
    switch function {
    case "initLedger": // i
        return s.initLedger(APIstub)
    case "putPatient": // i
        return s.putPatient(APIstub, args)
    case "putRecord": // i
        return s.putRecord(APIstub, args)
    case "putReport": // i
        return s.putReport(APIstub, args)
    case "query": // q
        return s.query(APIstub, args)
    case "get": // q
        return s.get(APIstub, args)
    case "getAll": // q
        return s.queryAllPatients(APIstub)
    case "getPubKey": // q
        return s.getPubKey(APIstub, args)
    case "updatePubKey": // i
        return s.updatePubKey(APIstub, args)
    case "refer": // i
        return s.refer(APIstub, args)
    }

	return shim.Error("Invalid Smart Contract function name.")
}

/** HKID
 */
func (s *SmartContract) get(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }
    peer := args[0]

    ciphertext, _ := APIstub.GetState(peer +"_"+args[1])
    if ciphertext == nil {
        return shim.Error(peer+"_"+args[1])
    }
    
    priKey, err := getPriKey(APIstub, peer)
    if err != nil {
        return shim.Error(err.Error())
    }

    aesKBytes, err := priKey.Decrypt(nil, ciphertext, nil)
    if err != nil {
        return shim.Error(err.Error())
    }
    factory.InitFactories(nil)
    ent, err := entities.NewAES256EncrypterEntity("ID", factory.GetDefault(), aesKBytes, nil)
    if err != nil {
        return shim.Error(err.Error())
    }
    ciphertext, _ = APIstub.GetState(args[1])
    if len(ciphertext) == 0 {
        return shim.Error("Nothing to decrypt")
    }
    patientAsBytes, err := ent.Decrypt(ciphertext)
    if err != nil {
        return shim.Error(err.Error())
    }

    patient := Patient{}
    json.Unmarshal(patientAsBytes, &patient)
    patientIndent, err := json.MarshalIndent(patient, "", "    ")
    if err != nil {
        return shim.Error(err.Error())
    }
	//patientAsBytes, _ := APIstub.GetState("Patient."+args[0])
	//return shim.Success(patientAsBytes)
    return shim.Success(patientIndent)
}

func getChannelName(peer string) string {
    num := "0123456789"
    hospital := strings.TrimRight(peer, num)
    if strings.HasSuffix(hospital, "pubhos") {
        return "publicchannel"
    }
    return "privatechannel"
}

func (s *SmartContract) putPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	hongkong := time.FixedZone("Hong Kong Time", int((8 * time.Hour).Seconds()))
    year, _ := strconv.Atoi(args[3])
    month, _ := strconv.Atoi(args[4])
    day, _ := strconv.Atoi(args[5])
    gender, _ := strconv.Atoi(args[6])
    patient := Patient{
        HKID: args[0],
        FirstName: args[1],
        LastName: args[2],
        Birth: time.Date(
            year,
            time.Month(month),
            day,
            0,0,0,0,
            hongkong,
        ),
        Gender: Gender(gender),
        Records: []Record{},
        Reports: []Report{},
    }
    return s.putEncrypt(APIstub, patient, args[7])
    //patientAsBytes, _ := json.Marshal(patient)
    //APIstub.PutState("Patient."+args[0], patient)
	//return shim.Success(patientAsBytes)
}

func initPutEncrypt(APIstub shim.ChaincodeStubInterface, patient Patient, hid string) sc.Response {
    aesKBytes := make([]byte, 32)
    rand.Reader.Read(aesKBytes)

    queryArgs := toChaincodeArgs("getPriKey")
    var channelName = ""
    num := "0123456789"
    hospital := strings.TrimRight(hid, num)
    if strings.HasSuffix(hospital, "prihos") {
        channelName = "privatechannel"
    } else
    if strings.HasSuffix(hospital, "pubhos") {
        channelName = "publicchannel"
    }
    response := APIstub.InvokeChaincode("disjoint", queryArgs, channelName)
    if response.Status != shim.OK {
        return shim.Error(response.Message)
    }
    priKeyAsBytes := response.Payload
    priKey, err := x509.ParsePKCS1PrivateKey(priKeyAsBytes)
    if err != nil {
        APIstub.PutState("ERRor", []byte(err.Error()))
    }
    pubKey := priKey.PublicKey
    ciphertext, err := rsa.EncryptPKCS1v15(
        //sha256.New(),
        rand.Reader,
        &pubKey,
        aesKBytes,
        //[]byte("end"),
    )
    if err != nil {
        APIstub.PutState("ERR", []byte(err.Error()))
    }
    APIstub.PutState(hid+"_"+patient.HKID, ciphertext)
    return shim.Success(aesKBytes)
}

func (s *SmartContract) putEncrypt(APIstub shim.ChaincodeStubInterface, patient Patient, hid string) sc.Response {
    patientAsBytesEncrypted, err := APIstub.GetState(patient.HKID)
    if err != nil {
        return shim.Error(err.Error())
    }

    var aesKBytes []byte
    if patientAsBytesEncrypted == nil {
        response := initPutEncrypt(APIstub, patient, hid)
        if response.Status != shim.OK {
            return shim.Error(response.Message)
        }
        aesKBytes = response.Payload
    } else {
        aesKBytesEncrypted, err := APIstub.GetState(hid+"_"+patient.HKID)
        if aesKBytesEncrypted == nil {
            return shim.Error("Unauthorized")
        }
        priKey, err := getPriKey(APIstub, hid)
        if err != nil {
            return shim.Error(err.Error())
        }
        aesKBytes, err = priKey.Decrypt(nil, aesKBytesEncrypted, nil)
    }

    factory.InitFactories(nil)
    ent, err := entities.NewAES256EncrypterEntity("ID", factory.GetDefault(), aesKBytes, nil)
    if err != nil {
        return shim.Error(err.Error())
    }
    patientAsBytes, _ := json.Marshal(patient)
    ciphertext, err := ent.Encrypt(patientAsBytes)
    if err != nil {
        return shim.Error(err.Error())
    }
    APIstub.PutState(patient.HKID, ciphertext)
    return shim.Success(nil)
}

/** HKID
 *  Datefrom
 *  DateTo
 *  Department Name
 *  RoomNumber
 *  BedNumber
 */
func (s *SmartContract) putRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

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
		/*
		Department: args[7],
		RoomNumber: args[8],
		BedNumber: args[9],
		*/
    }

    patientAsBytes, _ := APIstub.GetState("Patient."+args[0])
    patient := Patient{}
    json.Unmarshal(patientAsBytes, &patient)
    patient.Records = append(patient.Records, record)
    patientAsBytes, _ = json.Marshal(patient)
    APIstub.PutState("Patient."+args[0], patientAsBytes)

    nRecords, _ := APIstub.GetState("numRecords") //[]byte,error
    buf := bytes.NewBuffer(nRecords)
    var numRecords uint64
    binary.Read(buf, binary.LittleEndian, &numRecords)
    numRecords = numRecords + 1

    recordAsBytes, _ := json.Marshal(record)
    APIstub.PutState("Records"+strconv.FormatUint(numRecords,10), recordAsBytes)
    binary.Write(buf, binary.LittleEndian, numRecords)
    APIstub.PutState("numRecords", buf.Bytes())
	return shim.Success(recordAsBytes)
}

/** HKID
 *  Datefrom
 *  DateTo
 *  Report Type
 *  verified
 */
func (s *SmartContract) putReport(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

    year, _ := strconv.Atoi(args[1])
    month, _ := strconv.Atoi(args[2])
    day, _ := strconv.Atoi(args[3])
	yesno, _ := strconv.ParseBool(args[5])


    report := Report{
        Date: time.Date(
            year,
            time.Month(month), //change int to Month type
            day,
            0,0,0,0,
            time.UTC,
        ),
        Type: args[4],
        Yesno:  yesno, //change string "true" to bool "true"
    }

    patientAsBytes, _ := APIstub.GetState("Patient."+args[0])
    patient := Patient{}
    json.Unmarshal(patientAsBytes, &patient)
    patient.Reports = append(patient.Reports, report)
    patientAsBytes, _ = json.Marshal(patient)
    APIstub.PutState("Patient."+args[0], patientAsBytes)

    nReports, _ := APIstub.GetState("numReports")
    buf := bytes.NewBuffer(nReports)
    var numReports uint64
    binary.Read(buf, binary.LittleEndian, &numReports)
    numReports = numReports + 1

    reportAsBytes, _ := json.Marshal(report)
    APIstub.PutState("Reports"+strconv.FormatUint(numReports,10), reportAsBytes)
    binary.Write(buf, binary.LittleEndian, numReports)
    APIstub.PutState("numReports", buf.Bytes())
	return shim.Success(reportAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	//setup timezone
    s.putPatient(APIstub, []string{
        "Y000000(1)",
        "Eric",
        "Ho",
        "1996",
        "7",
        "3",
        "Male",
        "cli.icu.pubhos1",
    })
    s.putPatient(APIstub, []string{
        "Y000001(2)",
        "Yung",
        "Yeung",
        "2018",
        "1",
        "1",
        "Female",
        "cli.icu.pubhos2",
    })
    /*
	var i int = 0
	for i < len(patients) {
		fmt.Println("i is ", i)
		patientAsBytes, _ := json.Marshal(patients[i])
		APIstub.PutState("Patient." + patients[i].HKID, patientAsBytes)
		fmt.Println("Added", patients[i])
		i = i + 1
	}
    */
	return shim.Success(nil)
}

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, query []string) sc.Response {

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

func (s *SmartContract) queryAllPatients(APIstub shim.ChaincodeStubInterface) sc.Response {

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
		if delimit == true {
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

	fmt.Println("- queryAllPatients:\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) requestPatient(APIstub shim.ChaincodeStubInterface) sc.Response {

    patientInBytes, _ := APIstub.GetState("Y000000(1)")
    return shim.Success(patientInBytes)
}

func (s *SmartContract) updatePubKey(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }
    id := args[0]
    channelName := args[1]
    queryArgs := toChaincodeArgs("getPriKey")
    response := APIstub.InvokeChaincode("disjoint", queryArgs, channelName)
    if response.Status != shim.OK {
        return shim.Error(response.Message)
    }
    priKeyAsBytes := response.Payload
    priKey, err := x509.ParsePKCS1PrivateKey(priKeyAsBytes)
    if err != nil {
        return shim.Error(err.Error())
    }
    pubKey := priKey.PublicKey
    pubKeyAsBytes, _ := json.Marshal(PKCS1PublicKey{
        N: pubKey.N,
        E: pubKey.E,
    })
    APIstub.PutState(id, pubKeyAsBytes)
    return shim.Success(nil)

}

func (s *SmartContract) getPubKey(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    id := args[0]
    pubKeyAsBytes, err := APIstub.GetState(id)
    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success(pubKeyAsBytes)
}

func getPriKey(APIstub shim.ChaincodeStubInterface, peer string) (*rsa.PrivateKey, error) {
    queryArgs := toChaincodeArgs("getPriKey")
    response := APIstub.InvokeChaincode("disjoint", queryArgs, getChannelName(peer))
    if response.Status != shim.OK {
        return nil, errors.New(response.Message)
    }
    priKeyAsBytes := response.Payload
    return x509.ParsePKCS1PrivateKey(priKeyAsBytes)
}

func (s *SmartContract) refer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    
    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }
    referer := args[0]
    referee := args[1]
    pid := args[2]
    aesKBytesEncrypted, _ := APIstub.GetState(referer+"_"+pid)
    if aesKBytesEncrypted == nil {
        return shim.Error("Referer is unauthorized")
    }
    priKey, err := getPriKey(APIstub, referer)
    if err != nil {
        return shim.Error(err.Error())
    }
    aesKBytes, err := priKey.Decrypt(nil, aesKBytesEncrypted, nil)
    if err != nil {
        return shim.Error(err.Error())
    }
    response := s.getPubKey(APIstub, []string{referee})
    if response.Status != shim.OK {
        return shim.Error(response.Message)
    }
    refereePublicKeyAsBytes := response.Payload
    refereePublicKeyPKCS1 := PKCS1PublicKey{}
    json.Unmarshal(refereePublicKeyAsBytes, &refereePublicKeyPKCS1)
    refereePublicKey := rsa.PublicKey {
        N: refereePublicKeyPKCS1.N,
        E: refereePublicKeyPKCS1.E,
    }
    aesKBytesEncryptedNew, err := rsa.EncryptPKCS1v15(
        rand.Reader,
        &refereePublicKey,
        aesKBytes,
    )
    if err != nil {
        return shim.Error(err.Error())
    }
    APIstub.PutState(referee+"_"+pid, aesKBytesEncryptedNew)
    return shim.Success(nil)
}
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
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
