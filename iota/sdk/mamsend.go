package sdk

import (
	"log"

	"github.com/iotaledger/iota.go/mam/v1"
)

func PublishAndReturnState(message string, useTransmitter bool, seedFromStorage string, mamStateFromStorage string, mode string, sideKey string , transactionTag string) (string, string) {
	var t *mam.Transmitter = nil
	if useTransmitter == true {
		mamState := StringToMamState(mamStateFromStorage)
		t = ReconstructTransmitter(seedFromStorage, mamState)
	}

	transmitter, root := Publish(message, t, seedFromStorage, mode, sideKey,transactionTag)
	channel := transmitter.Channel()

	return MamStateToString(channel), root
}

func Publish(message string, t *mam.Transmitter, seed string, mode string, sideKey string, transactionTag string) (*mam.Transmitter, string) {
	transmitter := GetTransmitter(t,seed, mode, sideKey)

	root, err := transmitter.Transmit(message, transactionTag)
	if err != nil {
		log.Fatal(err)
	}

	return transmitter, root
}
