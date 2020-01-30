package utils

import (
	"github.com/astaxie/beego/context"
	"github.com/smartlon/supplynetwork/log"
	"strings"
)

var FilterToken = func(ctx *context.Context) {
	log.Info("current router path is ", ctx.Request.RequestURI)
	if ctx.Request.RequestURI != "/login" && ctx.Input.Header("Authorization") == "" {
		log.Error("without token, unauthorized !!")
		ctx.ResponseWriter.WriteHeader(401)
		ctx.ResponseWriter.Write([]byte("no permission"))
	}
	if ctx.Request.RequestURI != "/login" && ctx.Input.Header("Authorization") != "" {
		token := ctx.Input.Header("Authorization")
		token = strings.Split(token, "")[1]
		log.Info("curernttoken: ", token)
		// validate token
		// invoke ValidateToken in utils/token
		// invalid or expired todo res 401
		err := ValidateToken(token)
		if err != nil {
			ctx.ResponseWriter.WriteHeader(401)
			ctx.ResponseWriter.Write([]byte("token is invalid"))
		}

	}
}
