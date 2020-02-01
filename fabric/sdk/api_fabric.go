package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartlon/supplynetwork/fabric/utils"
	"github.com/smartlon/supplynetwork/log"
	"strings"
)

var action *chaincodeInvokeAction
var err error
func init() {
	action, err = newChaincodeInvokeAction()
	if err != nil {
		log.Infof(err.Error())
	}
}

func EnrollUser(enrollmentId, secret, orgName string) (token, message string, success bool) {
	_, err:= action.newUser(orgName,enrollmentId,secret)
	if err != nil {
		return "",err.Error(),false
	}
	token,err  = utils.GenerateToken(enrollmentId,orgName)
	if err != nil {
		return "",err.Error(),false
	}
	message = enrollmentId+ " logined successfully"

	return token,message ,true
}

// ChaincodeInvoke invoke chaincode
func ChaincodeInvoke(chaincodeID string, argsArray []Args,orgName,userName string) (result string, err error)  {
	log.Info("chaincode invoke...")
	if chaincodeID == "" {
		err = fmt.Errorf("must specify the chaincode ID")
		return "" , err
	}
	action.Set(Config().ChannelID,chaincodeID,[]Args{})
	if err != nil {
		log.Errorf("Error while initializing invokeAction: %v", err)
		return "" , err
	}
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go listener(action,chaincodeID,&wg)

	result, err = action.invoke(Config().ChannelID, chaincodeID, argsArray,orgName,userName)
	if err != nil {
		log.Errorf("Error while calling action.invoke(): %v", err)
	}
	//wg.Wait()
	return result, err
}

// ChaincodeInvokeByString call chaincode invoke of hyperledger fabric
func ChaincodeInvokeByString(chaincodeID, argsStr string) string {
	log.Infof("chaincode invoke: %s; %s; %s", chaincodeID, argsStr)
	arg, err := ParseArgs(argsStr)
	//if strings.EqualFold(arg.Func, "create") {
		//RegisterWalletByString(arg.Args[0], arg.Args[1], arg.Args[2])
	//}
	if err == nil {
		var argsArray []Args
		argsArray = append(argsArray, *arg)
		var ret string
		ret, err = ChaincodeInvoke(chaincodeID, argsArray)
		if err == nil {
			log.Info("chaincode invoke result: ", ret)
			return ret
		}
	}
	log.Errorf("%s %v", defaultResultJSON, err)
	return defaultResultJSON
}

// ChaincodeQuery call chaincode query of hyperledger fabric
func ChaincodeQuery(chaincodeID string, argsArray []Args) (result string, err error) {
	log.Info("chaincode query...")
	if chaincodeID == "" {
		err = fmt.Errorf("must specify the chaincode ID")
		return
	}
	var action *queryAction
	action, err = newQueryAction()
	if err != nil {
		log.Errorf("Error while initializing queryAction: %v", err)
		return
	}

	defer action.Terminate()

	result, err = action.query(Config().ChannelID, chaincodeID, argsArray)
	if err != nil {
		log.Errorf("Error while running queryAction: %v", err)
	} else if result == "" {
		err = errors.New("transaction not found")
	}
	if strings.HasPrefix(result, "v") {
		if i := strings.Index(result, ":"); i > -1 {
			result = result[i+1:]
		}
	}
	return
}

// ChaincodeQueryByString call chaincode query of hyperledger fabric
func ChaincodeQueryByString(chaincodeID, argsStr string) string {
	log.Infof("chaincode query: %s; %s; %s", chaincodeID, argsStr)
	argsArray, err := ArgsArray(argsStr)
	if err == nil {
		var ret string
		ret, err = ChaincodeQuery(chaincodeID, argsArray)
		if err == nil {
			log.Info("chaincode query result: ", ret)
			return ret
		}
	}
	log.Errorf("%s %v", defaultResultJSON, err)
	return defaultResultJSON
}

// TxInfo info of Tx
type TxInfo struct {
	Contract string `json:"contract,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Amount   string `json:"amount,omitempty"`
	GasUsed  string `json:"gasUsed,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	TxHash   string `json:"txHash,omitempty"`
	Height   string `json:"height,omitempty"`
	Status   string `json:"status,omitempty"`
}

// TxRegister registered Tx
type TxRegister struct {
	Key       string  `json:"key,omitempty"`
	ChainName string  `json:"chain,omitempty"`
	TokenName string  `json:"token,omitempty"`
	Addr      string  `json:"address,omitempty"`
	Amount    string  `json:"amount,omitempty"`
	GasUsed   string  `json:"gasUsed,omitempty"`
	GasPrice  string  `json:"gasPrice,omitempty"`
	Info      *TxInfo `json:"info,omitempty"`
}

// BlockRegister registered block
type BlockRegister struct {
	Height string        `json:"height,omitempty"`
	Txs    []*TxRegister `json:"transactions,omitempty"`
}

// RegisterBlock register block of source chain into hyperledger fabric
func RegisterBlock(block *BlockRegister) string {
	var err error
	var bytes []byte
	if bytes, err = json.Marshal(block); err != nil {
		log.Errorf("register block error: %v", err)
		return defaultResultJSON
	}
	a := Args{
		Func: "register",
		Args: []string{"block", string(bytes)}}
	var args []Args
	args = append(args, a)
	var ret string
	ret, err = ChaincodeInvoke("wallet", args)
	if err != nil {
		log.Errorf("register block error: %v", err)
		return defaultResultJSON
	}
	log.Info(ret)
	return ret
}



