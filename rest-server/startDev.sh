#!/usr/bin/env bash

# Stop and remove containers
docker ps -a | grep "ccapi" | awk '{print $1}' | xargs docker stop
docker ps -a | grep "ccapi" | awk '{print $1}' | xargs docker rm

docker-compose -f docker-compose-hlf1-org1.yaml -f docker-compose-hlf1-org2.yaml -f docker-compose-hlf1-org3.yaml -f docker-compose-temp.yaml down --volumes

# Start API
docker-compose -f docker-compose-temp.yaml -p intermediate-container up >> /dev/null
docker-compose -f docker-compose-hlf1-org1.yaml -p ccapi.org1.example.com up -d
docker-compose -f docker-compose-hlf1-org2.yaml -p ccapi.org2.example.com up -d
docker-compose -f docker-compose-hlf1-org3.yaml -p ccapi.org3.example.com up -d
