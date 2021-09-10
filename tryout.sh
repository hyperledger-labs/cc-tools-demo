#!/usr/bin/env bash

if [ $# -lt 1 ] ; then
  HOST="localhost"
else
  HOST=$1
fi

printf "Sending requests to ${HOST}"

printf '\n\nGet Header\n';
curl -k \
  "http://${HOST}/mainchannel/org1/api/query/getHeader" \
  -H 'cache-control: no-cache'

exit
printf '\n\nGet Transactions\n';
curl -k \
  "https://${HOST}/mainchannel/org1/api/query/getTx" \
  -H 'cache-control: no-cache'

printf '\n\nGet CreateAsset definition\n';
curl -k -X POST \
  "https://${HOST}/mainchannel/org1/api/query/getTx" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "txName": "createAsset"
      }'


printf '\n\nGet Asset Types\n';
curl -k \
  "https://${HOST}/mainchannel/org1/api/query/getSchema/" \
  -H 'cache-control: no-cache'

printf '\n\nGet person schema\n';
curl -k -X POST \
  "https://${HOST}/mainchannel/org1/api/query/getSchema" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "assetType": "person"
      }'

printf '\n\nCreate person\n'
curl -k -X POST \
  "https://${HOST}/mainchannel/org1/api/invoke/createAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "asset": [
    {
      "@assetType": "person",
      "name": "Maria",
      "id": "318.207.920-48"
    }
  ]
}'

printf '\n\nCreate book\n'
curl -k -X POST \
  "https://${HOST}/mainchannel/org2/api/invoke/createAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "asset": [
    {
      "@assetType": "book",
      "title": "Meu Nome é Maria",
      "author": "Maria Viana",
      "currentTenant": {
        "id": "318.207.920-48"
      },
      "genres": ["biography", "non-fiction"],
      "published": "2019-05-06T22:12:41Z"
    }
  ]
}'

printf '\n\nRead book\n';
curl -k -X POST \
  "https://${HOST}/mainchannel/org1/api/query/readAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "key": {
          "@assetType": "book",
          "author": "Maria Viana",
          "title": "Meu Nome é Maria"
        },
        "resolve": true
      }'

printf '\n\nUpdate person\n'
curl -k -X PUT \
  "https://${HOST}/mainchannel/org1/api/invoke/updateAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "update": {
      "@assetType": "person",
      "id": "318.207.920-48",
      "name": "Maria",
      "readerScore": 75
    }
}'

printf '\n\nRead person to check if it was updated\n';
curl -k -X POST \
  "https://${HOST}/mainchannel/org1/api/query/readAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "key": {
          "@assetType": "person",
          "id": "318.207.920-48"
        }
      }'

printf '\n\nQuery all books using couchdb queries\n';
curl -k -X POST \
  "https://${HOST}/mainchannel/org1/api/query/search" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "query": {
          "selector": {
            "@assetType": "book"
          }
        },
        "resolve": true
      }'

printf '\n\nDelete book\n'
curl -k -X DELETE \
  "https://${HOST}/mainchannel/org2/api/invoke/deleteAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "key": {
    "@assetType": "book",
    "title": "Meu Nome é Maria",
    "author": "Maria Viana"
  }
}'

printf '\n\nDelete person\n'
curl -k -X DELETE \
  "https://${HOST}/mainchannel/org1/api/invoke/deleteAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "key": {
    "@assetType": "person",
    "id": "318.207.920-48"
  }
}'

printf '\n\nRead person History\n';
curl -k \
  "https://${HOST}/mainchannel/org1/api/query/readAssetHistory" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "key": {
          "@assetType": "person",
          "id": "318.207.920-48"
        }
      }'

printf '\n'
