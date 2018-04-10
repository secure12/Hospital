#!/bin/bash
if [ $# -lt 4 ]; then
    cat <<EOF
Usage:
    ./invoke.sh  cliName/peerName   channelName     chaincodeName   chaincodeFunction   chaincodeArguments

Example:
    ./invoke.sh  cli.ce.prihos1     privatechannel  private         initLedger

Chaincode Functions:
EOF
    for FILE in "joint/joint.go" "disjoint/disjoint.go" "github.com/hyperledger/fabric/examples/chaincode/go/enccc_example/enccc_example.go"; do
        echo $FILE
        sed -n -e 's/^.*case.*\"\(.*\)\".*$/  - \1/p' chaincode/$FILE
    done
    exit 1
fi
i=1
while [ $i -le $# ]; do
    if [[ "${!i}" =~ ^'-'.* ]]; then
        break
    fi
    ((i++))
done
echo $i
QUERY=\'{\"Args\":\[\"${@:4:1}\"$(printf ",\"%s\"" ${@:5:i-5})\]}\'
set -ev
eval docker exec $1 peer chaincode invoke -C $2 -n $3 -c $QUERY --tls --cafile /etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem ${@:i}
