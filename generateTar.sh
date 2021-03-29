#!/usr/bin/env bash

# Make sure go mod is up to date
cd chaincode && go mod vendor && cd ..

# Clean rest-server folder
# cd rest-server && sudo rm -rf node_modules dist && cd ..

# Zip file
tar -czf cc-tools-demo.tar.gz chaincode
