package controller

import (
	"crypto/rand"
	"github.com/smartlon/supplynetwork/fabric/sdk"
	"encoding/json"
	"fmt"
	"github.com/smartlon/supplynetwork/log"
	"math/big"
	"strconv"
	"time"
)

const (
	CHAINCODEID = "logistic"

)

type UserReq struct {
	UserName string `json:"UserName"`
	PassWord string `json:"PassWord"`
}

func (lc *LogisticsController) RequestLogistic(){
	requestLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(requestLogisticReqBytes)
        
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) TransitLogistics(){
	transitLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(transitLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) DeliveryLogistics(){
	deliveryLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(deliveryLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) QueryLogistics(){
	queryLogisticReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryLogisticReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}

func (lc *LogisticsController) RecordContainer(){
	recordContainerReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(recordContainerReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
func (lc *LogisticsController) QueryContainer(){
	queryContainerReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryContainerReqBytes)
	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
type ContainerQueryResponse struct {
	Key    string `json:"Key"`
	Record    Container `json:"Record"`
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
	queryAllContainersReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryAllContainersReqBytes)
	var qr []ContainerQueryResponse
	err := json.Unmarshal([]byte(ret),&qr)
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(qr)
	var resp []Container
	for _,v := range qr {
		resp = append(resp,v.Record)
	}
	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": resp}
	lc.ServeJSON()
}

type logisticstransQueryResponse struct {
	Key    string `json:"Key"`
	Record    logisticstrans `json:"Record"`
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
	queryAllLogisticsReqBytes := lc.Ctx.Input.RequestBody
	code, message, ret := invokeController(queryAllLogisticsReqBytes)
	var qr []logisticstransQueryResponse
	err := json.Unmarshal([]byte(ret),&qr)
	if err != nil {
		fmt.Println(err.Error())
	}
	count := len(qr)
	var resp []logisticstrans
	for _,v := range qr {
		resp = append(resp,v.Record)
	}
	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": resp}
	lc.ServeJSON()
}


func invokeController(invokeReqBytes []byte)(code int, message, ret string){
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
	if  invokeReq.Func == "DeliveryLogistics" {
		timestamp := timeStamp()
		invokeReq.Args = append(invokeReq.Args,timestamp)
	}
	var argsArray []sdk.Args
	argsArray = append(argsArray, invokeReq)

	ret, err = sdk.ChaincodeInvoke(CHAINCODEID, argsArray)
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
