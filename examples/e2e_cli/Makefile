all: e2e

NODES += $(patsubst %,peer%,$(shell seq 1 4))
NODES += orderer

CHANNEL_NAME ?= mychannel

BASECRYPTO_PATH = build/crypto-config
ORDERER_PATH = $(BASECRYPTO_PATH)/ordererOrganizations/ordererOrg1/orderers/ordererOrg1Orderer1
PEER_PATH = $(BASECRYPTO_PATH)/peerOrganizations/peerOrg1/peers

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

%/genesis.block: %/configtx.yaml %/msp
	@mkdir -p ${@D}
	FABRIC_CFG_PATH=${@D} configtxgen -profile $(PROFILE) -outputBlock $@

%/channel.tx: %/configtx.yaml %/msp
	@mkdir -p ${@D}
	FABRIC_CFG_PATH=${@D} configtxgen -profile $(PROFILE) \
		-outputCreateChannelTx $@ \
		-channelID ${CHANNEL_NAME}

build/nodes/orderer: build/nodes/orderer/orderer.yaml
build/nodes/orderer: build/nodes/orderer/genesis.block
build/nodes/orderer: build/nodes/orderer/channel.tx
build/nodes/orderer: build/nodes/orderer/tls/ca.crt
build/nodes/orderer: CRYPTO_PATH=$(ORDERER_PATH)

build/nodes/peer1: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer1
build/nodes/peer2: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer2
build/nodes/peer3: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer3
build/nodes/peer4: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer4

build/nodes/peer%: build/nodes/peer%/core.yaml

build/nodes/%: build/nodes/%/msp build/nodes/%/tls build/nodes/%/configtx.yaml
	@echo "Built $@"

.PRECIOUS: %/msp
%/msp: build/.cryptogen
	@mkdir -p $@
	cp -R $(CRYPTO_PATH)/* $@

.PRECIOUS: %/tls
%/tls: build/.cryptogen
	@mkdir -p $@
	cp $(CRYPTO_PATH)/keystore/* $@/server.key
	cp $(CRYPTO_PATH)/signcerts/* $@/server.crt

.PRECIOUS: %/tls/ca.crt
%/tls/ca.crt: build/.cryptogen
	@mkdir -p ${@D}
	cp $(CRYPTO_PATH)/cacerts/* $@

clean:
	rm -rf build
