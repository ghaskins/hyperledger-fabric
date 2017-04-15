all: cluster

NODES += $(patsubst %,peer%,$(shell seq 1 4))
NODES += orderer

CHANNEL_NAME ?= mychannel

BASECRYPTO_PATH = build/crypto-config
ORDERER_PATH = $(BASECRYPTO_PATH)/ordererOrganizations/ordererOrg1/orderers/ordererOrg1Orderer1
PEER_PATH = $(BASECRYPTO_PATH)/peerOrganizations/peerOrg1/peers

PROFILE=SingleOrg

cluster: $(patsubst %,build/nodes/%,$(NODES))
	@docker-compose up

build/.cryptogen:
	@mkdir -p ${@D}
	cryptogen -baseDir ${@D} -ordererNodes 1 -peerOrgs 1 -peersPerOrg 4

.PRECIOUS: %.yaml
%.yaml:
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

%/orderer: %/orderer/orderer.yaml
%/orderer: %/orderer/genesis.block
%/orderer: %/orderer/channel.tx
%/orderer: %/orderer/tls/ca.crt
%/orderer: CRYPTO_PATH=$(ORDERER_PATH)

%/peer1: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer1
%/peer2: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer2
%/peer3: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer3
%/peer4: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer4

build/nodes/%: build/nodes/%/msp build/nodes/%/tls build/nodes/%/configtx.yaml build/nodes/%/core.yaml
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
	docker-compose down
	rm -rf build
