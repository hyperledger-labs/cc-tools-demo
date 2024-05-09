#!/usr/bin/env bash

echo "Starting development network ..."
# ./startDev.sh -n 1

echo "Running GoDog tests..."
echo "Tests may take a few minutes to complete..."
cd ./chaincode/tests; go test;
