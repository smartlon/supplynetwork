package casdk

import (
	"encoding/json"
	"fmt"
	"github.com/pretty66/casdk"
)

var caClients map[string]casdk.FabricCAClient
var err error

var filepath string = `D:\go\src\github.com\pretty66\casdk\caconfig.yaml`

func main() {

	caClients, err = casdk.NewCAClient("./caconfig.yaml", nil)
	if err != nil {
		panic(err)
	}
	for k,client := range caClients {
		res, err := client.GetCaInfo()
		if err != nil {
			fmt.Println(err.Error())
		}
		resB,_ := json.Marshal(res)
		fmt.Printf("key=%s, value=%s\n",k,string(resB))
	}
}