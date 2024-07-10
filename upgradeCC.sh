#!/usr/bin/env bash

if [[ $# -lt 2 || $# -gt 6 ]]; then
    printf 'Usage: ./upgradeCC.sh <version> <sequence> [-n <org-quantity>] [-ccaas <true | false>]\n'
    exit 1
fi

ORG_QNTY=3
DEPLOY_CCAAS=false
SKIP_COLL_GEN=false

version=$1
sequence=$2

while [[ $# -ge 1 ]] ; do
    key="$1"
    case $key in
        -n )
            ORG_QNTY=$2
            shift
            ;;
        -ccaas )
            DEPLOY_CCAAS=$2
            shift
            ;;
        -c )
            SKIP_COLL_GEN=true
            ;;
  esac
  shift
done

if [ $ORG_QNTY != 3 -a $ORG_QNTY != 1 ]
then
  echo 'WARNING: The number of organizations allowed is either 3 or 1.'
  echo 'Defaulting to 3 organizations.'
  ORG_QNTY=3
fi

if [ "$SKIP_COLL_GEN" = false ] ; then
  echo 'Generating collections configuration file...'
  if [ $ORG_QNTY == 1 ]
  then
    cd ./chaincode; go run . -g --orgs orgMSP; cd ..
  else
    cd ./chaincode; go run . -g --orgs org1MSP org2MSP org3MSP; cd ..
  fi
fi

CCCG_PATH="../chaincode/collections.json"

cd ./chaincode; go fmt ./...; cd ..

cd fabric
if [ "$DEPLOY_CCAAS" = "false" ]; then
  ./network.sh deployCC -ccn cc-tools-demo -ccp ../chaincode -ccl go -ccv $version -ccs $sequence -n $ORG_QNTY -cccg $CCCG_PATH -cci init
else
  ./network.sh deployCCAAS -ccn cc-tools-demo -ccp ../chaincode -ccv $version -ccs $sequence -n $ORG_QNTY -cccg $CCCG_PATH
fi
cd ..
