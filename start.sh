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
export FABRIC_START_TIMEOUT=15
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec icu.pubhos1.example.com peer channel create -o orderer.example.com:7050 -c publicchannel -f publicchannel.tx --tls --cafile /etc/hyperledger/orderers/tlscacerts/tlsca.example.com-cert.pem
docker exec icu.prihos1.example.com peer channel create -o orderer.example.com:7050 -c privatechannel -f privatechannel.tx --tls --cafile /etc/hyperledger/orderers/tlscacerts/tlsca.example.com-cert.pem
docker exec icu.pubhos1.example.com peer channel create -o orderer.example.com:7050 -c commonchannel -f commonchannel.tx --tls --cafile /etc/hyperledger/orderers/tlscacerts/tlsca.example.com-cert.pem

# Join peer0.org1.example.com to the channel.
docker exec icu.pubhos1.example.com peer channel join -b publicchannel.block
docker exec surgery.pubhos1.example.com peer channel join -b publicchannel.block
docker exec ir.pubhos1.example.com peer channel join -b publicchannel.block
docker exec icu.pubhos2.example.com peer channel join -b publicchannel.block
docker exec surgery.pubhos2.example.com peer channel join -b publicchannel.block
docker exec ir.pubhos2.example.com peer channel join -b publicchannel.block

docker exec icu.prihos1.example.com peer channel join -b privatechannel.block
docker exec surgery.prihos1.example.com peer channel join -b privatechannel.block
docker exec ir.prihos1.example.com peer channel join -b privatechannel.block

docker exec icu.pubhos1.example.com peer channel join -b commonchannel.block
docker exec surgery.pubhos1.example.com peer channel join -b commonchannel.block
docker exec ir.pubhos1.example.com peer channel join -b commonchannel.block
docker exec icu.pubhos2.example.com peer channel join -b commonchannel.block
docker exec surgery.pubhos2.example.com peer channel join -b commonchannel.block
docker exec ir.pubhos2.example.com peer channel join -b commonchannel.block
docker exec icu.prihos1.example.com peer channel join -b commonchannel.block
docker exec surgery.prihos1.example.com peer channel join -b commonchannel.block
docker exec ir.prihos1.example.com peer channel join -b commonchannel.block
