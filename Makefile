CAFILE=/etc/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
VERSION := $(shell date +%s)
ENCCC=github.com/hyperledger/fabric/examples/chaincode/go/enccc_example
INSTANTIATE=peer chaincode instantiate -o orderer.example.com:7050 -v $(VERSION) -c '{"Args":["init"]}' --tls --cafile $(CAFILE)
SILENT=

all: set install sleep5 instantiate sleep10 keys sleep6 updatekeys joint_invoke disjoint_invoke #joint disjoint enccc sleep joint_invoke disjoint_invoke

set:
	set -ev

install: joint_install disjoint_install enccc_install

instantiate: joint_instantiate disjoint_instantiate #enccc_instantiate

joint: joint_install joint_instantiate

disjoint: disjoint_install disjoint_instantiate

enccc: enccc_install enccc_instantiate

sleep5:
	sleep 5

sleep10:
	sleep 10 

sleep6:
	sleep 6

joint_install:
	# install chaincode
	docker exec -d cli.icu.pubhos1 peer chaincode install -n joint -v $(VERSION) -p joint
	#docker exec -d cli.surgery.pubhos1 peer chaincode install -n joint -v $(VERSION) -p joint
	#docker exec -d cli.ir.pubhos1 peer chaincode install -n joint -v $(VERSION) -p joint
	docker exec -d cli.icu.pubhos2 peer chaincode install -n joint -v $(VERSION) -p joint
	#docker exec -d cli.surgery.pubhos2 peer chaincode install -n joint -v $(VERSION) -p joint
	#docker exec -d cli.ir.pubhos2 peer chaincode install -n joint -v $(VERSION) -p joint
	docker exec -d cli.icu.prihos1 peer chaincode install -n joint -v $(VERSION) -p joint
	#docker exec -d cli.surgery.prihos1 peer chaincode install -n joint -v $(VERSION) -p joint
	#docker exec -d cli.ir.prihos1 peer chaincode install -n joint -v $(VERSION) -p joint

joint_instantiate:
	# instantiate chaincode
	docker exec $(SILENT) cli.icu.pubhos1 $(INSTANTIATE) -n joint -C commonchannel
	#docker exec $(SILENT) cli.icu.prihos1 $(INSTANTIATE) -n joint -C privatechannel
	#docker exec $(SILENT) cli.icu.pubhos2 $(INSTANTIATE) -n joint -C commonchannel

joint_invoke:
	./invoke.sh cli.icu.pubhos1 commonchannel joint initLedger
	#./invoke.sh cli.icu.prihos1 privatechannel joint initLedger
	#./invoke.sh cli.icu.pubhos2 commonchannel joint initLedger

disjoint_install:
	# install chaincode
	docker exec -d cli.icu.pubhos1 peer chaincode install -n disjoint -v $(VERSION) -p disjoint
	#docker exec -d cli.surgery.pubhos1 peer chaincode install -n disjoint -v $(VERSION) -p disjoint
	#docker exec -d cli.ir.pubhos1 peer chaincode install -n disjoint -v $(VERSION) -p disjoint
	docker exec -d cli.icu.pubhos2 peer chaincode install -n disjoint -v $(VERSION) -p disjoint
	#docker exec -d cli.surgery.pubhos2 peer chaincode install -n disjoint -v $(VERSION) -p disjoint
	#docker exec -d cli.ir.pubhos2 peer chaincode install -n disjoint -v $(VERSION) -p disjoint
	docker exec -d cli.icu.prihos1 peer chaincode install -n disjoint -v $(VERSION) -p disjoint
	#docker exec -d cli.surgery.prihos1 peer chaincode install -n disjoint -v $(VERSION) -p disjoint
	#docker exec -d cli.ir.prihos1 peer chaincode install -n disjoint -v $(VERSION) -p disjoint

disjoint_instantiate:
	# instantiate chaincode
	docker exec $(SILENT) cli.icu.pubhos1 $(INSTANTIATE) -n disjoint -C publicchannel
	docker exec $(SILENT) cli.icu.prihos1 $(INSTANTIATE) -n disjoint -C privatechannel

disjoint_invoke:
	./invoke.sh cli.icu.pubhos1 publicchannel disjoint initLedger
	./invoke.sh cli.icu.prihos1 privatechannel disjoint initLedger

enccc_install:
	# install chaincode
	docker exec -d cli.icu.pubhos1 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)
	#docker exec -d cli.surgery.pubhos1 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)
	#docker exec -d cli.ir.pubhos1 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)
	docker exec -d cli.icu.pubhos2 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)
	#docker exec -d cli.surgery.pubhos2 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)
	#docker exec -d cli.ir.pubhos2 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)
	docker exec -d cli.icu.prihos1 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)
	#docker exec -d cli.surgery.prihos1 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)
	#docker exec -d cli.ir.prihos1 peer chaincode install -n enccc -v $(VERSION) -p $(ENCCC)

enccc_instantiate:
	# instantiate chaincode
	docker exec $(SILENT) cli.icu.pubhos1 $(INSTANTIATE) -n enccc -C commoncchannel
	#docker exec $(SILENT) cli.icu.prihos1 $(INSTANTIATE) -n enccc -C privatechannel
	#docker exec $(SILENT) cli.icu.pubhos2 $(INSTANTIATE) -n enccc -C commonchannel

keys:
	./invoke.sh cli.icu.pubhos1 publicchannel disjoint genPriKey
	./invoke.sh cli.icu.prihos1 privatechannel disjoint genPriKey

updatekeys:
	./invoke.sh cli.icu.pubhos1 commonchannel joint updatePubKey cli.icu.pubhos1 publicchannel
	#./invoke.sh cli.surgery.pubhos1 commonchannel joint updatePubKey cli.surgery.pubhos1 publicchannel
	#./invoke.sh cli.ir.pubhos1 commonchannel joint updatePubKey cli.ir.pubhos1 publicchannel
	./invoke.sh cli.icu.pubhos2 commonchannel joint updatePubKey cli.icu.pubhos2 publicchannel
	#./invoke.sh cli.surgery.pubhos2 commonchannel joint updatePubKey cli.surgery.pubhos2 publicchannel
	#./invoke.sh cli.ir.pubhos2 commonchannel joint updatePubKey cli.ir.pubhos2 publicchannel
	./invoke.sh cli.icu.prihos1 commonchannel joint updatePubKey cli.icu.prihos1 privatechannel
	#./invoke.sh cli.surgery.prihos1 commonchannel joint updatePubKey cli.surgery.prihos1 privatechannel
	#./invoke.sh cli.ir.prihos1 commonchannel joint updatePubKey cli.ir.prihos1 privatechannel
