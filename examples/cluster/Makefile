all: cluster

NODES += $(patsubst %,peer%.org1,$(shell seq 1 4))
NODES += orderer
NODES += cli

CHANNEL_NAME ?= mychannel

BASECRYPTO_PATH = build/cryptogen
BASEPEER_PATH= $(BASECRYPTO_PATH)/peerOrganizations/org1
ORDERER_PATH = $(BASECRYPTO_PATH)/ordererOrganizations/orderer/orderers/primary.orderer
PEER_PATH = $(BASEPEER_PATH)/peers
CLI_PATH = $(BASEPEER_PATH)/users/Admin@org1

PROFILE=SingleOrg

CHANNEL_SPEC=build/nodes/cli/channel.tx

cluster: nodes
	@docker-compose up -d
	FABRIC_CFG_PATH=build/nodes/cli ./configure.sh $(CHANNEL_NAME) $(CHANNEL_SPEC)

nodes: $(patsubst %,build/nodes/%,$(NODES))

build/.cryptogen: cryptogen.yaml
	mkdir -p ${@D}
	cryptogen generate --config cryptogen.yaml --output $(BASECRYPTO_PATH)
	touch $@

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

build/nodes/peer1.org1: CRYPTO_PATH=$(PEER_PATH)/peer1.org1
build/nodes/peer2.org1: CRYPTO_PATH=$(PEER_PATH)/peer2.org1
build/nodes/peer3.org1: CRYPTO_PATH=$(PEER_PATH)/peer3.org1
build/nodes/peer4.org1: CRYPTO_PATH=$(PEER_PATH)/peer4.org1

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
