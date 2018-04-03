#!/bin/bash
FILE="$3.go"
if [ $# -lt 4 ]; then
    cat <<EOF
Usage:
    ./query.sh  cliName/peerName    channelName     chaincodeName   chaincodeFunction   chaincodeArguments

Example:
    ./query.sh  cli.ce.prihos1      privatechannel  private         getAll

Chaincode Functions:
EOF
    sed -n -e 's/^.*function.*==.*\"\(.*\)\".*$/  - \1/p' $FILE
    exit 1
fi
QUERY=\'{\"Args\":\[\"${@:4:1}\"$(printf ",\"%s\"" ${@:5})\]}\'
set -ev
eval docker exec $1 peer chaincode query -C $2 -n $3 -c $QUERY --tls --cafile /etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
