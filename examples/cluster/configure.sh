#!/bin/bash

set -x
set -e

CHANNEL_NAME=$1
CHANNEL_SPEC=$2
PEERS=$3

peer channel create \
     -o primary.orderer:7050 \
     -c $CHANNEL_NAME \
     -f $CHANNEL_SPEC \
     --tls true \
     --cafile build/cryptogen/ordererOrganizations/orderer/ca/orderer-cert.pem

echo $PEERS

for PEER in $PEERS
do
    CORE_PEER_ADDRESS=$PEER:7051 peer channel join -b $CHANNEL_NAME.block
done
