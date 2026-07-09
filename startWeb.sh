#!/usr/bin/env bash

DOCKER_COMPOSE_CMD="docker compose"
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

$DOCKER_COMPOSE_CMD -f ccapi/web-client/docker-compose-goinitus.yaml up -d