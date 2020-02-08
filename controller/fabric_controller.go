package controller

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/smartlon/supplynetwork/fabric/sdk"
	"github.com/smartlon/supplynetwork/log"
	"math/big"
	"strconv"
	"time"
)

const (
	CHAINCODEID = "supcc"

)

func (lc *LogisticsController) EnrollUser(){
	loginUserReqBytes := lc.Ctx.Input.RequestBody
	var loginUserReq UserReq
	err := json.Unmarshal(loginUserReqBytes,&loginUserReq)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "token": ""}
		lc.ServeJSON()
	}
	token, msg, success := sdk.EnrollUser(loginUserReq.UserName,loginUserReq.PassWord,loginUserReq.OrgName)
	lc.Data["json"] = map[string]interface{}{"success": success,"msg": msg, "token": token}
	lc.ServeJSON()
}

func (lc *LogisticsController) RecordProduct(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	recordProductReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(recordProductReqBytes,orgName ,userName)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
type ProductQueryResponse struct {
	Record    []Product `json:"result"`
}
type Product struct {
	ProductID string `json:"ProductID"`
	ProductType string `json:"ProductType"`
	Description string `json:"Description"`
	Timestamp string `json:"Timestamp"`
	Status string `json:"Status"`
	Holder  string `json:"Holder"`
}

func (lc *LogisticsController) QueryAllProduct(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	queryAllProductReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryAllProductReqBytes,orgName ,userName)
	var qr ProductQueryResponse
	err = json.Unmarshal([]byte(ret),&qr)
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(qr.Record)
	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": qr.Record}
	lc.ServeJSON()
}

func (lc *LogisticsController) RequestLogistic(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	requestLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(requestLogisticReqBytes,orgName ,userName)
        
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) TransitLogistics(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	transitLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(transitLogisticReqBytes,orgName ,userName)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) DeliveryLogistics(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	deliveryLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(deliveryLogisticReqBytes,orgName ,userName)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}

func (lc *LogisticsController) QueryLogistics(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	queryLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryLogisticReqBytes,orgName ,userName)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}


func (lc *LogisticsController) RecordContainer(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	recordContainerReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(recordContainerReqBytes,orgName ,userName)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}


func (lc *LogisticsController) QueryContainer(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	queryContainerReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryContainerReqBytes,orgName ,userName)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
type ContainerQueryResponse struct {
	Record    []Container `json:"result"`
}
type Container struct {
	ContainerID string `json:"ContainerID"`
	Description string `json:"Description"`
	Timestamp string `json:"Timestamp"`
	Location  string `json:"Location"`
	Used      string   `json:"Used"`  //true is used and false is not used
	Holder  string `json:"Holder"`
}

func (lc *LogisticsController) QueryAllContainers(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	queryAllContainersReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryAllContainersReqBytes,orgName ,userName)
	var qr ContainerQueryResponse
	err = json.Unmarshal([]byte(ret),&qr)
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(qr.Record)
	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": qr.Record}
	lc.ServeJSON()
}

type logisticstransQueryResponse struct {
	Record    []logisticstrans `json:"result"`
}

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
type MAMChannel struct {
	ContainerID       string       `json:"ContainerID"`
	Root              string       `json:"Root"`
	SideKey           string       `json:"SideKey"`
}
func (lc *LogisticsController) QueryAllLogistics(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	queryAllLogisticsReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryAllLogisticsReqBytes,orgName,userName)
	retBytes := bytes.Trim([]byte(ret),`\x00`)
	var qr logisticstransQueryResponse
	err = json.Unmarshal(retBytes,&qr)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	count := len(qr.Record)

	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": qr.Record}
	lc.ServeJSON()
}

type ParticipantQueryResponse struct {
	Record    []Participant `json:"result"`
}

type Participant struct {
	UserName string `json:"UserName"`
	Affiliation string `json:"Affiliation"`
	Location string `json:"Location"`
}

func (lc *LogisticsController) QueryParticipant(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	queryParticipantReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryParticipantReqBytes,orgName,userName)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}

func (lc *LogisticsController) QueryAllParticipant(){
	orgName,userName,err := VerifyToken(lc.Ctx)
	if err != nil {
		lc.Data["json"] = map[string]interface{}{"code": 201,"msg": err.Error(), "data": ""}
		lc.ServeJSON()
	}
	queryAllParticipantReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryAllParticipantReqBytes,orgName,userName)
	fmt.Println("QueryAllParticipant",[]byte(ret))
	var qr ParticipantQueryResponse

	err = json.Unmarshal([]byte(ret),&qr)
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(qr.Record)
	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": qr.Record}
	lc.ServeJSON()
}

func invokeController(invokeReqBytes []byte,orgName, userName string)(code int, message, ret string){
	var invokeReq sdk.Args
	err := json.Unmarshal(invokeReqBytes,&invokeReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	if invokeReq.Func == "RecordContainer" {
		timestamp := timeStamp()
		seed := generateRandomSeedString(81)
		sidekey := generateRandomSeedString(81)
		argscomposite := []string{timestamp,seed,sidekey}
		invokeReq.Args = append(invokeReq.Args,argscomposite...)
	}
	if invokeReq.Func == "TransitLogistics"  {
		timestamp := timeStamp()
		sidekey := generateRandomSeedString(81)
		argscomposite := []string{sidekey,timestamp}
		invokeReq.Args = append(invokeReq.Args,argscomposite...)
	}
	if  invokeReq.Func == "DeliveryLogistics" || invokeReq.Func == "RecordProduct"{
		timestamp := timeStamp()
		invokeReq.Args = append(invokeReq.Args,timestamp)
	}
	invokeReq.Args = append(invokeReq.Args,userName)
	var argsArray []sdk.Args
	argsArray = append(argsArray, invokeReq)

	ret, err = sdk.ChaincodeInvoke(CHAINCODEID, argsArray,orgName,userName)
	if err != nil {
		log.Error(err)
		message = err.Error()
		code = 201
	}else {
		message = "invoke " +invokeReq.Func+ " success"
		code = 200
	}
	return
}


func timeStamp() string {
	return strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
}

func generateRandomSeedString(length int) string {
	seed := ""
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ9"

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(27))
		if err != nil {
			fmt.Println(err)
		}
		seed += string(alphabet[n.Int64()])
	}
	return seed
}
