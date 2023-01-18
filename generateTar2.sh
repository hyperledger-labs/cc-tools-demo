#!/usr/bin/env bash

if [ $# -ne 1 ] ; then
    printf 'Usage: ./generateTar2.sh <label>\n'
    exit
fi

LABEL=$1

rm -f cc-tools-demo.tar.gz

# Make sure go mod is up to date
cd chaincode && go mod vendor && cd ..

export FABRIC_CFG_PATH=fabric2/config
peer lifecycle chaincode package chaincode.tar.gz --path chaincode --lang golang --label cc-tools-demo_${LABEL}

# Compress file without rest-server (GoFabric will use the standard CC API)
tar -czf cc-tools-demo.tar.gz chaincode.tar.gz

# Compress file with rest-server (GoFabric will use the one provided)
# tar -c --exclude=node_modules -zf cc-tools-demo.tar.gz chaincode.tar.gz rest-server

rm -f chaincode.tar.gz