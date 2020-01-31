package controller

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	. "github.com/smartlon/supplynetwork/iota/sdk"
)

type LogisticsController struct {
	beego.Controller
}


type MAMTransmitReq struct {
	Message IoTData `json:"Message"`
	Seed string     `json:"Seed"`
	SideKey string  `json:"SideKey"`
}
type MAMReceiveReq struct {
	Root string `json:"Root"`
	SideKey string `json:"SideKey"`
}

func (lc *LogisticsController) MAMTransmit(){
	mamReqBytes := lc.Ctx.Input.RequestBody
	var mamReq MAMTransmitReq
	err := json.Unmarshal(mamReqBytes,&mamReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	var data IoTData
	data = mamReq.Message
	iotDataBytes,err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	_,root := MAMTransmit(string(iotDataBytes),mamReq.Seed,"restricted",mamReq.SideKey,"")
	var code int
	var message string
	var ret string
	if err != nil {
		code = 201
		message = "failed to transmit mam tx"
		ret = err.Error()
	}else {
		code = 200
		message = "successed to transmit mam tx"
		ret = root
	}

	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()

}


func (lc *LogisticsController) MAMReceive(){
	mamReqBytes := lc.Ctx.Input.RequestBody
	var mamReq MAMReceiveReq
	err := json.Unmarshal(mamReqBytes,&mamReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	mamMessages := MAMReceive(mamReq.Root,"restricted",mamReq.SideKey)
	var code int
	var message string
	var ret []string
	if err != nil {
		code = 201
		message = "failed to receive mam tx"
		ret = append(ret,err.Error())
	}else {
		code = 200
		message = "successed to receive mam tx"
		ret = mamMessages
	}
	count := len(ret)
	lc.Data["json"] = map[string]interface{}{"code": code,"count": count,"msg": message, "data": ret}
	lc.ServeJSON()
}

func (lc *LogisticsController) GetNodeInfo(){
	nodeInfoBytes,err := NodeInfo()
	var code int
	var message,ret string
	if err != nil {
		code = 201
		message = "failed to get node info"
		ret = err.Error()
	}else {
		code = 200
		message = "successed to get node info"
		ret = string(nodeInfoBytes)
	}

	lc.Data["json"] = map[string]interface{}{"code": code,"msg": message, "data": ret}
	lc.ServeJSON()
}
