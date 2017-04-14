#!/bin/bash

CHANNEL_NAME=$1
if [ -z "$1" ]; then
	echo "Setting channel to default name 'mychannel'"
	CHANNEL_NAME="mychannel"
fi

echo "Channel name - "$CHANNEL_NAME
echo

echo "Generating genesis block"
configtxgen -profile TwoOrgs -outputBlock crypto/orderer/orderer.block

echo "Generating channel configuration transaction"
configtxgen -profile TwoOrgs -outputCreateChannelTx crypto/orderer/channel.tx -channelID $CHANNEL_NAME
