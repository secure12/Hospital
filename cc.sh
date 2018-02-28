set -ev

# install chaincode
#docker exec cli peer chaincode install -n test -v 1.0 -p ./myprj/

# instantiate chaincode
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -n test -v 1.0 -C mychannel -c '{"Args":["init"]}' --tls --cafile /etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
