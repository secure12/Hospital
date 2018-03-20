FILE="test.go"
if [ "$2" != "invoke" ] && [ "$2" != "query" ]; then
    cat <<EOF
Usage:
    ./query.sh  invoke  cliName/peerName    channelName chaincodeName   chaincodeFunction   chaincodeArguments
    ./query.sh  query   cliName/peerName    channelName chaincodeName   chaincodeFunction   chaincodeArguments

Chaincode Functions:
EOF
    sed -n -e 's/^.*function.*==.*\"\(.*\)\".*$/  - \1/p' $FILE
    exit 1
fi
QUERY=\'{\"Args\":\[\"${@:5:1}\"$(printf ",\"%s\"" ${@:6})\]}\'
set -ev
eval docker exec $1 peer chaincode $2 -C $3 -n $4 -c $QUERY --tls --cafile /etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
