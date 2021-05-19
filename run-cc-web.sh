#!/usr/bin/env bash

if [ $# -lt 1 ] ; then
  PORT="8080"
else
  PORT=$1
fi

docker run -p 0.0.0.0:${PORT}:80/tcp --name cc-web${PORT} goledger/cc-webclient:latest
