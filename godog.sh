#!/usr/bin/env bash

which godog
if [ "$?" -ne 0 ]; then
  echo "ERROR: godog tool not found. Please install it to run godog tests."
  echo "  Read more about the tool on: https://github.com/cucumber/godog"
  echo "exiting..."
  exit 1
fi

echo "Running GoDog tests..."
echo "You may be prompted to insert system password for cert generation."
cd ./chaincode/tests; rm -rf ca channel-artifacts crypto-config; godog run;