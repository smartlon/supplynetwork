package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
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
	ProductType string `json:"ProductType"`
	Description string `json:"Description"`
	Timestamp string `json:"Timestamp"`
	Status string `json:"Status"`
	Holder  string `json:"Holder"`
}

type Participant struct {
	UserName string `json:"UserName"`
	Affiliation string `json:"Affiliation"`
	Location string `json:"Location"`
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

	fmt.Println("Initiate the chaincode")
	return shim.Success(nil)
}

func (s *SmartContract) RecordParticipant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	_,_,affiliation , err := ABAC(stub)
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	participant := Participant{UserName:args[0],Affiliation:affiliation,Location:args[1]}
	participantAsBytes,err := json.Marshal(participant)
	if err != nil {
		return shim.Error(err.Error())
	}
	participantKey, err := stub.CreateCompositeKey("Participant", []string{affiliation,args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(participantKey, participantAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record participant: %s", participantKey))
	}
	fmt.Printf("- RecordParticipant:\nkey = %s, value = %s\n", participantKey,string(participantAsBytes))
	return shim.Success(nil)
}

func (s *SmartContract) QueryAllParticipant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	resultsIterator, err := stub.GetStateByPartialCompositeKey("Participant", []string{})
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

	fmt.Printf("- queryAllParticipant:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func ABAC(stub shim.ChaincodeStubInterface) (string,string, string, error) {
	// Get the client ID object
	cIdn, err := cid.New(stub)
	if err != nil {
		return "","","",err
	}
	id, err := cIdn.GetID()
	if err != nil {
		return "","","",err
	}

	mspid, err := cIdn.GetMSPID()
	if err != nil {
		return "","","",err
	}
	affiliation, ok, err := cIdn.GetAttributeValue("hf.Affiliation")
	if err != nil {
		return "","","",err
	}
	if !ok {
		return "","","",err
	}
	return id,mspid,affiliation,nil
}

func (s *SmartContract) RecordProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id,_,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	product := Product{ProductID:args[0],ProductType:args[1],Description:args[2],Timestamp:args[3],Holder:id,Status:"false"}
	productAsBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}
	productKey, err := stub.CreateCompositeKey("Product", []string{affiliation,args[4],args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(productKey, productAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record product: %s", productKey))
	}
	fmt.Printf("- RecordProduct:\nkey = %s, value = %s\n", productKey,string(productAsBytes))
	return shim.Success(nil)
}

func (s *SmartContract) QueryAllProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	_,_,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	resultsIterator, err := stub.GetStateByPartialCompositeKey("Product", []string{affiliation,args[0]})
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

	fmt.Printf("- queryAllProduct:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) RecordContainer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id,mspid,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if mspid != "DelivererMSP"{
		return shim.Error(fmt.Sprintf("the user %s is not belong to the deliverer",id))
	}
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	container := Container{ ContainerID:args[0],Description: args[1], Location: args[2], Timestamp: args[3], Holder: id,Used:"false" }
	containerAsBytes, err := json.Marshal(container)
	if err != nil {
		return shim.Error(err.Error())
	}
	containerKey, err := stub.CreateCompositeKey("Container", []string{affiliation,args[6],args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(containerKey, containerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record container: %s", containerKey))
	}

	mode := "restricted"
	seed := args[4]
	sideKey := args[5]
	iotaPayload := IotaPayload{ContainerID:containerKey,Seed: seed, MamState: "", Root: "", Mode: mode, SideKey: sideKey}
	iotaPayloadAsBytes, err := json.Marshal(iotaPayload)
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaKey, err := stub.CreateCompositeKey("IotaPayload", []string{affiliation,args[6],args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(iotaKey, iotaPayloadAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("- RecordProduct:\ncontainer key = %s, value = %s\n", containerKey,string(containerAsBytes))
	fmt.Printf("- RecordProduct:\nIotaPayload key = %s, value = %s\n", iotaKey,string(iotaPayloadAsBytes))
	return shim.Success(nil)
}

func (s *SmartContract) QueryContainer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id,mspid,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if mspid != "DelivererMSP"{
		return shim.Error(fmt.Sprintf("the user %s is not belong to the deliverer",id))
	}
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	containerKey, err := stub.CreateCompositeKey("Container", []string{affiliation,args[1],args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	containerAsBytes, err := stub.GetState(containerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	err = json.Unmarshal(containerAsBytes, &container)
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaKey, err := stub.CreateCompositeKey("IotaPayload", []string{affiliation,args[1],args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaPayloadAsBytes, err := stub.GetState(iotaKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if iotaPayloadAsBytes == nil {
		return shim.Error("Could not locate IOTA state object")
	}
	iotaPayload := IotaPayload{}
	err = json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)
	if err != nil {
		return shim.Error(err.Error())
	}
	mamstate := map[string]interface{}{}
	mamstate["seed"] = iotaPayload.Seed
	mamstate["root"] = iotaPayload.Root
	mamstate["sideKey"] = iotaPayload.SideKey
	out := map[string]interface{}{}
	out["container"] = container
	out["mamstate"] = mamstate
	result, err := json.Marshal(out)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func (s *SmartContract) QueryAllContainers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id,mspid,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if mspid != "DelivererMSP"{
		return shim.Error(fmt.Sprintf("the user %s is not belong to the deliverer",id))
	}
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	resultsIterator, err := stub.GetStateByPartialCompositeKey("Container", []string{affiliation,args[0]})
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
	case "RecordParticipant":
		return t.RecordParticipant(stub, args)
	case "QueryAllParticipant":
		return t.QueryAllParticipant(stub, args)
	case "RecordProduct":
		return t.RecordProduct(stub, args)
	case "QueryAllProduct":
		return t.QueryAllProduct(stub, args)
	case "RecordContainer":
		return t.RecordContainer(stub, args)
	case "QueryContainer":
		return t.QueryContainer(stub, args)
	case "QueryAllContainers":
		return t.QueryAllContainers(stub, args)
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
	_,_,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	productKey, err := stub.CreateCompositeKey("Product", []string{affiliation,args[4],args[1]})
	if err != nil {
		return shim.Error(err.Error())
	}
	productAsBytes, err := stub.GetState(productKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if productAsBytes == nil {
		return shim.Error("Could not locate product")
	}
	product := Product{}
	err = json.Unmarshal(productAsBytes, &product)
	if err != nil {
		return shim.Error(err.Error())
	}
	product.Status = "true"
	productAsBytes, err = json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(productKey, productAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to deliver product: %s", productKey))
	}
	sellerParticipant := Participant{}
	sellerKey, err := stub.CreateCompositeKey("Participant", []string{affiliation,args[4]})
	if err != nil {
		return shim.Error(err.Error())
	}
	participantAsBytes,err  := stub.GetState(sellerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(participantAsBytes,sellerParticipant)
	if err != nil {
		return shim.Error(err.Error())
	}
	buyerParticipant := Participant{}
	buyerKey, err := stub.CreateCompositeKey("Participant", []string{args[2]})
	if err != nil {
		return shim.Error(err.Error())
	}
	participantAsBytes,err  = stub.GetState(buyerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(participantAsBytes,buyerParticipant)
	if err != nil {
		return shim.Error(err.Error())
	}
	logistiParticipant := Participant{}
	logistiKey, err := stub.CreateCompositeKey("Participant", []string{args[3]})
	if err != nil {
		return shim.Error(err.Error())
	}
	participantAsBytes,err  = stub.GetState(logistiKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(participantAsBytes,logistiParticipant)
	if err != nil {
		return shim.Error(err.Error())
	}
	var logobj = logisticstrans{
		LogisticstranID: args[0],
		ProductID: args[1],
		ProductType: product.ProductType,
		BuyerID: args[2],
		BuyerLocation: buyerParticipant.Location,
		SellerID: affiliation+args[4],
		SellerLocation: sellerParticipant.Location,
		LogisticsID: args[3],
		LogisticsLocation: logistiParticipant.Location,
	}
	logobj.Status = "Requested"
	logobjasBytes, err := json.Marshal(logobj)
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticstranKey, err := stub.CreateCompositeKey("logisticstrans", []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(logisticstranKey, logobjasBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to request logistics: %s", logisticstranKey))
	}
	fmt.Printf("- RequestLogistic:\nlogisticstrans key = %s, value = %s\n", logisticstranKey,string(logobjasBytes))
	return shim.Success(nil)
}

//TransitLogistics at the same time measuring the temp details from logistics
func (t *SmartContract) TransitLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id,mspid,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if mspid != "DelivererMSP"{
		return shim.Error(fmt.Sprintf("the user %s is not belong to the deliverer",id))
	}
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 5")
	}
	logisticstranKey, err := stub.CreateCompositeKey("logisticstrans", []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticsAsBytes, err := stub.GetState(logisticstranKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	var logisticobj logisticstrans
	json.Unmarshal(logisticsAsBytes, &logisticobj)

	sideKey := args[2]
	logisticobj.MAMChannel.SideKey = sideKey
	timestamp := args[3]
	logisticobj.JourneyStartTime = timestamp
	containerKey, err := stub.CreateCompositeKey("Container", []string{affiliation,args[4],args[1]})
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticobj.MAMChannel.ContainerID = containerKey
	containerAsBytes, err := stub.GetState(containerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	err = json.Unmarshal(containerAsBytes, &container)
	if err != nil {
		return shim.Error(err.Error())
	}
	container.Description = args[0]
	if  container.Used == "true" {
		return shim.Error("Container is used")
	}
	container.Used = "true"
	container.Timestamp = timestamp
	container.Location = logisticobj.SellerLocation
	iotaKey, err := stub.CreateCompositeKey("IotaPayload", []string{affiliation,args[4],args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaPayloadAsBytes, err := stub.GetState(iotaKey )
	if err != nil {
		return shim.Error(err.Error())
	}
	if iotaPayloadAsBytes == nil {
		return shim.Error("Could not locate IOTA state object")
	}
	iotaPayload := IotaPayload{}
	err = json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaPayload.SideKey=sideKey
	if logisticobj.Status != "Requested" {
		fmt.Println("we cannnot transit  the product which was not requested")
		return shim.Error("we cannnot transit  the product which was not requested")
	}

	logisticobj.Status = "Ready-Transit"
	logisticsAsBytes, err = json.Marshal(logisticobj)
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(logisticstranKey, logisticsAsBytes)
	containerAsBytes, err = json.Marshal(container)
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(containerKey, containerAsBytes)
	iotaPayloadAsBytes, err = json.Marshal(iotaPayload)
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(iotaKey, iotaPayloadAsBytes)
	err = stub.SetEvent(`{"From":"Fabric","To":"Iota","Func":"CreateChannel"}`, iotaPayloadAsBytes)
	if err != nil {
		fmt.Println("Could not set event for loan application creation", err)
	}
	return shim.Success(nil)
}

func (t *SmartContract) InTransitLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 3")
	}
	containerAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	err = json.Unmarshal(containerAsBytes, &container)
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticstranKey, err := stub.CreateCompositeKey("logisticstrans", []string{container.Description})
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticsAsBytes, err := stub.GetState(logisticstranKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	var logisticobj logisticstrans
	json.Unmarshal(logisticsAsBytes, &logisticobj)
	logisticobj.LogisticstranID = container.Description
	logisticobj.MAMChannel.Root = args[1]
	iotaKey, err := stub.CreateCompositeKey("IotaPayload", []string{strings.TrimLeft(args[0],"Container")})
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaPayloadAsBytes, err := stub.GetState(iotaKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if iotaPayloadAsBytes == nil {
		return shim.Error("Could not locate IOTA state object")
	}
	iotaPayload := IotaPayload{}
	err=json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaPayload.Root = args[1]
	iotaPayload.MamState = args[2]
	if logisticobj.Status != "Ready-Transit" {
		fmt.Println("we cannnot transit  the product which was not Ready_Transit")
		return shim.Error("we cannnot transit  the product which was not Ready_Transit")
	}

	logisticobj.Status = "In-Transit"
	logisticsAsBytes, err = json.Marshal(logisticobj)
	if err != nil {
		return shim.Error(err.Error())
	}
	err= stub.PutState(logisticstranKey, logisticsAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaPayloadAsBytes, err = json.Marshal(iotaPayload)
	if err != nil {
		return shim.Error(err.Error())
	}
	err=stub.PutState(iotaKey, iotaPayloadAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SmartContract) DeliveryLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id,mspid,_ , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if mspid != "DelivererMSP"{
		return shim.Error(fmt.Sprintf("the user %s is not belong to the deliverer",id))
	}
	if len(args) != 2 {
		return shim.Error("Invalid   no of arg for delivery function ")

	}
	logisticstranKey, err := stub.CreateCompositeKey("logisticstrans", []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticsasbytes1, err := stub.GetState(logisticstranKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	var logisticobj1 logisticstrans

	json.Unmarshal(logisticsasbytes1, &logisticobj1)

	if logisticobj1.Status != "In-Transit" {
		fmt.Println("we cannnot delivery the product which is not in In_Transit")
		return shim.Error("we cannnot delivery the product which is not in In_Transit")
	}
	timestamp := args[1]
	logisticobj1.JourneyEndTime = timestamp
	logisticobj1.Status = "Wait-Sign"
	logisticsasbytes1, err = json.Marshal(logisticobj1)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(logisticstranKey, logisticsasbytes1)
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaKey, err := stub.CreateCompositeKey("IotaPayload", []string{strings.TrimLeft(logisticobj1.MAMChannel.ContainerID,"Container")})
	if err != nil {
		return shim.Error(err.Error())
	}
	iotaPayloadAsBytes, err := stub.GetState(iotaKey )
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent(`{"From":"Fabric","To":"Iota","Func":"DeliveryLogistics"}`, iotaPayloadAsBytes)
	if err != nil {
		fmt.Println("Could not set event for loan application creation", err)
	}
	return shim.Success(nil)
}

func (t *SmartContract) SignLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Invalid   no of arg for Sign function ")

	}
	containerAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	err=json.Unmarshal(containerAsBytes, &container)
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticstranKey, err := stub.CreateCompositeKey("logisticstrans", []string{container.Description})
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticsasbytes1, err := stub.GetState(logisticstranKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	var logisticobj1 logisticstrans

	err=json.Unmarshal(logisticsasbytes1, &logisticobj1)
	if err != nil {
		return shim.Error(err.Error())
	}
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
		product := Product{}
		productKey, err := stub.CreateCompositeKey("Product", []string{logisticobj1.SellerID,logisticobj1.ProductID})
		if err != nil {
			return shim.Error(err.Error())
		}
		productAsBytes,err := stub.GetState( productKey )
		if err != nil {
			return shim.Error(err.Error())
		}
		err=json.Unmarshal(productAsBytes,product)
		if err != nil {
			return shim.Error(err.Error())
		}
		product.Status = "false"
		productAsBytes,err = json.Marshal(product)
		if err != nil {
			return shim.Error(err.Error())
		}
		err=stub.PutState(productKey,productAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	logisticobj1.Count = strconv.Itoa(count)
	containerAsBytes, err = json.Marshal(container)
	if err != nil {
		return shim.Error(err.Error())
	}
	err=stub.PutState(args[0], containerAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticsasbytes1, err = json.Marshal(logisticobj1)
	if err != nil {
		return shim.Error(err.Error())
	}
	err=stub.PutState(logisticstranKey, logisticsasbytes1)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


func (t *SmartContract) QueryLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Invalid   no of arg for Query function ")
	}
	logisticKey, err := stub.CreateCompositeKey("logisticstrans", []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	logisticsasbytes1, err := stub.GetState(logisticKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(logisticsasbytes1)
}

func (t *SmartContract) QueryAllLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	_,mspid,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if len(args) != 1 {
		return shim.Error("Invalid   no of arg for Query function ")
	}
	holder := affiliation+args[0]
	resultsIterator, err := stub.GetStateByPartialCompositeKey("logisticstrans", []string{})
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
		logis := logisticstrans{}
		json.Unmarshal(queryResponse.Value,logis)
		if mspid != "DelivererMSP" {
			if !(logis.BuyerID == holder || logis.SellerID == holder) {
				continue
			}
		}else {
			if logis.LogisticstranID != holder {
				continue
			}
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

