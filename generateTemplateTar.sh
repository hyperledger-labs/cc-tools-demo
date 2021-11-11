#!/usr/bin/env bash

# Make sure go mod is up to date
cd chaincode && go mod vendor && cd ..

# Compress file without rest-server
# GoFabric will use the standard CC API from Docker Hub
tar \
--exclude=chaincode/assettypes/customAssets.go \
--exclude=chaincode/collections.json \
--exclude=chaincode/header/header.go \
-czf cc-tools-demo.tar.gz chaincode

# Compress file with rest-server
# GoFabric will use the one provided
# tar \
# --exclude=node_modules \
# --exclude=chaincode/assettypes/customAssets.go \
# --exclude=chaincode/collections.json \
# --exclude=chaincode/header/header.go \
# -czf cc-tools-demo.tar.gz chaincode rest-server
