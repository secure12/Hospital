#!/bin/bash
FILE="test.go"
if [ $# -lt 4 ]; then
    cat <<EOF
Usage:
    ./invoke.sh  cliName/peerName   channelName     chaincodeName   chaincodeFunction   chaincodeArguments

Example:
    ./invoke.sh  cli2               privatechannel  test            initLedger

Chaincode Functions:
EOF
    sed -n -e 's/^.*function.*==.*\"\(.*\)\".*$/  - \1/p' $FILE
    exit 1
fi
QUERY=\'{\"Args\":\[\"${@:4:1}\"$(printf ",\"%s\"" ${@:5})\]}\'
set -ev
eval docker exec $1 peer chaincode invoke -C $2 -n $3 -c $QUERY --tls --cafile /etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
