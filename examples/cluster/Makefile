PEERS += $(patsubst %,peer%,$(shell seq 1 4))
NODES += $(PEERS)
NODES += orderer
NODES += cli
NODES += ca

CHANNEL_NAME ?= mychannel

CRYPTOGEN = build/cryptogen
ORDERER_ORG = $(CRYPTOGEN)/ordererOrganizations/orderer.net
PEER_ORG= $(CRYPTOGEN)/peerOrganizations/org1.net

CA_PATH = $(PEER_ORG)/ca
ORDERER_PATH = $(ORDERER_ORG)/orderers
PEER_PATH = $(PEER_ORG)/peers
USERS_PATH = $(PEER_ORG)/users

CHANNEL_TXNS=build/channel.tx build/anchor.tx

mspmap.orderer    := $(ORDERER_PATH)/orderer.orderer.net
mspmap.peer1      := $(PEER_PATH)/peer1.org1.net
mspmap.peer2      := $(PEER_PATH)/peer2.org1.net
mspmap.peer3      := $(PEER_PATH)/peer3.org1.net
mspmap.peer4      := $(PEER_PATH)/peer4.org1.net
mspmap.cli        := $(USERS_PATH)/Admin@org1.net

COMPOSE=docker-compose -f compose/docker-compose.yaml
DRUN=$(COMPOSE) run --rm cli

TLS ?= true
export TLS_ENABLED=$(TLS)

help:
	@cat usage.txt

all: compose

compose: nodes
	$(COMPOSE) up -d ca $(PEERS)
	$(DRUN) ./configure.sh $(CHANNEL_NAME) "$(CHANNEL_TXNS)" "$(PEERS)"

nodes: $(patsubst %,build/nodes/%,$(NODES))

$(CRYPTOGEN): config/cryptogen.yaml
	@mkdir -p ${@D}
	cryptogen generate --config $< --output $@

.PRECIOUS: %.yaml
%.yaml:
	@mkdir -p ${@D}
	cp config/${@F} $@

%/genesis.block: build/configtx.yaml build/core.yaml
	@mkdir -p ${@D}
	FABRIC_CFG_PATH=build configtxgen -profile SampleOrg -outputBlock $@

%.tx: build/configtx.yaml build/core.yaml

%/channel.tx:
	@mkdir -p ${@D}
	FABRIC_CFG_PATH=build configtxgen -profile SampleChannel \
		-channelID ${CHANNEL_NAME} \
		-outputCreateChannelTx $@

%/anchor.tx:
	@mkdir -p ${@D}
	FABRIC_CFG_PATH=build configtxgen -profile SampleChannel \
		-channelID ${CHANNEL_NAME} \
		-outputAnchorPeersUpdate $@ \
		-asOrg Org1MSP

.PRECIOUS: %/msp
%/msp: $(CRYPTOGEN)
	$(eval NODE = ${patsubst build/nodes/%/msp,%,${@}})
	@mkdir -p ${@D}
	cp -R $(mspmap.${NODE})/* ${@D}

build/nodes/orderer: build/nodes/orderer/orderer.yaml
build/nodes/orderer: build/nodes/orderer/genesis.block
build/nodes/cli: $(CHANNEL_TXNS)

build/nodes/ca:
	@mkdir -p $@/tls
	cp $(CA_PATH)/*_sk $@/tls/ca.key
	cp $(CA_PATH)/*.pem $@/tls/ca.crt

build/nodes/%: build/nodes/%/msp build/nodes/%/configtx.yaml build/nodes/%/core.yaml
	@echo "Built $@"

clean:
	$(COMPOSE) down
	rm -rf build
