package sdk

import (
	"encoding/json"
	"fmt"
)

type IoTData struct {
	ContainerID        string `json:"ContainerID"`
	Temperature string `json:"Temperature"`
	Location    string `json:"Location"`
	Time        string `json:"Time"`
	Status        string `json:"Status"`
}

func MAMTransmit(message string, seed string, mode string, sideKey string, transactionTag string) (string, string){
	transmitter, root := Publish(message, nil, seed, mode, sideKey,transactionTag)
	channel := transmitter.Channel()
	return MamStateToString(channel), root
}

func  MAMReceive( root string, mode string, sideKey string) ([]string ){
	channelMessages := Fetch(root,mode,sideKey)
	return channelMessages
}


func NodeInfo() ([]byte, error) {
	iotaApi := GetApi()
	nodeInfo, err := iotaApi.GetNodeInfo()
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}
	nodeInfoBytes,err := json.Marshal(nodeInfo)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}
	return nodeInfoBytes,nil
}
