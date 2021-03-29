#!/usr/bin/env bash

# Stop and remove containers
docker-compose stop && yes | docker-compose rm
docker images -a | grep "dev-peer" | awk '{print $1}' | xargs docker rmi -f
docker-compose down --volumes
docker network create cc-tools-demo-net
docker volume rm $(docker volume ls -q)

yes | docker volume prune
