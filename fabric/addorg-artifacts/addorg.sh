#!/usr/bin/env bash

export IMAGE_TAG=1.4

ORDERER_HOST=orderer0.org1.example.com
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/orderers/$ORDERER_HOST/msp/tlscacerts/ca-cert.pem
CHANNEL_NAME=sys-channel

CMD="peer channel update -f ./addorg-artifacts/sys-update-envelope.pb -c $CHANNEL_NAME -o orderer0.org1.example.com:7050 --tls --cafile $ORDERER_CA"
docker exec cli $CMD

sleep 10
rm -rf bootstrap-block.pb
# Fetch sys-channel config block
CMD="peer channel fetch config ./addorg-artifacts/bootstrap-block.pb -o $ORDERER_HOST:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA"
docker exec cli $CMD
configtxlator proto_decode --input bootstrap-block.pb  --type common.Block > bootstrap-block.json

docker-compose -f ../docker-compose-org2.yaml down
docker-compose -f ../docker-compose-org2.yaml up -d

# Wait for sys RAFT leader to be elected
echo "Waiting for sys RAFT to stabilize"
sleep 15
CHANNEL_NAME=mainchannel
CMD="peer channel update -f ./addorg-artifacts/main-update-envelope.pb -c $CHANNEL_NAME -o orderer0.org1.example.com:7050 --tls --cafile $ORDERER_CA"
docker exec cli $CMD
