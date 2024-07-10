#!/usr/bin/env bash

cd "$(dirname "$0")"

ORG_QNTY=3

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
    cd ../ccapi; docker-compose -f docker-compose-1org.yaml down; docker-compose -f docker-compose-1org.yaml up -d --build; cd ..
else
    cd ../ccapi; docker-compose down; docker-compose up -d --build; cd ..
fi