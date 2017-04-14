all: e2e

NODES += $(patsubst %,peerOrg1Peer%,$(shell seq 1 4))
NODES += orderer

PEER_PATH = build/crypto-config/peerOrganizations/peerOrg1/peers

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

build/nodes/orderer: build/.cryptogen
	@echo $@
	@mkdir -p $@

build/nodes/peer%: build/.cryptogen build/nodes/peer%/configtx.yaml build/nodes/peer%/core.yaml
	@echo $@
	@mkdir -p $@/tls
	cp -R $(PEER_PATH)/${@F} $@/msp
	cp $(PEER_PATH)/${@F}/keystore/* $@/tls/peer.key
	cp $(PEER_PATH)/${@F}/signcerts/* $@/tls/peer.crt

clean:
	rm -rf build
