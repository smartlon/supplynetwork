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

type Participatant struct {
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

func (s *SmartContract) RecordParticipatant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	participatant := Participatant{UserName:args[0],Affiliation:args[1],Location:args[2]}
	participatantAsBytes,err := json.Marshal(participatant)
	if err != nil {
		return shim.Error(err.Error())
	}
	participatantID := args[1]+ args[0]
	err = stub.PutState(participatantID, participatantAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record participatant: %s", participatantID))
	}
	return shim.Success(nil)
}

func (s *SmartContract) QueryAllParticipatant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	startKey :=  ""
	endKey :=  ""
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

	fmt.Printf("- queryAllParticipatant:\n%s\n", buffer.String())
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
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	holder := affiliation + "." + id
	productID := holder + "." + args[0]
	product := Product{ProductID:args[0],ProductType:args[1],Description:args[2],Timestamp:args[3],Holder:holder,Status:"false"}

	productAsBytes, _ := json.Marshal(product)
	err = stub.PutState(productID, productAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record product: %s", productID))
	}
	return shim.Success(nil)
}

func (s *SmartContract) QueryAllProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id,_,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	startKey := affiliation + "." + id + "."  + "product1"
	endKey := affiliation + "." + id + "." + "product999"

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
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	holder := affiliation + "." + id
	container := Container{ ContainerID:args[0],Description: args[1], Location: args[2], Timestamp: args[3], Holder: holder,Used:"false" }
	containerID := holder + "."  + args[0]
	containerAsBytes, _ := json.Marshal(container)
	err = stub.PutState(containerID, containerAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record container: %s", containerID))
	}

	mode := "restricted"
	seed := args[4]
	sideKey := args[5]
	iotaID := containerID + "."  + "iotapayload"
	iotaPayload := IotaPayload{ContainerID:containerID,Seed: seed, MamState: "", Root: "", Mode: mode, SideKey: sideKey}
	iotaPayloadAsBytes, _ := json.Marshal(iotaPayload)
	stub.PutState(iotaID, iotaPayloadAsBytes)
	fmt.Println("New Asset", args[0], container, seed, mode, sideKey)

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
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	containerID := affiliation + "." + id + "."  + args[0]
	containerAsBytes, _ := stub.GetState(containerID)
	if containerAsBytes == nil {
		return shim.Error("Could not locate container")
	}
	container := Container{}
	json.Unmarshal(containerAsBytes, &container)
	iotaID := containerID + "."  + "iotapayload"
	iotaPayloadAsBytes, _ := stub.GetState(iotaID)
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
	id,mspid,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if mspid != "DelivererMSP"{
		return shim.Error(fmt.Sprintf("the user %s is not belong to the deliverer",id))
	}
	startKey := affiliation + "." + id + "."  +"container1"
	endKey := affiliation + "." + id + "."  +"container999"

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
	case "RecordParticipatant":
		return t.RecordParticipatant(stub, args)
	case "QueryAllParticipatant":
		return t.QueryAllParticipatant(stub, args)
	case "RecordProduct":
		return t.RecordProduct(stub, args)
	case "QueryAllProduct":
		return t.QueryAllProduct(stub, args)
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
	id,_,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	productID := affiliation + "." + id + "."  + args[1]
	productAsBytes, _ := stub.GetState(productID)
	if productAsBytes == nil {
		return shim.Error("Could not locate product")
	}
	product := Product{}
	json.Unmarshal(productAsBytes, &product)
	product.Status = "true"
	productAsBytes, err = json.Marshal(product)
	err = stub.PutState(productID, productAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to deliver product: %s", args[1]))
	}
	sellerParticipatant := Participatant{}
	participatantAsBytes,err := stub.GetState(affiliation + "." + id)
	err = json.Unmarshal(participatantAsBytes,sellerParticipatant)
	if err != nil {
		return shim.Error(err.Error())
	}
	buyerParticipatant := Participatant{}
	participatantAsBytes,err = stub.GetState(args[2])
	err = json.Unmarshal(participatantAsBytes,buyerParticipatant)
	if err != nil {
		return shim.Error(err.Error())
	}
	logistiParticipatant := Participatant{}
	participatantAsBytes,err = stub.GetState(args[3])
	err = json.Unmarshal(participatantAsBytes,logistiParticipatant)
	if err != nil {
		return shim.Error(err.Error())
	}
	var logobj = logisticstrans{
		LogisticstranID: args[0],
		ProductID: args[1],
		ProductType: product.ProductType,
		BuyerID: args[2],
		BuyerLocation: buyerParticipatant.Location,
		SellerID: affiliation + "." + id,
		SellerLocation: sellerParticipatant.Location,
		LogisticsID: args[3],
		LogisticsLocation: logistiParticipatant.Location,
	}
	logobj.Status = "Requested"
	logisticsID := args[0]
	logobjasBytes, _ := json.Marshal(logobj)
	err = stub.PutState(logisticsID, logobjasBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to request logistics: %s", logisticsID))
	}
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
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 3")
	}
	hoder := affiliation + "." + id
	logisticsAsBytes, _ := stub.GetState(args[0])
	var logisticobj logisticstrans
	json.Unmarshal(logisticsAsBytes, &logisticobj)
	logisticobj.MAMChannel.ContainerID = hoder + "." + args[1]
	sideKey := args[2]
	logisticobj.MAMChannel.SideKey = sideKey
	timestamp := args[3]
	logisticobj.JourneyStartTime = timestamp

	containerAsBytes, _ := stub.GetState(hoder + "." + args[1])
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
	container.Location = logisticobj.SellerLocation
	iotaPayloadAsBytes, _ := stub.GetState(hoder + "."+container.ContainerID + "." +"iotapayload" )
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
	stub.PutState(hoder + "." + container.ContainerID, containerAsBytes)
	iotaPayloadAsBytes, _ = json.Marshal(iotaPayload)
	stub.PutState(hoder + "."+container.ContainerID + "." +"iotapayload", iotaPayloadAsBytes)
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
	iotaID := logisticobj.MAMChannel.ContainerID + "."  + "iotapayload"
	iotaPayloadAsBytes, _ := stub.GetState(iotaID)
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
	stub.PutState(iotaID, iotaPayloadAsBytes)
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
	iotaPayloadAsBytes, _ := stub.GetState(logisticobj1.MAMChannel.ContainerID + "." + "iotapayload" )
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
		product := Product{}
		productAsBytes,_ := stub.GetState(logisticobj1.SellerID + "." + logisticobj1.ProductID)
		json.Unmarshal(productAsBytes,product)
		product.Status = "false"
		productAsBytes,_ = json.Marshal(product)
		stub.PutState(logisticobj1.BuyerID + "." + logisticobj1.ProductID,productAsBytes)
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
	id,mspid,affiliation , err := ABAC(stub)
	if err != nil {
		return shim.Error("There was an error trying to retrieve the attribute")
	}
	holder := affiliation + "." + id
	startKey := ""
	endKey := ""
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

