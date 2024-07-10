#!/usr/bin/bash

cd "$(dirname "$0")"

while getopts "o:h" opt; do
    case $opt in
        o) orgs+=("$OPTARG");;
        #...
        h)
            echo "Usage: ./generateCollection.sh [-o <org name>]"
            echo "  -o: Include an organization in the collection configuration file"
            echo "      This option can be used multiple times to include multiple organizations"
            echo "      If no organizations are specified, the default is to include any organization found in the readers"
            echo ""
            echo "Example: ./generateCollection.sh -o org1MSP -o org2MSP -o org3MSP"
            exit 0
            ;;
    esac
done
shift $((OPTIND -1))

if [ ${#orgs[@]} -gt 0 ]
then
    cd ../chaincode; go run . -g --orgs ${orgs[@]}; cd ..
else
    cd ../chaincode; go run . -g; cd ..
fi