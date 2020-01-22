/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
	"github.com/smartlon/supplynetwork/iota/sdk"
	"strconv"
	"sync"
	"time"
)

type IotaPayload struct {
	ContainerID string `json:"ContainerID"`
	Seed        string `json:"Seed"`
	MamState    string `json:"MamState"`
	Root        string `json:"Root"`
	Mode       	string `json:"Mode"`
	SideKey     string `json:"SideKey"`
}

func  listener(action *chaincodeInvokeAction,chaincode string,wg *sync.WaitGroup) error {
	defer 	wg.Done()
	ec, err := action.EventClient(event.WithBlockEvents())
	if err != nil {
		fmt.Println("failed to create client")
		return err
	}

	registrationCreateChannel, notifierCreateChannel, err := ec.RegisterChaincodeEvent(chaincode, `{"From":"Fabric","To":"Iota","Func":"CreateChannel"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event: CreateChannel")
		return err
                
	}
	registrationDeliveryLogistics, notifierDeliveryLogistics, err := ec.RegisterChaincodeEvent(chaincode, `{"From":"Fabric","To":"Iota","Func":"DeliveryLogistics"}`)
	if err != nil {
		fmt.Println("failed to register chaincode event: DeliveryLogistics")
		return err
	}
	defer  unregister(ec,[]fab.Registration{registrationCreateChannel,registrationDeliveryLogistics})
	var iotapayload IotaPayload

	select {
	case ccEvent,ok := <-notifierCreateChannel:
		if !ok {
			return errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event")
		}
		err := json.Unmarshal(ccEvent.Payload,&iotapayload)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("received chaincode event CreateChannel:  %v\n", iotapayload)
		timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
		iotdata := &sdk.IoTData{
			iotapayload.ContainerID,
			"",
			"",
			timestamp,
			"start",
		}
		iotdatabytes,_ := json.Marshal(iotdata)
		fmt.Println(iotdatabytes)
		//mamstat, root:= sdk.MAMTransmit(string(iotdatabytes),iotapayload.Seed,iotapayload.Mode,iotapayload.SideKey,"FABRIC")
		//
		//var argsArray []Args
		//argsArray = append(argsArray, Args{"InTransitLogistics",[]string{iotapayload.ContainerID,root,mamstat}})
		//_, err = action.invoke(Config().ChannelID, chaincode, argsArray)
		//if err !=nil {
		//	fmt.Printf("fabric  failed to callback for InTransitLogistics :  %v\n", err.Error())
		//}
	case ccEvent,ok := <-notifierDeliveryLogistics:
		if !ok {
			return errors.WithMessage(err, "unexpected closed channel while waiting for chaincode event")
		}
		err := json.Unmarshal(ccEvent.Payload,&iotapayload)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("received chaincode event DeliveryLogistics:  %v\n", iotapayload)
		timestamp := strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
		iotdata := &sdk.IoTData{
			iotapayload.ContainerID,
			"",
			"",
			timestamp,
			"end",
		}
		iotdatabytes,_ := json.Marshal(iotdata)
		fmt.Println(iotdatabytes)
		//sdk.MAMTransmit(string(iotdatabytes),iotapayload.Seed,iotapayload.Mode,iotapayload.SideKey,"FABRIC")
		//var argsArray []Args
		//iotdatas := sdk.MAMReceive(iotapayload.Root,iotapayload.Mode,iotapayload.SideKey)
		//var temperature string
		//for _,iotdata := range iotdatas {
		//	var iot sdk.IoTData
		//	err := json.Unmarshal([]byte(iotdata),&iot)
		//	if err != nil {
		//		fmt.Printf("fabric  failed to callback for SignLogistics :  %v\n", err.Error())
		//	}
		//	if len(iot.Temperature) != 0 {
		//		temperature = temperature + iot.Temperature + ","
		//	}
		//}
		//temperature = strings.TrimSuffix(temperature,",")
		//fmt.Println("SignLogistics temperature: ",temperature)
		//argsArray = append(argsArray, Args{"SignLogistics",[]string{iotapayload.ContainerID,temperature}})
		//_, err = action.invoke(Config().ChannelID, chaincode, argsArray)
		//if err !=nil {
		//	fmt.Printf("fabric  failed to callback for SignLogistics :  %v\n", err.Error())
		//}
	case <-time.After(time.Second * 3):
		fmt.Println("Exit while waiting for chaincode event")
		}
	return nil
}

func unregister(ec *event.Client,registrations []fab.Registration) {
	for _,registration := range registrations {
		ec.Unregister(registration)
	}
}






