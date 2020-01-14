
//working invoke

//open and try this commnds in cli

//docker exec -it cli bash

CC_NAME=logistic

peer chaincode invoke -o orderer.logisticstransfer.com:7050  --tls --cafile  /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/logisticstransfer.com/orderers/orderer.logisticstransfer.com/msp/tlscacerts/tlsca.logisticstransfer.com-cert.pem  -C "logchannel" -n $CC_NAME -c '{"Args":["RequestLogistic","product2","medical","seller1","xian","buyer1","beijing"]}'
"product1","medical","seller1","xian","buyer1","beijing"
//working invoke
peer chaincode invoke -o orderer.logisticstransfer.com:7050  --tls --cafile  /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/logisticstransfer.com/orderers/orderer.logisticstransfer.com/msp/tlscacerts/tlsca.logisticstransfer.com-cert.pem  -C "logchannel" -n $CC_NAME -c '{"Args":["RequestLogistic","{\"Name\":\"Dinesh\",\"ID\":\"2\"}"]}'

//query by id
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile  /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C "logchannel" -n $CC_NAME -c '{"Args":["QueryLogistics","product1"]}'

/home/lgao/go/src/github.com/smartlon/gateway/adapter/fabric/network/fixtures/config