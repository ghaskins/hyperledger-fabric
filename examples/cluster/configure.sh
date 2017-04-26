#!/bin/bash

set -x
set -e

CHANNEL_NAME=$1
CHANNEL_TXNS=$2
PEERS=$3

CA_CRT=build/cryptogen/ordererOrganizations/orderer.net/ca/ca.orderer.net-cert.pem

for TXN in $CHANNEL_TXNS
do
    peer channel create -o orderer:7050 \
         -c $CHANNEL_NAME \
         -f $TXN \
         --tls $CORE_PEER_TLS_ENABLED \
         --cafile $CA_CRT
done

for PEER in $PEERS
do
    CORE_PEER_ADDRESS=$PEER:7051 peer channel join -b $CHANNEL_NAME.block
done
