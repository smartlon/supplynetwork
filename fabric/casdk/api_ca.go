package casdk

import "github.com/smartlon/supplynetwork/fabric/utils"

var CaClients map[string]*FabricCAClient
var AdminIdns map[string]*Identity

func init() {
	var err error
	CaClients = make(map[string]*FabricCAClient)
	AdminIdns = make(map[string]*Identity)
	err = NewCAClient("/home/lgao/go/src/github.com/smartlon/supplynetwork/fabric/casdk/caconfig.yaml", nil)
	if err != nil {
		panic(err)
	}
}

func EnrollUser(enrollmentId, secret, orgName string) (token, message string, success bool) {
	enrollReq := CaEnrollmentRequest{
		EnrollmentId: enrollmentId,
		Secret:       secret,
	}
	idn, err := CaClients[orgName].Enroll(enrollReq)
	AdminIdns[orgName]=idn
	if err != nil {
		return "",err.Error(),false
	}
	token,err  = utils.GenerateToken(enrollmentId,orgName)
	if err != nil {
		return "",err.Error(),false
	}
	message = enrollmentId+ " logined successfully"
	return token,message,true
}

func RegisterUser(enrollmentId, userType, secret, orgName string)(secretRes, message string, success bool){
	req := CARegistrationRequest{
		EnrolmentId:    enrollmentId,
		Type:           userType,
		Secret:         secret,
		MaxEnrollments: -1,
		Attrs:          nil,
		CAName:         CaClients[orgName].ServerInfo.CAName,
	}
	idn := AdminIdns[orgName]
	secretRes, err := CaClients[orgName].Register(idn, &req)
	if err != nil {
		return "",err.Error(),false
	}
	message = enrollmentId+ " registered successfully"
	return secretRes,message,true
}

func RevokeUser( cert, orgName string)(secretRes, message string, success bool) {
	serial, aki, err := CaClients[orgName].GetCertSerialAki([]byte(cert))
	if err != nil {
		return "",err.Error(),false
	}
	req := CARevocationRequest{
		//EnrollmentId: "ca5", // 根据注册用户注销其证书
		Serial: serial,
		AKI:    aki,
		GenCRL: true,
	}
	idn := AdminIdns[orgName]
	_, err = CaClients[orgName].Revoke(idn, &req)
	if err != nil {
		return "",err.Error(),false
	}
	message = serial+ " registered successfully"
	return secretRes,message,true
}

func GetAllUser(orgName string)(Identities []CaIdentityResponse,message string, success bool)  {
	idn := AdminIdns[orgName]
	caName := CaClients[orgName].ServerInfo.CAName
	allUser, err := CaClients[orgName].GetIdentities(idn,caName)
	if err != nil {
		return []CaIdentityResponse{},err.Error(),false
	}
	return allUser.Identities,"",true
}

