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

# Create the channel
docker exec peer0.org1.example.com peer channel create -o orderer.example.com:7050 -c mychannel -f channel.tx --tls --cafile /etc/hyperledger/orderers/tlscacerts/tlsca.example.com-cert.pem
# Join peer0.org1.example.com to the channel.
docker exec peer0.org1.example.com peer channel join -b mychannel.block
docker exec peer1.org1.example.com peer channel join -b mychannel.block
docker exec peer0.org2.example.com peer channel join -b mychannel.block
docker exec peer1.org2.example.com peer channel join -b mychannel.block
