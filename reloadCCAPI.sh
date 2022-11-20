#!/usr/bin/env bash

# Script used to start the development environment.
# cd ./rest-server && ./startDev.sh
cd ./ccapi && docker-compose down && docker-compose up -d
