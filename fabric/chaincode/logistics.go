package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"log"
	"math/big"
	"strconv"
	"strings"
)

type SmartContract struct {
}

type MAMChannel struct {
	ContainerID       string       `json:"ContainerID"`
	Root              string       `json:"Root"`
	SideKey           string       `json:"SideKey"`
}

// logisticstrans type
type logisticstrans struct {
	//product might be food,medical,other itmes
	//Product id should be unique such as FISH123,Prawns456,ICECREAM789
	LogisticstranID   string       `json:"LogisticstranID"`
	ProductID         string       `json:"ProductID"`
	ProductType       string       `json:"ProductType"`
	SellerID          string       `json:"SellerID"`
	SellerLocation    string       `json:"SellerLocation"`
	BuyerID           string       `json:"BuyerID"`
	BuyerLocation     string       `json:"BuyerLocation"`
	LogisticsID       string       `json:"LogisticsID"`
	LogisticsLocation string       `json:"LogisticsLocation"`
	JourneyStartTime  string       `json:",JourneyStartTime"`
	JourneyEndTime    string       `json:",JourneyEndTime"`
	Status            string       `json:"Status"`
	Count            string       `json:"Count"`
	MAMChannel        MAMChannel   `json:"MAMChannel"`
}

type Container struct {
	ContainerID string `json:"ContainerID"`
	Description string `json:"Description"`
	Timestamp string `json:"Timestamp"`
	Location  string `json:"Location"`
	Used      string   `json:"Used"`  //true is used and false is not used
	Holder  string `json:"Holder"`
}

type IotaPayload struct {
	ContainerID string `json:"ContainerID"`
	Seed        string `json:"Seed"`
	MamState    string `json:"MamState"`
	Root        string `json:"Root"`
	Mode       	string `json:"Mode"`
	SideKey     string `json:"SideKey"`
}

type Product struct {
	ProductID string `json:"ProductID"`
	Description string `json:"Description"`
	Timestamp string `json:"Timestamp"`
	Location  string `json:"Location"`
	Holder  string `json:"Holder"`
}

func GenerateRandomSeedString(length int) string {
	seed := ""
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ9"

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(27))
		if err != nil {
			log.Fatal(err)
		}
		seed += string(alphabet[n.Int64()])
	}
	return seed
}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("Error with chaincode")
	} else {
		fmt.Println("Chaincode installed successfully")
	}
}

//Init logisticstrans
func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
	//containers := []Container{
	//	Container{ContainerID: "container1",Description: "", Location: "67.0006, -70.5476", Timestamp: timestamp, Holder: "Freight Forwarder",Used:"false"},
	//	Container{ContainerID: "container2",Description: "", Location: "91.2395, -49.4594", Timestamp: timestamp, Holder: "Freight Forwarder",Used:"false"},
	//	Container{ContainerID: "container3",Description: "", Location: "58.0148, 59.01391", Timestamp: timestamp, Holder: "Freight Forwarder",Used:"false"},
	//}
	//
	//i := 0
	//for i < len(containers) {
	//	containerAsBytes, _ := json.Marshal(containers[i])
	//	stub.PutState(containers[i].ContainerID, containerAsBytes)
	//
	//	mode := "restricted"
	//	seed := GenerateRandomSeedString(81)
	//	sideKey := GenerateRandomSeedString(81)
	//
	//	iotaPayload := IotaPayload{ContainerID:containers[i].ContainerID,Seed: seed, MamState: "", Root: "", Mode: mode, SideKey: sideKey}
	//	iotaPayloadAsBytes, _ := json.Marshal(iotaPayload)
	//	stub.PutState("iotapayload" + containers[i].ContainerID, iotaPayloadAsBytes)
	//
	//	fmt.Println("New Asset", strconv.Itoa(i+1), containers[i], seed, mode, sideKey)
	//	i = i + 1
	//}
	fmt.Println("Initiate the chaincode")
	return shim.Success(nil)
}

func (s *SmartContract) RecordContainer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	container := Container{ ContainerID:args[0],Description: args[1], Location: args[2], Timestamp: args[4], Holder: args[3],Used:"false" }

	containerAsBytes, _ := json.Marshal(container)
	err := stub.PutState(args[0], containerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record container: %s", args[0]))
	}

	mode := "restricted"
	seed := args[5]
	sideKey := args[6]

	iotaPayload := IotaPayload{ContainerID:args[0],Seed: seed, MamState: "", Root: "", Mode: mode, SideKey: sideKey}
	iotaPayloadAsBytes, _ := json.Marshal(iotaPayload)
	stub.PutState("iotapayload" + container.ContainerID, iotaPayloadAsBytes)
	fmt.Println("New Asset", args[0], container, seed, mode, sideKey)

	return shim.Success(nil)
}

