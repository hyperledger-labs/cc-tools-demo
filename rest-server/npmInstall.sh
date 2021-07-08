#!/usr/bin/env bash
echo "Node package install for GoLedger CC-Tools Rest server"

# npm install inside containuer
sudo rm -rf node_modules
docker-compose -f docker-compose-npm-install.yaml up

