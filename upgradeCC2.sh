#!/usr/bin/env bash

if [ $# -ne 2 ] ; then
    printf 'Usage: ./upgradeCC.sh <version> <sequence>\n'
    exit
fi

cd ./chaincode; go fmt ./...; cd ..

version=$1
sequence=$2

cd fabric2
./network.sh deployCC -ccn cc-tools-demo -ccp ../chaincode -ccl go -ccv $version -ccs $sequence -cci init
cd ..
