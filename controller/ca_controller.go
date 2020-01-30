package controller

import (
	"encoding/json"
	"fmt"
	"github.com/smartlon/supplynetwork/fabric/casdk"
	"github.com/smartlon/supplynetwork/log"
)

type UserReq struct {
	UserName string `json:"UserName"`
	PassWord string `json:"PassWord"`
	OrgName string `json:"OrgName"`
}

func (lc *LogisticsController) EnrollUser(){
	registerUserReqBytes := lc.Ctx.Input.RequestBody
	var loginUserReq UserReq
	err := json.Unmarshal(registerUserReqBytes,&loginUserReq)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "token": ""}
		lc.ServeJSON()
	}
	token, msg, success := casdk.EnrollUser(loginUserReq.UserName,loginUserReq.PassWord,loginUserReq.OrgName)
	lc.Data["json"] = map[string]interface{}{"success": success,"msg": msg, "token": token}
	lc.ServeJSON()
}

func (lc *LogisticsController) RegisterUser(){
	registerUserReqBytes := lc.Ctx.Input.RequestBody
	var registerUserReq UserReq
	err := json.Unmarshal(registerUserReqBytes,&registerUserReq)
	if err != nil {
		fmt.Println(err.Error())
	}
	var ret string
	//ret, err = sdk.RegisterUser(registerUserReq.UserName, registerUserReq.PassWord)
	message := "register user success"
	if err != nil {
		log.Error(err)
		message = "register user fail"
	}
	lc.Data["json"] = map[string]interface{}{"code": 200,"message": message, "result": ret}
	lc.ServeJSON()
}