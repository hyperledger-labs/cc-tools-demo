#!/usr/bin/env bash

cd "$(dirname "$0")"

ORG_QNTY=3
DOCKER_COMPOSE_CMD="docker compose"

# Check docker compose version
if docker compose &>/dev/null; then
    DOCKER_COMPOSE_CMD="docker compose"
else
    if command -v docker-compose &>/dev/null; then
        DOCKER_COMPOSE_CMD="docker-compose"
    else
        errorln "docker compose or docker-compose command not found. Please install the latest version of docker compose"
        exit 1
    fi
fi

while getopts n: opt; do
    case $opt in
        n)  ORG_QNTY=${OPTARG}
            ;;
    esac
done

if [ $ORG_QNTY != 3 -a $ORG_QNTY != 1 ]
then
  echo 'WARNING: The number of organizations allowed is either 3 or 1.'
  echo 'Defaulting to 3 organizations.'
  ORG_QNTY=3
fi

## This brings up API in Go
if [ $ORG_QNTY == 1 ]
then
    cd ../ccapi; $DOCKER_COMPOSE_CMD -f docker-compose-1org.yaml down; $DOCKER_COMPOSE_CMD -f docker-compose-1org.yaml up -d --build; cd ..
else
    cd ../ccapi; $DOCKER_COMPOSE_CMD down; $DOCKER_COMPOSE_CMD up -d --build; cd ..
fi