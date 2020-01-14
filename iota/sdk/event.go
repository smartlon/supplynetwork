package sdk

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"strconv"
	"strings"
)

type Transaction struct {
	Type         string `json:"type"`
	Hash         string `json:"hash"`
	Address      string `json:"address"`
	Value        int    `json:"value"`
	ObsoleteTag  string `json:"obsolete_tag"`
	Timestamp    int64  `json:"timestamp"`
	CurrentIndex int    `json:"current_index"`
	LastIndex    int    `json:"last_index"`
	BundleHash   string `json:"bundle_hash"`
	TrunkTxHash  string `json:"trunk_tx_hash"`
	BranchTxHash string `json:"branch_tx_hash"`
	ArrivalTime  int64  `json:"arrival_time"`
	Tag          string `json:"tag"`
	Status       string `json:"status"`
	Inherent_lat int64  `json:"inherent_latency"`
	Confirm_lat  int64  `json:"confirm_latency"`
}

type ConfTx struct {
	Hash         string `json:"hash"`
	Address	     string `json:"address"`
	TrunkTxHash  string `json:"trunk_tx_hash"`
	BranchTxHash string `json:"branch_tx_hash"`
	BundleHash   string `json:"bundle_hash"`
}

func StartTxFeed() {
	socket, err := zmq4.NewSocket(zmq4.SUB)
	must(err)
	socket.SetSubscribe("tx")
	err = socket.Connect(zmqAddress)
	must(err)

	fmt.Printf("started tx feed\n")
	for {
		msg, err := socket.Recv(0)
		must(err)

		tx := buildTxFromZMQData(msg)
		//fmt.Printf("received tx: %s\n",tx)
		if tx == nil {
			fmt.Printf("tx: receive error! message format error\n")
			continue
		} else if tx.Type == "tx_trytes" {
			//fmt.Printf("tx: trytes received. Skip.\n")
			continue
		}
		if tx.Tag == chainTag {
			fmt.Printf("received tx: %s\n",tx)
		}

	}
}

func StartConfirmationFeed() {
	socket, err := zmq4.NewSocket(zmq4.SUB)
	must(err)
	socket.SetSubscribe("sn")
	err = socket.Connect(zmqAddress)
	must(err)

	fmt.Printf("started confirmation feed\n")
	for {
		msg, err := socket.Recv(0)
		must(err)
		tx := buildConfirmFromZMQData(msg)
		if tx == nil {
			fmt.Printf("confirm: receive error! message format error\n")
			continue
		}

	}
}

func buildTxFromZMQData(msg string) *Transaction {
	msgSplit := strings.Split(msg, " ")
	if len(msgSplit) != 13 {
		if msgSplit[0] == "tx_trytes" {
			return &Transaction{Type:msgSplit[0]}
		} else {
			return nil
		}
	}
	var err error
	tx := &Transaction{}
	tx.Type = msgSplit[0]
	tx.Hash = msgSplit[1]
	tx.Address = msgSplit[2]
	tx.Value, err = strconv.Atoi(msgSplit[3])
	if err != nil {
		return nil
	}
	tx.ObsoleteTag = msgSplit[4]
	tx.Timestamp, err = strconv.ParseInt(msgSplit[5], 10, 64)
	if err != nil {
		return nil
	}
	tx.CurrentIndex, err = strconv.Atoi(msgSplit[6])
	if err != nil {
		return nil
	}
	tx.LastIndex, err = strconv.Atoi(msgSplit[7])
	if err != nil {
		return nil
	}
	tx.BundleHash = msgSplit[8]
	tx.TrunkTxHash = msgSplit[9]
	tx.BranchTxHash = msgSplit[10]
	tx.ArrivalTime, err = strconv.ParseInt(msgSplit[11], 10, 64)
	if err != nil {
		return nil
	}
	tx.Tag = msgSplit[12]
	return tx
}

func buildConfirmFromZMQData(msg string) *ConfTx {
	msgSplit := strings.Split(msg, " ")
	if len(msgSplit) != 7 {
		return nil
	}
	msgSplit = msgSplit[2:]
	tx := &ConfTx{}
	tx.Hash = msgSplit[0]
	tx.Address = msgSplit[1]
	tx.TrunkTxHash = msgSplit[2]
	tx.BranchTxHash = msgSplit[3]
	tx.BundleHash = msgSplit[4]
	return tx
}
