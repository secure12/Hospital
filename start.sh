#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml up -d

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=5
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# toggle silent
SLIENT=

# Create the channel
docker exec $SLIENT icu.pubhos1.example.com peer channel create -o orderer.example.com:7050 -c publicchannel -f publicchannel.tx --tls --cafile /etc/hyperledger/orderers/tlscacerts/tlsca.example.com-cert.pem
docker exec $SLIENT icu.prihos1.example.com peer channel create -o orderer.example.com:7050 -c privatechannel -f privatechannel.tx --tls --cafile /etc/hyperledger/orderers/tlscacerts/tlsca.example.com-cert.pem
docker exec $SLIENT icu.pubhos1.example.com peer channel create -o orderer.example.com:7050 -c commonchannel -f commonchannel.tx --tls --cafile /etc/hyperledger/orderers/tlscacerts/tlsca.example.com-cert.pem

# Join peer0.org1.example.com to the channel.
docker exec $SLIENT icu.pubhos1.example.com peer channel join -b publicchannel.block
#docker exec $SLIENT surgery.pubhos1.example.com peer channel join -b publicchannel.block
#docker exec $SLIENT ir.pubhos1.example.com peer channel join -b publicchannel.block
docker exec $SLIENT icu.pubhos2.example.com peer channel join -b publicchannel.block
#docker exec $SLIENT surgery.pubhos2.example.com peer channel join -b publicchannel.block
#docker exec $SLIENT ir.pubhos2.example.com peer channel join -b publicchannel.block

docker exec $SLIENT icu.prihos1.example.com peer channel join -b privatechannel.block
#docker exec $SLIENT surgery.prihos1.example.com peer channel join -b privatechannel.block
#docker exec $SLIENT ir.prihos1.example.com peer channel join -b privatechannel.block

docker exec $SLIENT icu.pubhos1.example.com peer channel join -b commonchannel.block
#docker exec $SLIENT surgery.pubhos1.example.com peer channel join -b commonchannel.block
#docker exec $SLIENT ir.pubhos1.example.com peer channel join -b commonchannel.block
docker exec $SLIENT icu.pubhos2.example.com peer channel join -b commonchannel.block
#docker exec $SLIENT surgery.pubhos2.example.com peer channel join -b commonchannel.block
#docker exec $SLIENT ir.pubhos2.example.com peer channel join -b commonchannel.block
docker exec $SLIENT icu.prihos1.example.com peer channel join -b commonchannel.block
#docker exec $SLIENT surgery.prihos1.example.com peer channel join -b commonchannel.block
#docker exec $SLIENT ir.prihos1.example.com peer channel join -b commonchannel.block

SLIENT=
