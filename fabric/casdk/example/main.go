package main

import (
	"fmt"
	"github.com/smartlon/supplynetwork/fabric/casdk"
	"os"
)

func main() {
	//	var idresp casdk.CAGetIdentityResponse
	var idsresp casdk.CAListAllIdentitesResponse
	//初始化CAClient
	err := casdk.InitCASDK("./", "caconfig.yaml")
	if err != nil {
		fmt.Println(err)
	}
	//注册admin证书
	enrollRequest := casdk.CaEnrollmentRequest{EnrollmentId: "admin", Secret: "adminpw"}
	_, _, err = casdk.Enroll(casdk.CA, enrollRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//注册peer证书
		attr := []casdk.CaRegisterAttribute{{
			Name: "Revoker",
			Value: "true",
			ECert: true,
		},
		}
		rr := casdk.CARegistrationRequest{
			EnrolmentId: "peer2",
			Affiliation: "org1.department1",
			Type: "peer",
			Attrs: attr,
		}
		err = casdk.Register(casdk.CA, casdk.ID, &rr)
	
		if err != nil {
			fmt.Println(err)
		}
	//撤销证书
		req := casdk.CARevocationRequest{EnrollmentId: "peer1", Reason: "aacompromise", GenCRL: true}
		casdk.Revoke(casdk.CA, casdk.ID, &req)
	//查询单一id
		idresp, err := casdk.GetIndentity(casdk.CA, casdk.ID, "peer1", "ca1")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(idresp)
	//查询所有id
	idsresp, err = casdk.GetIndentities(casdk.CA, casdk.ID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(idsresp)
}
