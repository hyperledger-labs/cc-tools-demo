#!/usr/bin/env bash

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

if [ "$SKIP_COLL_GEN" = false ]
then
    ./startDev.sh -n $ORG_QNTY
else
    ./startDev.sh -n $ORG_QNTY -c
fi