all: e2e

NODES := $(patsubst %,peer%,$(shell seq 1 4))
NODES += orderer

PEER_PATH = build/crypto-config/peerOrganizations/peerOrg1/peers

e2e: $(patsubst %,build/nodes/%,$(NODES))

build:
	mkdir -p $@

build/.cryptogen: build
	cryptogen -baseDir ${@D} -ordererNodes 1 -peerOrgs 1 -peersPerOrg 4
	touch $@

build/nodes/peer%: build/.cryptogen
	$(eval ID = ${patsubst build/nodes/peer%,%,${@}})
	@echo $@
	mkdir -p $@/tls
	cp -R $(PEER_PATH)/peerOrg1Peer${ID} $@/msp
	cp $(PEER_PATH)/peerOrg1Peer${ID}/keystore/* $@/tls/peer.key
	cp $(PEER_PATH)/peerOrg1Peer${ID}/signcerts/* $@/tls/peer.crt

build/nodes/orderer: build/.cryptogen
	@echo $@
	mkdir -p $@

clean:
	rm -rf build
