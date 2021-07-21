#!/usr/bin/env bash

# Get env file
source .env

while getopts "g" opt; do
    case $opt in
        g)
            GENERATE_CERT=true
            ;;
    esac
done

# Stop and remove containers
docker ps -a | grep "ccapi" | awk '{print $1}' | xargs docker stop
docker ps -a | grep "ccapi" | awk '{print $1}' | xargs docker rm

# Create docker-compose file
if [ $HTTPS = "true" ]
then
  sed "s/PORT/443/g" template/docker-compose-org1-template2.yaml > docker-compose-org1.yaml
  sed "s/PORT/443/g" template/docker-compose-org2-template2.yaml > docker-compose-org2.yaml
  sed "s/PORT/443/g" template/docker-compose-org3-template2.yaml > docker-compose-org3.yaml
  # Generate self-signed cert
else
  sed "s/PORT/80/g" template/docker-compose-org1-template2.yaml > docker-compose-org1.yaml
  sed "s/PORT/80/g" template/docker-compose-org2-template2.yaml > docker-compose-org2.yaml
  sed "s/PORT/80/g" template/docker-compose-org3-template2.yaml > docker-compose-org3.yaml
fi

docker-compose -f docker-compose-org1.yaml -f docker-compose-org2.yaml -f docker-compose-org3.yaml down --volumes

cd scripts
if [ "$HTTPS" == true ]; then
  if [ "$GENERATE_CERT" == true ]; then
    ./generate-dummy-cert.sh -d $DOMAIN
    if [ "$LETS_ENCRYPT" == true ];then
      ./letsencrypt-init.sh -g
    fi
  else
    if [ "$LETS_ENCRYPT" == true ];then
        ./letsencrypt-init.sh
    fi
  fi
fi
cd ..

# Start API
if [ "$GENERATE_CERT" != true ]; then
  docker-compose -f docker-compose-org1.yaml -p ccapi.org1.example.com up -d
  docker-compose -f docker-compose-org2.yaml -p ccapi.org2.example.com up -d
  docker-compose -f docker-compose-org3.yaml -p ccapi.org3.example.com up -d
fi
