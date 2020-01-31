package controller

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/smartlon/supplynetwork/fabric/casdk"
	"github.com/smartlon/supplynetwork/fabric/utils"
	"github.com/smartlon/supplynetwork/log"
	"strings"
)

type UserReq struct {
	UserName string `json:"UserName"`
	PassWord string `json:"PassWord"`
	UserType string `json:"UserType"`
	OrgName string `json:"OrgName"`
}

func verifyToken(ctx *context.Context)(orgName string,err error) {
	if ctx.Input.Header("Authorization") != "" {
		authorization := ctx.Input.Header("Authorization")
		token := strings.Split(authorization, " ")[1]
		log.Info("curernttoken: ", token)
		orgName, err := utils.GetOrgNameFromValidateToken(token)
		if err != nil {
			return "",err
		}
		return orgName,nil
	}
	return "",fmt.Errorf("Authorization is empty")
}

func (lc *LogisticsController) EnrollUser(){
	loginUserReqBytes := lc.Ctx.Input.RequestBody
	var loginUserReq UserReq
	err := json.Unmarshal(loginUserReqBytes,&loginUserReq)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "token": ""}
		lc.ServeJSON()
	}
	token, msg, success := casdk.EnrollUser(loginUserReq.UserName,loginUserReq.PassWord,loginUserReq.OrgName)
	lc.Data["json"] = map[string]interface{}{"success": success,"msg": msg, "token": token}
	lc.ServeJSON()
}

func (lc *LogisticsController) GetAllUser(){
	orgName,err := verifyToken(lc.Ctx)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "user": ""}
		lc.ServeJSON()
	}
	allUser, msg, success := casdk.GetAllUser(orgName)
	lc.Data["json"] = map[string]interface{}{"success": success,"msg": msg, "user": allUser}
	lc.ServeJSON()
}

func (lc *LogisticsController) RegisterUser(){
	orgName,err := verifyToken(lc.Ctx)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "secret": ""}
		lc.ServeJSON()
	}
	registerUserReqBytes := lc.Ctx.Input.RequestBody
	var registerUserReq UserReq
	err = json.Unmarshal(registerUserReqBytes,&registerUserReq)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "secret": ""}
		lc.ServeJSON()
	}
	secret, msg, success := casdk.RegisterUser(registerUserReq.UserName,registerUserReq.UserType,registerUserReq.PassWord,orgName)
	lc.Data["json"] = map[string]interface{}{"success": success,"msg": msg, "secret": secret}
	lc.ServeJSON()
}