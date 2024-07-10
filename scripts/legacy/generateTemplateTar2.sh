#!/usr/bin/env bash

cd "$(dirname "$0")"
cd ../..

# Make sure go mod is up to date
cd chaincode && GOWORK=off go mod vendor && cd ..

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
# mkdir -p tmp/rest-server
# cp -r ccapi/* tmp/rest-server
# tar \
# --exclude=vendor \
# --exclude=src/assettypes/customAssets.go \
# --exclude=src/collections.json \
# --exclude=src/header/header.go \
# -czf cc-tools-demo.tar.gz src -C tmp rest-server
# rm -rf tmp


rm -rf src