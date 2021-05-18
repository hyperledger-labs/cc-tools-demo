#!/usr/bin/env bash

if [ $# -lt 1 ] ; then
  PORT="8080"
else
  PORT=$1
fi

docker run -p 127.0.0.1:${PORT}:80/tcp --name cc-web${PORT} goledger/cc-webclient:latest
