#!/usr/bin/env bash

which godog
if [ "$?" -ne 0 ]; then
  echo "ERROR: godog tool not found. Please install it to run godog tests."
  echo "  Read more about the tool on: https://github.com/cucumber/godog"
  echo "exiting..."
  exit 1
fi

cd ./chaincode/tests; godog run;