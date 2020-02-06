package rooter

import (
	"github.com/astaxie/beego"
	"github.com/smartlon/supplynetwork/controller"
)

func init() {
	beego.Router("/iota/nodeinfo", &controller.LogisticsController{},"get:GetNodeInfo")
	beego.Router("/iota/mamtransmit", &controller.LogisticsController{},"post:MAMTransmit")
	beego.Router("/iota/mamreceive", &controller.LogisticsController{},"post:MAMReceive")

	beego.Router("/fabric/requestlogistic", &controller.LogisticsController{},"post:RequestLogistic")
	beego.Router("/fabric/transitlogistics", &controller.LogisticsController{},"post:TransitLogistics")
	beego.Router("/fabric/deliverylogistics", &controller.LogisticsController{},"post:DeliveryLogistics")
	beego.Router("/fabric/querylogistics", &controller.LogisticsController{},"post:QueryLogistics")
	beego.Router("/fabric/queryalllogistics", &controller.LogisticsController{},"post:QueryAllLogistics")

	beego.Router("/fabric/recordcontainer", &controller.LogisticsController{},"post:RecordContainer")
	beego.Router("/fabric/querycontainer", &controller.LogisticsController{},"post:QueryContainer")
	beego.Router("/fabric/queryallcontainers", &controller.LogisticsController{},"post:QueryAllContainers")

	//user management
	beego.Router("/queryallparticipants", &controller.LogisticsController{},"post:QueryAllParticipant")
	beego.Router("/registeruser", &controller.LogisticsController{},"post:RegisterUser")
	beego.Router("/enrolluser", &controller.LogisticsController{},"post:EnrollUser")
	beego.Router("/enrollca", &controller.LogisticsController{},"post:EnrollCA")
	beego.Router("/getalluser", &controller.LogisticsController{},"post:GetAllUser")
	beego.Router("/revokeuser", &controller.LogisticsController{},"post:RevokeUser")
	beego.Router("/addaffiliation", &controller.LogisticsController{},"post:AddAffiliation")

}
