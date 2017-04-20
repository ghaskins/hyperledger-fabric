#!/bin/bash

set -x
set -e

CHANNEL_NAME=$1
CHANNEL_SPEC=$2

dockerip () {
    docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $1 | head -n1
}

orderer=$(dockerip orderer)

echo "orderer0:" $orderer

peer channel create -o $orderer:7050 -c $CHANNEL_NAME -f $CHANNEL_SPEC
peer channel join -b $CHANNEL_NAME.block
