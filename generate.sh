#!/bin/sh
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
set -u
export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}

# remove previous crypto material and config transactions
rm -fr config
rm -fr crypto-config

# generate crypto material
cryptogen generate --config=./crypto-config.yaml --output=crypto-config
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

# generate genesis block for orderer
mkdir config
configtxgen -profile ThreeOrgOrdererGenesis -outputBlock ./config/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi

CHANNEL_NAME=publicchannel
# generate channel configuration transaction
configtxgen -profile PublicOrgChannel -outputCreateChannelTx ./config/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi
# generate anchor peer transaction
configtxgen -profile PublicOrgChannel -outputAnchorPeersUpdate ./config/PubHos1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg PubHos1MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for PubHos1MSP..."
  exit 1
fi
# generate anchor peer transaction
configtxgen -profile PublicOrgChannel -outputAnchorPeersUpdate ./config/PubHos2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg PubHos2MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for PubHos2MSP..."
  exit 1
fi

CHANNEL_NAME=privatechannel
# generate channel configuration transaction
configtxgen -profile PrivateOrgChannel -outputCreateChannelTx ./config/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi
# generate anchor peer transaction
configtxgen -profile PrivateOrgChannel -outputAnchorPeersUpdate ./config/PriHos1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg PriHos1MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for PriHos1MSP..."
  exit 1
fi

CHANNEL_NAME=commonchannel
# generate channel configuration transaction
configtxgen -profile CommonOrgChannel -outputCreateChannelTx ./config/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi
# generate anchor peer transaction
configtxgen -profile CommonOrgChannel -outputAnchorPeersUpdate ./config/PubHos1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg PubHos1MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for PubHos1MSP..."
  exit 1
fi
# generate anchor peer transaction
configtxgen -profile CommonOrgChannel -outputAnchorPeersUpdate ./config/PubHos2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg PubHos2MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for PubHos2MSP..."
  exit 1
fi
# generate anchor peer transaction
configtxgen -profile CommonOrgChannel -outputAnchorPeersUpdate ./config/PriHos1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg PriHos1MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for PriHos1MSP..."
  exit 1
fi

KEY=$(cd crypto-config/peerOrganizations/pubhos1.example.com/ca/ && ls *_sk) ################## DONT KNOW
sed -i "/FABRIC_CA_SERVER_CA_KEYFILE/c\      - FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/$KEY" docker-compose.yml
