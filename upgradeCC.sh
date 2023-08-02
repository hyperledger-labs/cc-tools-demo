#!/usr/bin/env bash

if [[ $# -lt 2 || $# -gt 4 ]]; then
    printf 'Usage: ./upgradeCC.sh <version> <sequence> [-n <org-quantity>]\n'
    exit 1
fi

ORG_QNTY=3
SKIP_COLL_GEN=false

while getopts n:c opt; do
    case $opt in
        n)  ORG_QNTY=${OPTARG}
            ;;
        c)  SKIP_COLL_GEN=true
            ;;
    esac
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
        cd ./chaincode; go run . -g --orgs-1; cd ..
    else
        cd ./chaincode; go run . -g --orgs-3; cd ..
    fi
fi

CCCG_PATH="../chaincode/collections.json"

# if [ $ORG_QNTY == 1 ]
# then
#   CCCG_PATH="../chaincode/collections2-org.json"
# else
#   CCCG_PATH="../chaincode/collections.json"
# fi

cd ./chaincode; go fmt ./...; cd ..

version=$1
sequence=$2

cd fabric
./network.sh deployCC -ccn cc-tools-demo -ccp ../chaincode -ccl go -ccv $version -ccs $sequence -n $ORG_QNTY -cccg $CCCG_PATH -cci init
cd ..
