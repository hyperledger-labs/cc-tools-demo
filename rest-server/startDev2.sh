#!/usr/bin/env bash

# Stop and remove containers
docker ps -a | grep "mainchannel" | awk '{print $1}' | xargs docker stop
docker ps -a | grep "mainchannel" | awk '{print $1}' | xargs docker rm

docker-compose -f docker-compose-hlf2-org1.yaml -f docker-compose-hlf2-org2.yaml -f docker-compose-hlf2-org3.yaml -f docker-compose-temp.yaml down --volumes

# Start reverse proxy
cd nginx;
docker-compose up -d
cd ..;

# Start API
docker-compose -f docker-compose-temp.yaml -p intermediate-container up >> /dev/null
docker-compose -f docker-compose-hlf2-org1.yaml -p mainchannel.org1.example.com up -d
docker-compose -f docker-compose-hlf2-org2.yaml -p mainchannel.org2.example.com up -d
docker-compose -f docker-compose-hlf2-org3.yaml -p mainchannel.org3.example.com up -d
