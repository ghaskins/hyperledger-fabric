all: e2e

NODES += $(patsubst %,peerOrg1Peer%,$(shell seq 1 4))
NODES += orderer

CRYPTO_PATH = build/crypto-config
ORDERER_PATH = $(CRYPTO_PATH)/ordererOrganizations/ordererOrg1/orderers/ordererOrg1Orderer1
PEER_PATH = $(CRYPTO_PATH)/peerOrganizations/peerOrg1/peers

PROFILE=SingleOrg

e2e: $(patsubst %,build/nodes/%,$(NODES))

build:
	mkdir -p $@

build/.cryptogen: build
	cryptogen -baseDir ${@D} -ordererNodes 1 -peerOrgs 1 -peersPerOrg 4
	@touch $@

.PRECIOUS: build/nodes/%.yaml
build/nodes/%.yaml:
	@mkdir -p ${@D}
	cp ${@F} $@

%/genesis.block: %/configtx.yaml
	mkdir -p ${@D}
	FABRIC_CFG_PATH=${@D} configtxgen -profile $(PROFILE) -outputBlock ${@F}

build/nodes/orderer: build/nodes/orderer/configtx.yaml
build/nodes/orderer: build/nodes/orderer/orderer.yaml
build/nodes/orderer: build/nodes/orderer/genesis.block
build/nodes/orderer: build/.cryptogen
	@echo $@
	@mkdir -p $@/tls
	cp -R $(ORDERER_PATH) $@/msp
	cp $(ORDERER_PATH)/keystore/* $@/tls/peer.key
	cp $(ORDERER_PATH)/signcerts/* $@/tls/peer.crt

build/nodes/peer%: build/.cryptogen build/nodes/peer%/configtx.yaml build/nodes/peer%/core.yaml
	@echo $@
	@mkdir -p $@/tls
	cp -R $(PEER_PATH)/${@F} $@/msp
	cp $(PEER_PATH)/${@F}/keystore/* $@/tls/peer.key
	cp $(PEER_PATH)/${@F}/signcerts/* $@/tls/peer.crt

clean:
	rm -rf build
