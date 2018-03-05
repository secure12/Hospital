VERSION="$(shell date +%s)"
test: test.go
	set -ev
	# install chaincode
	mkdir -p volumes/client
	sudo cp test.go volumes/client
	docker exec cli peer chaincode install -n test -v $(VERSION) -p myprj
	# instantiate chaincode
	docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -n test -v $(VERSION) -C mychannel -c '{"Args":["init"]}' --tls --cafile /etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

