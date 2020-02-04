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

func EnrollCA(enrollmentId, secret, orgName string) (token, message string,cert, priKey []byte, success bool) {
	enrollReq := CaEnrollmentRequest{
		EnrollmentId: enrollmentId,
		Secret:       secret,
	}
	idn, err := CaClients[orgName].Enroll(enrollReq)
	if enrollmentId == "admin" {
		AdminIdns[orgName]=idn
	}
	if err != nil {
		return "",err.Error(),nil,nil,false
	}
	token,err  = utils.GenerateToken(enrollmentId,orgName)
	if err != nil {
		return "",err.Error(),nil,nil,false
	}
	message = enrollmentId+ " logined successfully"
	cert = idn.GetPemCert()
	priKey,err = idn.GetPemPrivateKey()
	if err != nil {
		return "",err.Error(),nil,nil,false
	}
	return token,message,cert,priKey ,true
}

func RegisterUser(enrollmentId, userType, secret,affiliation string, orgName string)(secretRes, message string, success bool){
	req := CARegistrationRequest{
		EnrolmentId:    enrollmentId,
		Type:           userType,
		Secret:         secret,
		MaxEnrollments: -1,
		Affiliation:  affiliation,
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

func AddAffiliation(affiliation string, orgName string)(string,  string,  bool){
	req := affliationRequest{
		Name:affiliation,
		CaName:orgName,
	}
	idn := AdminIdns[orgName]
	affiliation, err := CaClients[orgName].AddAffiliations(idn, &req)
	if err != nil {
		return "",err.Error(),false
	}
	message := affiliation+ " added successfully"
	return affiliation,message,true
}

func RevokeUser( enrollmentId, orgName string)(caRevokeResult *CARevokeResult, message string, success bool) {
	//serial, aki, err := CaClients[orgName].GetCertSerialAki([]byte(cert))
	//if err != nil {
	//	return "",err.Error(),false
	//}
	req := CARevocationRequest{
		EnrollmentId: enrollmentId, // 根据注册用户注销其证书
		//Serial: serial,
		//AKI:    aki,
		GenCRL: true,
	}
	idn := AdminIdns[orgName]
	caRevokeResult, err := CaClients[orgName].Revoke(idn, &req)
	if err != nil {
		return &CARevokeResult{},err.Error(),false
	}
	message = enrollmentId+ " revoked successfully"
	return caRevokeResult,message,true
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

