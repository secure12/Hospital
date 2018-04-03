CAFILE=/etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
VERSION="$(shell date +%s)"

all: test.go
	make test

test: test.go
	set -ev
	# install chaincode
	sudo mkdir -p volumes/pubhos1/hospital
	sudo cp test.go volumes/pubhos1/hospital
	#
	sudo mkdir -p volumes/pubhos1/enccc
	sudo cp -r enccc volumes/pubhos1/
	#
	sudo mkdir -p volumes/pubhos2/hospital
	sudo cp test.go volumes/pubhos2/hospital
	#
	sudo mkdir -p volumes/pubhos2/enccc
	sudo cp -r enccc volumes/pubhos2/
	#
	sudo mkdir -p volumes/prihos1/hospital
	sudo cp test.go volumes/prihos1/hospital
	#
	sudo mkdir -p volumes/prihos1/enccc
	sudo cp -r enccc volumes/prihos1/
	#
	docker exec cli.icu.pubhos1 peer chaincode install -n test -v $(VERSION) -p hospital
	docker exec cli.surgery.pubhos1 peer chaincode install -n test -v $(VERSION) -p hospital
	docker exec cli.ir.pubhos1 peer chaincode install -n test -v $(VERSION) -p hospital
	docker exec cli.icu.pubhos2 peer chaincode install -n test -v $(VERSION) -p hospital
	docker exec cli.surgery.pubhos2 peer chaincode install -n test -v $(VERSION) -p hospital
	docker exec cli.ir.pubhos2 peer chaincode install -n test -v $(VERSION) -p hospital
	docker exec cli.icu.prihos1 peer chaincode install -n test -v $(VERSION) -p hospital
	docker exec cli.surgery.prihos1 peer chaincode install -n test -v $(VERSION) -p hospital
	docker exec cli.ir.prihos1 peer chaincode install -n test -v $(VERSION) -p hospital
	#
	docker exec cli.icu.pubhos1 peer chaincode install -n enccc -v $(VERSION) -p enccc
	docker exec cli.surgery.pubhos1 peer chaincode install -n enccc -v $(VERSION) -p enccc
	docker exec cli.ir.pubhos1 peer chaincode install -n enccc -v $(VERSION) -p enccc
	docker exec cli.icu.pubhos2 peer chaincode install -n enccc -v $(VERSION) -p enccc
	docker exec cli.surgery.pubhos2 peer chaincode install -n enccc -v $(VERSION) -p enccc
	docker exec cli.ir.pubhos2 peer chaincode install -n enccc -v $(VERSION) -p enccc
	docker exec cli.icu.prihos1 peer chaincode install -n enccc -v $(VERSION) -p enccc
	docker exec cli.surgery.prihos1 peer chaincode install -n enccc -v $(VERSION) -p enccc
	docker exec cli.ir.prihos1 peer chaincode install -n enccc -v $(VERSION) -p enccc
	# instantiate chaincode
	docker exec cli.icu.pubhos1 peer chaincode instantiate -o orderer.example.com:7050 -n test -v $(VERSION) -C publicchannel -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
	docker exec cli.icu.pubhos1 peer chaincode instantiate -o orderer.example.com:7050 -n enccc -v $(VERSION) -C publicchannel -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
	sleep 5
	./invoke.sh cli.icu.pubhos1 publicchannel test initLedger
	# instantiate chaincode
	docker exec cli.icu.prihos1 peer chaincode instantiate -o orderer.example.com:7050 -n test -v $(VERSION) -C privatechannel -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
	docker exec cli.icu.prihos1 peer chaincode instantiate -o orderer.example.com:7050 -n enccc -v $(VERSION) -C privatechannel -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
	sleep 5
	./invoke.sh cli.icu.prihos1 privatechannel test initLedger
	# instantiate chaincode
	docker exec cli.icu.pubhos2 peer chaincode instantiate -o orderer.example.com:7050 -n test -v $(VERSION) -C commonchannel -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
	docker exec cli.icu.pubhos2 peer chaincode instantiate -o orderer.example.com:7050 -n enccc -v $(VERSION) -C commonchannel -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
	sleep 5
	./invoke.sh cli.icu.pubhos2 commonchannel test initLedger
