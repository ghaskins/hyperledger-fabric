all: cluster

NODES += $(patsubst %,peer%,$(shell seq 1 4))
NODES += orderer
NODES += cli

CHANNEL_NAME ?= mychannel

BASECRYPTO_PATH = build/crypto-config
BASEPEER_PATH= $(BASECRYPTO_PATH)/peerOrganizations/peerOrg1
ORDERER_PATH = $(BASECRYPTO_PATH)/ordererOrganizations/ordererOrg1/orderers/ordererOrg1Orderer1
PEER_PATH = $(BASEPEER_PATH)/peers
CLI_PATH = $(BASEPEER_PATH)/users/peerOrg1Admin

PROFILE=SingleOrg

CHANNEL_SPEC=build/nodes/cli/channel.tx

cluster: nodes
	@docker-compose up -d
	FABRIC_CFG_PATH=build/nodes/cli ./configure.sh $(CHANNEL_NAME) $(CHANNEL_SPEC)

nodes: $(patsubst %,build/nodes/%,$(NODES))

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

build/nodes/orderer: build/nodes/orderer/orderer.yaml
build/nodes/orderer: build/nodes/orderer/genesis.block
build/nodes/orderer: CRYPTO_PATH=$(ORDERER_PATH)

build/nodes/peer1: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer1
build/nodes/peer2: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer2
build/nodes/peer3: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer3
build/nodes/peer4: CRYPTO_PATH=$(PEER_PATH)/peerOrg1Peer4

build/nodes/cli: $(CHANNEL_SPEC)
build/nodes/cli: CRYPTO_PATH=$(CLI_PATH)

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
	cp $(CRYPTO_PATH)/cacerts/* $@/ca.crt

clean:
	docker-compose down
	rm -rf build
