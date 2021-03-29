#!/usr/bin/env bash

if [ $# -lt 1 ] ; then
  HOST="localhost"
else
  HOST=$1
fi

printf "Sending requests to ${HOST}"

printf '\n\nGet Header\n';
curl -k \
  "http://${HOST}:80/api/query/getHeader" \
  -H 'cache-control: no-cache'


printf '\n\nGet Transactions\n';
curl -k \
  "http://${HOST}:80/api/query/getTx" \
  -H 'cache-control: no-cache'

printf '\n\nGet CreateAsset definition\n';
curl -k -X POST \
  "http://${HOST}:80/api/query/getTx" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "txName": "createAsset"
      }'


printf '\n\nGet Asset Types\n';
curl -k \
  "http://${HOST}:80/api/query/getSchema/" \
  -H 'cache-control: no-cache'

printf '\n\nGet person schema\n';
curl -k -X POST \
  "http://${HOST}:80/api/query/getSchema" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "assetType": "person"
      }'

printf '\n\nCreate person\n'
curl -k -X POST \
  "http://${HOST}/api/invoke/createAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "asset": [
    {
      "@assetType": "person",
      "name": "Maria",
      "cpf": "318.207.920-48",
      "readerScore": 70
    }
  ]
}'

printf '\n\nCreate book\n'
curl -k -X POST \
  "http://${HOST}:80/api/invoke/createAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "asset": [
    {
      "@assetType": "book",
      "title": "Meu Nome é Maria",
      "author": "Maria Viana",
      "currentTenant": {
        "name": "Maria"
      },
      "genres": ["biography", "non-fiction"],
      "published": "2019-05-06T22:12:41Z"
    }
  ]
}'

printf '\n\nRead book\n';
curl -k -X POST \
  "http://${HOST}:80/api/query/readAsset" \
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
  "http://${HOST}/api/invoke/updateAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "update": {
      "@assetType": "person",
      "name": "Maria",
      "readerScore": 75
    }
}'

printf '\n\nRead person to check if it was updated\n';
curl -k \
  "http://${HOST}:80/api/query/readAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "key": {
          "@assetType": "person",
          "name": "Maria"
        }
      }'

printf '\n\nQuery all books using couchdb queries\n';
curl -k \
  "http://${HOST}/api/query/search?@request=eyJxdWVyeSI6eyJzZWxlY3RvciI6eyJuYW1lIjoiTWFyaWEifX19" \
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
  "http://${HOST}/api/invoke/deleteAsset" \
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
  "http://${HOST}:80/api/invoke/deleteAsset" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "key": {
    "@assetType": "person",
    "name": "Maria"
  }
}'

printf '\n\nRead person History\n';
curl -k \
  "http://${HOST}/api/query/readAssetHistory?@request=eyJrZXkiOnsiQGFzc2V0VHlwZSI6InNhbXBsZVBlcnNvbiIsIm5hbWUiOiJNYXJpYSJ9fQ==" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
        "key": {
          "@assetType": "person",
          "name": "Maria"
        }
      }'

printf '\n'
