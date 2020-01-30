package casdk

import "github.com/smartlon/supplynetwork/fabric/utils"

var CaClients map[string]FabricCAClient

func init() {
	var err error
	CaClients, err = NewCAClient("./caconfig.yaml", nil)
	if err != nil {
		panic(err)
	}
}

func EnrollUser(enrollmentId, secret, orgName string) (token, message string, success bool) {
	enrollReq := CaEnrollmentRequest{
		EnrollmentId: enrollmentId,
		Secret:       secret,
	}
	_, err := CaClients[orgName].Enroll(enrollReq)
	if err != nil {
		return "",err.Error(),false
	}
	token,err  = utils.GenerateToken(enrollmentId,orgName)
	if err != nil {
		return "",err.Error(),false
	}
	message = enrollmentId+ "logined successfully"
	return token,message,true
}