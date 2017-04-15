#!/bin/bash

CHANNEL_NAME=$1
if [ -z "$1" ]; then
	echo "Setting channel to default name 'mychannel'"
	CHANNEL_NAME="mychannel"
fi

echo "Channel name - "$CHANNEL_NAME
echo

#Backup the original configtx.yaml
cp ../../sampleconfig/configtx.yaml ../../sampleconfig/configtx.yaml.orig
cp configtx.yaml ../../sampleconfig/configtx.yaml

cd $PWD/../../
echo "Building configtxgen"
make configtxgen

echo "Generating genesis block"
./build/bin/configtxgen -profile TwoOrgs -outputBlock orderer.block
mv orderer.block examples/e2e_cli/crypto/orderer/orderer.block

echo "Generating channel configuration transaction"
./build/bin/configtxgen -profile TwoOrgs -outputCreateChannelTx channel.tx -channelID $CHANNEL_NAME
mv channel.tx examples/e2e_cli/crypto/orderer/channel.tx

#reset configtx.yaml file to its original
cp sampleconfig/configtx.yaml.orig sampleconfig/configtx.yaml
rm sampleconfig/configtx.yaml.orig
