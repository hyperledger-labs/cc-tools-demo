#!/usr/bin/env bash

# Docs: https://hyperledger-fabric.readthedocs.io/en/release-1.4/channel_update_tutorial.html

# Create org2.json
configtxgen -configPath . -printOrg org2MSP > ./org2.json

# Create vars to shorten commands
ORDERER_HOST=orderer0.org1.example.com
NEW_ORDERERS="orderer0.org2.example.com:7050"
PEER_HOST=peer0.org1.example.com
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/orderers/$ORDERER_HOST/msp/tlscacerts/ca-cert.pem
CHANNEL_NAME=sys-channel

# Fetch sys-channel config block
CMD="peer channel fetch config ./addorg-artifacts/sys-block.pb -o $ORDERER_HOST:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA"
docker exec cli $CMD

# Translate sys-block.pb to JSON
configtxlator proto_decode --input sys-block.pb --type common.Block > sys-block.json

# Extract config data from sys-block.json
cat sys-block.json | jq .data.data[0].payload.data.config > sys-config.json
current_hash=$(jq -M -r '.header | "\(.data_hash)"' sys-block.json)
current_block_number=$(($(jq -M -r '.header | "\(.number)"' sys-block.json)+1))

# Add org2.json to sys-channel config groups
NEW_ORDERER_TLS=$(base64 ../crypto-config/peerOrganizations/org2.example.com/orderers/orderer0.org2.example.com/tls/server.crt -w 0)
NEW_CONSENTER=$(jq -n --arg new_orderer_tls $NEW_ORDERER_TLS  '{"client_tls_cert": $new_orderer_tls, "host": "orderer0.org2.example.com", "port": 7050, "server_tls_cert": $new_orderer_tls}')

rm sys-config-updated.json
jq --argjson new_consenter "$NEW_CONSENTER" \
    '.channel_group.groups.Orderer.values.ConsensusType.value.metadata.consenters += [$new_consenter]'\
    sys-config.json | \
jq --arg new_orderers $NEW_ORDERERS \
    '.channel_group.values.OrdererAddresses.value.addresses += [$new_orderers]' | \
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"org2MSP":$neworgjson[0]}}}}}' \
    --slurpfile neworgjson ./org2.json | \
jq -s '.[0] * {"channel_group":{"groups":{"Consortiums":{"groups": {"SampleConsortium": {"groups": {"org2MSP":$neworgjson[0]}}}}}}}' \
    --slurpfile neworgjson ./org2.json | \
jq -s '.[0] * {"channel_group":{"groups":{"Orderer":{"groups": {"org2MSP":$neworgjson[0]}}}}}' \
    --slurpfile neworgjson ./org2.json \
> ./sys-config-updated.json

# Encode sys-config.json and sys-config-updated.json
configtxlator proto_encode --input sys-config.json --type common.Config --output sys-config.pb
configtxlator proto_encode --input sys-config-updated.json --type common.Config --output sys-config-updated.pb

# Compute update transaction
configtxlator compute_update --channel_id $CHANNEL_NAME --original sys-config.pb --updated sys-config-updated.pb --output sys-update.pb

# Decode sys-update.json
configtxlator proto_decode --input sys-update.pb --type common.ConfigUpdate | jq . > sys-update.json

# Envelope update in block
echo '{"payload":{"header":{"channel_header":{"channel_id":"sys-channel", "type":2}},"data":{"config_update":'$(cat sys-update.json)'}}}' | jq . > sys-update-envelope.json

# Encode update envelope
configtxlator proto_encode --input sys-update-envelope.json --type common.Envelope --output sys-update-envelope.pb

# Change channel name
CHANNEL_NAME=mainchannel

# Fetch mainchannel config block
CMD="peer channel fetch config ./addorg-artifacts/main-block.pb -o $ORDERER_HOST:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA"
docker exec cli $CMD

# Translate main-block.pb to JSON
configtxlator proto_decode --input main-block.pb --type common.Block > main-block.json

# Extract config data from main-block.json
cat main-block.json | jq .data.data[0].payload.data.config > main-config.json

# Add org2.json to mainchannel config groups
rm main-config-updated.json
jq --argjson new_consenter "$NEW_CONSENTER" \
    '.channel_group.groups.Orderer.values.ConsensusType.value.metadata.consenters += [$new_consenter]'\
    main-config.json | \
jq --arg new_orderers $NEW_ORDERERS \
    '.channel_group.values.OrdererAddresses.value.addresses += [$new_orderers]' | \
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"org2MSP":$neworgjson[0]}}}}}' \
    --slurpfile neworgjson ./org2.json | \
jq -s '.[0] * {"channel_group":{"groups":{"Orderer":{"groups": {"org2MSP":$neworgjson[0]}}}}}' \
    --slurpfile neworgjson ./org2.json \
> ./main-config-updated.json

# Encode main-config.json and main-config-updated.json
configtxlator proto_encode --input main-config.json --type common.Config --output main-config.pb
configtxlator proto_encode --input main-config-updated.json --type common.Config --output main-config-updated.pb

# Compute update transaction
configtxlator compute_update --channel_id $CHANNEL_NAME --original main-config.pb --updated main-config-updated.pb --output main-update.pb

# Decode main-update.json
configtxlator proto_decode --input main-update.pb --type common.ConfigUpdate | jq . > main-update.json

# Envelope update in block
echo '{"payload":{"header":{"channel_header":{"channel_id":"mainchannel", "type":2}},"data":{"config_update":'$(cat main-update.json)'}}}' | jq . > main-update-envelope.json

# Encode update envelope
configtxlator proto_encode --input main-update-envelope.json --type common.Envelope --output main-update-envelope.pb
