#!/usr/bin/env bash

# Clear unused images and volumes
docker rmi $(docker images --quiet --filter "dangling=true")
docker volume rm $(docker volume ls -qf dangling=true)

# Script used to start the development environment.
cd ./chaincode; go fmt ./...; cd ..
cd ./fabric2; ./startDev.sh; cd ..
cd ./rest-server; ./startDev2.sh; cd ..
