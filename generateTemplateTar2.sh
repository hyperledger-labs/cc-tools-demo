#!/usr/bin/env bash

# Make sure go mod is up to date
cd chaincode && go mod vendor && cd ..

mkdir src
cp -r chaincode/* src

# Compress file without rest-server
# GoFabric will use the standard CC API from Docker Hub
tar \
--exclude=src/assettypes/customAssets.go \
--exclude=src/collections.json \
--exclude=src/header/header.go \
-czf cc-tools-demo.tar.gz src

# Compress file with rest-server
# GoFabric will use the one provided
# tar \
# --exclude=node_modules \
# --exclude=src/assettypes/customAssets.go \
# --exclude=src/collections.json \
# --exclude=src/header/header.go \
# -czf cc-tools-demo.tar.gz src rest-server


rm -rf src