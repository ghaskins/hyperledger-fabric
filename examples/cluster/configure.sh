#!/bin/bash

set -x
set -e

CHANNEL_NAME=$1
CHANNEL_SPEC=$2

peer channel create \
     -o primary.orderer:7050 \
     -c $CHANNEL_NAME \
     -f $CHANNEL_SPEC \
     --tls true \
     --cafile build/cryptogen/ordererOrganizations/orderer/ca/orderer-cert.pem

peer channel join -b $CHANNEL_NAME.block
