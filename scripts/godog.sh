#!/usr/bin/env bash

cd "$(dirname "$0")"
CUR_DIR=$(pwd)

echo "Starting development network ..."
# ./startDev.sh -n 1

echo "Running GoDog tests..."
echo "Tests may take a few minutes to complete..."
cd ../chaincode/tests; go test;
