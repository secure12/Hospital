CAFILE=/etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
VERSION="$(shell date +%s)"
test: test.go
	set -ev
	# install chaincode
	mkdir -p volumes/client
	sudo cp test.go volumes/client
	docker exec cli peer chaincode install -n test -v $(VERSION) -p myprj
	# instantiate chaincode
	docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -n test -v $(VERSION) -C publicchannel -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
	./invoke.sh cli publicchannel test initLedger
	docker exec cli2 peer chaincode install -n test -v $(VERSION) -p myprj
	# instantiate chaincode
	docker exec cli2 peer chaincode instantiate -o orderer.example.com:7050 -n test -v $(VERSION) -C privatechannel -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
	./invoke.sh cli2 privatechannel test initLedger
