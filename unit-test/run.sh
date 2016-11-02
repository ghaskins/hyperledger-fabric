#!/bin/bash

set -e

TAG=$1
export GO_LDFLAGS=$2

BASEIMAGE="hyperledger/fabric-peer"
IMAGE=$BASEIMAGE

if [ "$TAG" != "" ]
then
    IMAGE="$BASEIMAGE:$TAG"
fi

echo "Running unit tests using $IMAGE"

echo "Cleaning membership services folder"
rm -rf membersrvc/ca/.ca/

echo -n "Obtaining list of tests to run.."
# Some examples don't play nice with `go test`
export PKGS=`go list github.com/hyperledger/fabric/... 2> /dev/null | \
                                                  grep -v /vendor/ | \
                                                  grep -v /build/ | \
	                                          grep -v /examples/chaincode/chaintool/ | \
						  grep -v /examples/chaincode/go/asset_management | \
						  grep -v /examples/chaincode/go/utxo | \
						  grep -v /examples/chaincode/go/rbac_tcerts_no_attrs`
echo "DONE!"

echo "Running tests..."
docker-compose up --abort-on-container-exit