func (s *SmartContract) QueryContainer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	containerAsBytes, _ := stub.GetState(args[0])
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	json.Unmarshal(containerAsBytes, &container)

	iotaPayloadAsBytes, _ := stub.GetState("iotapayload" + container.ContainerID)
	if iotaPayloadAsBytes == nil {
		return shim.Error("Could not locate IOTA state object")
	}
	iotaPayload := IotaPayload{}
	json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)

	mamstate := map[string]interface{}{}
	mamstate["seed"] = iotaPayload.Seed
	mamstate["root"] = iotaPayload.Root
	mamstate["sideKey"] = iotaPayload.SideKey
	out := map[string]interface{}{}
	out["container"] = container
	out["mamstate"] = mamstate
	result, _ := json.Marshal(out)
	return shim.Success(result)
}

func (s *SmartContract) QueryAllContainers(stub shim.ChaincodeStubInterface) pb.Response {
	startKey := "container1"
	endKey := "container999"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllContainers:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//Invoke logisticstrans
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fun, args := stub.GetFunctionAndParameters()
	fmt.Println("Arguements for function  ", fun)
	switch fun {
	case "RecordContainer":
		return t.RecordContainer(stub, args)
	case "QueryContainer":
		return t.QueryContainer(stub, args)
	case "QueryAllContainers":
		return t.QueryAllContainers(stub)
	case "RequestLogistic":
		return t.RequestLogistic(stub, args)
	case "TransitLogistics":
		return t.TransitLogistics(stub, args)
	case "InTransitLogistics":
		return t.InTransitLogistics(stub, args)
	case "DeliveryLogistics":
		return t.DeliveryLogistics(stub, args)
	case "SignLogistics":
		return t.SignLogistics(stub, args)
	case "QueryLogistics":
		return t.QueryLogistics(stub, args)
	case "QueryAllLogistics":
		return t.QueryAllLogistics(stub, args)
	}
	fmt.Println("Function not found!")
	return shim.Error("Recieved unknown function invocation!")
}

//Genlogistics for

func (t *SmartContract) RequestLogistic(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 9 {
		fmt.Println("Incorrect number of arguments. Expecting 9")
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}
	var logobj = logisticstrans{
		LogisticstranID: args[0],
		ProductID: args[1],
		ProductType: args[2],
		BuyerID: args[3],
		BuyerLocation: args[4],
		SellerID: args[5],
		SellerLocation: args[6],
		LogisticsID: args[7],
		LogisticsLocation: args[8],
	}
	logobj.Status = "Requested"

	logobjasBytes, _ := json.Marshal(logobj)
	err := stub.PutState(args[0], logobjasBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to request logistics: %s", args[0]))
	}
	return shim.Success(nil)
}

//TransitLogistics at the same time measuring the temp details from logistics
func (t *SmartContract) TransitLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 3")
	}
	logisticsAsBytes, _ := stub.GetState(args[0])
	var logisticobj logisticstrans
	json.Unmarshal(logisticsAsBytes, &logisticobj)
	logisticobj.LogisticstranID = args[0]
	logisticobj.MAMChannel.ContainerID = args[1]
	sideKey := args[3]
	logisticobj.MAMChannel.SideKey = sideKey
	timestamp := args[4]
	logisticobj.JourneyStartTime = timestamp

	containerAsBytes, _ := stub.GetState(args[1])
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	json.Unmarshal(containerAsBytes, &container)
	container.Description = args[0]
	if  container.Used == "true" {
		return shim.Error("Container is used")
	}
	container.Used = "true"
	container.Timestamp = timestamp
	container.Location = args[2]
	iotaPayloadAsBytes, _ := stub.GetState("iotapayload" + container.ContainerID)
	if iotaPayloadAsBytes == nil {
		return shim.Error("Could not locate IOTA state object")
	}
	iotaPayload := IotaPayload{}
	json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)
	iotaPayload.SideKey=sideKey
	if logisticobj.Status != "Requested" {
		fmt.Println("we cannnot transit  the product which was not requested")
		return shim.Error("we cannnot transit  the product which was not requested")
	}

	logisticobj.Status = "Ready-Transit"
	logisticsAsBytes, _ = json.Marshal(logisticobj)
	stub.PutState(args[0], logisticsAsBytes)
	containerAsBytes, _ = json.Marshal(container)
	stub.PutState(container.ContainerID, containerAsBytes)
	iotaPayloadAsBytes, _ = json.Marshal(iotaPayload)
	stub.PutState("iotapayload" + container.ContainerID, iotaPayloadAsBytes)
	err := stub.SetEvent(`{"From":"Fabric","To":"Iota","Func":"CreateChannel"}`, iotaPayloadAsBytes)
	if err != nil {
		fmt.Println("Could not set event for loan application creation", err)
	}
	return shim.Success(nil)
}

