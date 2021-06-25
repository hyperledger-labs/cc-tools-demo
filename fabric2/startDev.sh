#!/usr/bin/env bash

./network.sh down
docker network create cc-tools-demo-net
./network.sh up createChannel
./network.sh deployCC -ccn cc-tools-demo -ccp ../chaincode -ccl go