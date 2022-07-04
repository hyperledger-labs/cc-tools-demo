#!/usr/bin/env bash

version='0'
ORG_QTY=3

while [ $OPTIND -le "$#" ]
do
    if getopts n: opt
    then
        case $opt
        in
            n) ORG_QTY=${OPTARG}
            ;;
        esac
    else
        version="${!OPTIND}"
        ((OPTIND++))
    fi
done

if [ $version = '0' ]
then
    printf 'Usage: ./upgradeCC.sh <version>\n'
    printf 'Flags:\n'
    printf '    -n   Quantity of organizations (Min: 1, Max: 3, Default: 3)\n'
    exit
fi

if [ $ORG_QTY -gt 3 -o $ORG_QTY -lt 1 ]
then
  echo 'WARNING: The maximum number of organizations allowed is 3 and the minimum is 1.'
  echo 'Defaulting to 3 organizations.'
  ORG_QTY=3
fi

cd ./chaincode; go fmt ./...; cd ..

IDENTITIES="{
    \"role\": {
            \"name\": \"member\",
            \"mspId\": \"orgMSP\"
    }
}"

POLICIES="{
    \"signed-by\": 0
}"

curl -sS -k -X POST \
    'https://localhost:3000/api/v1/network/channel/chaincode/install' \
    -H 'Content-Type: application/json' \
    --header 'gofabricversion: 0.9.0' \
    -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
    -d "{
            \"channelName\": \"mainchannel\",
            \"chaincode\": \"cc-tools-demo\",
            \"chaincodeVersion\": \"${version}\"
        }" > /dev/null

if [ $ORG_QTY -gt 1 ]
then
    curl -sS -k -X POST \
        'https://localhost:3001/api/v1/network/channel/chaincode/install' \
        --header 'gofabricversion: 0.9.0' \
        -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
        -H 'Content-Type: application/json' \
        -d "{
                \"channelName\": \"mainchannel\",
                \"chaincode\": \"cc-tools-demo\",
                \"chaincodeVersion\": \"${version}\"
            }" > /dev/null

    IDENTITIES="{
        \"role\": {
                \"name\": \"member\",
                \"mspId\": \"org1MSP\"
            }
        },
        {
            \"role\": {
                \"name\": \"member\",
                \"mspId\": \"org2MSP\"
            }
        }
    }"

    POLICIES="{
        \"signed-by\": 1
    },
    {
        \"signed-by\": 2
    }"
fi

if [ $ORG_QTY -eq 3 ]
then
    curl -sS -k -X POST \
        'https://localhost:3002/api/v1/network/channel/chaincode/install' \
        --header 'gofabricversion: 0.9.0' \
        -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
        -H 'Content-Type: application/json' \
        -d "{
                \"channelName\": \"mainchannel\",
                \"chaincode\": \"cc-tools-demo\",
                \"chaincodeVersion\": \"${version}\"
            }" > /dev/null

    IDENTITIES="{
        \"role\": {
                \"name\": \"member\",
                \"mspId\": \"org1MSP\"
            }
        },
        {
            \"role\": {
                \"name\": \"member\",
                \"mspId\": \"org2MSP\"
            }
        },
        {
            \"role\": {
                \"name\": \"member\",
                \"mspId\": \"org3MSP\"
            }
    }"
fi

printf "\nUpgrade chaincode to version $1\n"

curl -k -X POST \
    'https://localhost:3000/api/v1/network/channel/chaincode/upgrade' \
    --header 'gofabricversion: 0.9.0' \
    -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
    -H 'Content-Type: application/json' \
    -d "{
            \"channelName\": \"mainchannel\",
            \"chaincode\": \"cc-tools-demo\",
            \"chaincodeVersion\": \"${version}\",
            \"chaincodeType\": \"golang\",
            \"endorsement\": {
              \"identities\": [
                  ${IDENTITIES}
              ],
              \"policy\": {
                  \"1-of\": [
                      {
                          \"signed-by\": 0
                      }
                  ]
              }
          }
        }"

printf '\n'
