#!/bin/bash
export IMAGE_TAG=latest
export CHANNEL_NAME=logchannel
echo "Generating cryto material for peers..."
[ -d ./supply-network/crypto-config ] || mkdir ./supply-network/crypto-config

cryptogen generate --config=./supply-network/crypto-config.yaml --output="./supply-network/crypto-config"

[ -d ./supply-network/channel-artifacts ] || mkdir ./supply-network/channel-artifacts

echo "Generating channel artifacts and genesis block..."
configtxgen -configPath ./supply-network -profile SupplyOrdererGenesis -outputBlock ./supply-network/channel-artifacts/genesis.block
configtxgen -configPath ./supply-network -profile SupplyChannel -outputCreateChannelTx ./supply-network/channel-artifacts/channel.tx -channelID $CHANNEL_NAME
configtxgen -configPath ./supply-network -profile SupplyChannel -outputAnchorPeersUpdate ./supply-network/channel-artifacts/ProducerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg ProducerMSP
configtxgen -configPath ./supply-network -profile SupplyChannel -outputAnchorPeersUpdate ./supply-network/channel-artifacts/ManufacturerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg ManufacturerMSP
configtxgen -configPath ./supply-network -profile SupplyChannel -outputAnchorPeersUpdate ./supply-network/channel-artifacts/DelivererMSPanchors.tx -channelID $CHANNEL_NAME -asOrg DelivererMSP
configtxgen -configPath ./supply-network -profile SupplyChannel -outputAnchorPeersUpdate ./supply-network/channel-artifacts/RetailerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg RetailerMSP

 CURRENT_DIR=$PWD
 cd ./supply-network/base
 cp docker-compose-base-template.yaml docker-compose-base.yaml
 OPTS="-i"
 cd $CURRENT_DIR
 cd ./supply-network/crypto-config/peerOrganizations/producer.example.com/ca/
 PRIV_KEY=$(ls *_sk)
 cd $CURRENT_DIR
 cd ./supply-network/base
 sed $OPTS "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-base.yaml
 
 cd $CURRENT_DIR
 cd ./supply-network/crypto-config/peerOrganizations/manufacturer.example.com/ca/
 PRIV_KEY=$(ls *_sk)
 cd $CURRENT_DIR
 cd ./supply-network/base
 sed $OPTS "s/CA2_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-base.yaml
 
 
 cd $CURRENT_DIR
 cd ./supply-network/crypto-config/peerOrganizations/deliverer.example.com/ca/
 PRIV_KEY=$(ls *_sk)
 cd $CURRENT_DIR
 cd ./supply-network/base
 sed $OPTS "s/CA3_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-base.yaml
 
 
 cd $CURRENT_DIR
 cd ./supply-network/crypto-config/peerOrganizations/retailer.example.com/ca/
 PRIV_KEY=$(ls *_sk)
 cd $CURRENT_DIR
 cd ./supply-network/base
 sed $OPTS "s/CA4_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-base.yaml
  
  
 cd $CURRENT_DIR
 ./supply-network/scripts/ccp-generate.sh
