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

func VerifyToken(ctx *context.Context)(orgName, usernName string,err error) {
	if ctx.Input.Header("Authorization") != "" {
		authorization := ctx.Input.Header("Authorization")
		token := strings.Split(authorization, " ")[1]
		log.Info("curernttoken: ", token)
		orgName, userName, err := utils.GetUserInfoFromValidateToken(token)
		if err != nil {
			return "","",err
		}
		return orgName,userName,nil
	}
	return "","",fmt.Errorf("Authorization is empty")
}

func (lc *LogisticsController) EnrollCA(){
	loginUserReqBytes := lc.Ctx.Input.RequestBody
	var loginUserReq UserReq
	err := json.Unmarshal(loginUserReqBytes,&loginUserReq)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "token": "","certpem":"","prikeypem":""}
		lc.ServeJSON()
	}
	token, msg,cert,priKey, success := casdk.EnrollCA(loginUserReq.UserName,loginUserReq.PassWord,loginUserReq.OrgName)
	lc.Data["json"] = map[string]interface{}{"success": success,"msg": msg, "token": token,"certpem":cert,"prikeypem":priKey}
	lc.ServeJSON()
}

func (lc *LogisticsController) GetAllUser(){
	orgName,_,err := VerifyToken(lc.Ctx)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "user": "","count": 0}
		lc.ServeJSON()
	}
	allUser, msg, success := casdk.GetAllUser(orgName)
	count := len(allUser)
	lc.Data["json"] = map[string]interface{}{"success": success,"msg": msg, "user": allUser,"count": count}
	lc.ServeJSON()
}

func (lc *LogisticsController) RevokeUser(){
	orgName,_,err := VerifyToken(lc.Ctx)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "revokedlist": "","crl":"","count": 0}
		lc.ServeJSON()
	}
	revokeUserReqBytes := lc.Ctx.Input.RequestBody
	var revokeUserReq UserReq
	err = json.Unmarshal(revokeUserReqBytes,&revokeUserReq)
	if err != nil {
		fmt.Println(err.Error())
		lc.Data["json"] = map[string]interface{}{"success": false,"msg": err.Error(), "revokedlist": "","crl":"","count": 0}
		lc.ServeJSON()
	}
	caRevokeResult, msg, success := casdk.RevokeUser(revokeUserReq.UserName,orgName)
	count := len(caRevokeResult.RevokedCertificates)
	lc.Data["json"] = map[string]interface{}{"success": success,"msg": msg, "revokedlist": caRevokeResult.RevokedCertificates,"crl":caRevokeResult.CRL,"count": count}
	lc.ServeJSON()
}

func (lc *LogisticsController) RegisterUser(){
	orgName,_,err := VerifyToken(lc.Ctx)
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