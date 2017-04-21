all: cluster

PEERS += $(patsubst %,peer%.org1,$(shell seq 1 4))
NODES += $(PEERS)
NODES += orderer
NODES += cli

CHANNEL_NAME ?= mychannel

BASECRYPTO_PATH = build/cryptogen
ORDERER_ORG = $(BASECRYPTO_PATH)/ordererOrganizations/orderer
PEER_ORG= $(BASECRYPTO_PATH)/peerOrganizations/org1

ORDERER_PATH = $(ORDERER_ORG)/orderers
PEER_PATH = $(PEER_ORG)/peers
USERS_PATH = $(PEER_ORG)/users

PROFILE=SingleOrg

CHANNEL_SPEC=build/nodes/cli/channel.tx

mspmap.orderer    := $(ORDERER_PATH)/primary.orderer
mspmap.peer1.org1 := $(PEER_PATH)/peer1.org1
mspmap.peer2.org1 := $(PEER_PATH)/peer2.org1
mspmap.peer3.org1 := $(PEER_PATH)/peer3.org1
mspmap.peer4.org1 := $(PEER_PATH)/peer4.org1
mspmap.cli        := $(USERS_PATH)/Admin@org1

cluster: nodes
	docker-compose up -d $(PEERS)
	docker-compose run --rm cli ./configure.sh $(CHANNEL_NAME) $(CHANNEL_SPEC)

nodes: $(patsubst %,build/nodes/%,$(NODES))

$(BASECRYPTO_PATH): cryptogen.yaml
	@mkdir -p ${@D}
	cryptogen generate --config $< --output $@

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
		-channelID ${CHANNEL_NAME} \
		-outputCreateChannelTx $@

.PRECIOUS: %/msp
%/msp: $(BASECRYPTO_PATH)
	$(eval NODE = ${patsubst build/nodes/%/msp,%,${@}})
	@mkdir -p $@
	cp -R $(mspmap.${NODE})/* $@

.PRECIOUS: %/tls
%/tls: %/msp
	$(eval MSP = ${@D}/msp)
	@mkdir -p $@
	cp $(MSP)/keystore/* $@/server.key
	cp $(MSP)/signcerts/* $@/server.crt
	cp $(MSP)/cacerts/* $@/ca.crt

build/nodes/orderer: build/nodes/orderer/orderer.yaml
build/nodes/orderer: build/nodes/orderer/genesis.block
build/nodes/cli: $(CHANNEL_SPEC)

build/nodes/%: build/nodes/%/msp build/nodes/%/tls build/nodes/%/configtx.yaml build/nodes/%/core.yaml
	@echo "Built $@"

clean:
	docker-compose down
	rm -rf build