func (t *SmartContract) InTransitLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 3")
	}
	containerAsBytes, _ := stub.GetState(args[0])
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	json.Unmarshal(containerAsBytes, &container)
	logisticsAsBytes, _ := stub.GetState(container.Description)

	var logisticobj logisticstrans
	json.Unmarshal(logisticsAsBytes, &logisticobj)
	logisticobj.LogisticstranID = container.Description
	logisticobj.MAMChannel.Root = args[1]
	iotaPayloadAsBytes, _ := stub.GetState("iotapayload" + logisticobj.MAMChannel.ContainerID)
	if iotaPayloadAsBytes == nil {
		return shim.Error("Could not locate IOTA state object")
	}
	iotaPayload := IotaPayload{}
	json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)
	iotaPayload.Root = args[1]
	iotaPayload.MamState = args[2]
	if logisticobj.Status != "Ready-Transit" {
		fmt.Println("we cannnot transit  the product which was not Ready_Transit")
		return shim.Error("we cannnot transit  the product which was not Ready_Transit")
	}

	logisticobj.Status = "In-Transit"
	logisticsAsBytes, _ = json.Marshal(logisticobj)
	stub.PutState(logisticobj.LogisticstranID, logisticsAsBytes)
	iotaPayloadAsBytes, _ = json.Marshal(iotaPayload)
	stub.PutState("iotapayload" + logisticobj.MAMChannel.ContainerID, iotaPayloadAsBytes)
	return shim.Success(nil)
}

func (t *SmartContract) DeliveryLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Invalid   no of arg for delivery function ")

	}
	logisticsasbytes1, _ := stub.GetState(args[0])
	var logisticobj1 logisticstrans

	json.Unmarshal(logisticsasbytes1, &logisticobj1)

	if logisticobj1.Status != "In-Transit" {
		fmt.Println("we cannnot delivery the product which is not in In_Transit")
		return shim.Error("we cannnot delivery the product which is not in In_Transit")
	}
	timestamp := args[1]
	logisticobj1.JourneyEndTime = timestamp
	logisticobj1.Status = "Wait-Sign"
	logisticsasbytes1, _ = json.Marshal(logisticobj1)
	stub.PutState(args[0], logisticsasbytes1)
	iotaPayloadAsBytes, _ := stub.GetState("iotapayload" + logisticobj1.MAMChannel.ContainerID)
	err := stub.SetEvent(`{"From":"Fabric","To":"Iota","Func":"DeliveryLogistics"}`, iotaPayloadAsBytes)
	if err != nil {
		fmt.Println("Could not set event for loan application creation", err)
	}
	return shim.Success(nil)
}

func (t *SmartContract) SignLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Invalid   no of arg for Sign function ")

	}
	containerAsBytes, _ := stub.GetState(args[0])
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	json.Unmarshal(containerAsBytes, &container)
	logisticsasbytes1, _ := stub.GetState(container.Description)
	var logisticobj1 logisticstrans

	json.Unmarshal(logisticsasbytes1, &logisticobj1)

	if logisticobj1.Status != "Wait-Sign" {
		fmt.Println("we cannnot delivery the product which is not in Wait_Sign")
		return shim.Error("we cannnot delivery the product which is not in Wait_Sign")
	}
	fmt.Println("length of the logibj journry in  device", logisticobj1.JourneyEndTime)
	fmt.Println("length of the logibj  journey out timefrrom device", logisticobj1.JourneyStartTime)

	count := 0
	tempStr := strings.Split(args[1],",")
	for i := 0; i < len(tempStr); i++ {
		temp, _ := strconv.Atoi(tempStr[i])
		if temp >= 20 {
			count++
		} else {
			count = 0
		}
	}
	if count >= 3 {
		logisticobj1.Status = "Rejected from Buyer"

	} else {
		logisticobj1.Status = "Accepted  from Buyer"
		container.Used = "false"
	}
	logisticobj1.Count = strconv.Itoa(count)
	containerAsBytes, _ = json.Marshal(container)
	stub.PutState(args[0], containerAsBytes)
	logisticsasbytes1, _ = json.Marshal(logisticobj1)
	stub.PutState(container.Description, logisticsasbytes1)

	return shim.Success(nil)
}


func (t *SmartContract) QueryLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Invalid   no of arg for Query function ")
	}
	logisticsasbytes1, _ := stub.GetState(args[0])
	return shim.Success(logisticsasbytes1)
}

func (t *SmartContract) QueryAllLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	startKey := "logistictran0"
	endKey := "logistictran999"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllLogistics:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

