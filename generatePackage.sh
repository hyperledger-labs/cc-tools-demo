#!/usr/bin/bash

# Default values for the flags
FLAG_CCAPI="none"
FLAG_LABEL="1.0"
SKIP_COLL_GEN=false

# You can change this if you want to avoid using the --name flag
FLAG_NAME="cc-tools-demo"

# Process command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --label | -l)
            if [[ $# -gt 1 ]]; then
                FLAG_LABEL=$2
                shift 2
            else
                echo "Error: --label flag requires a value."
                exit 1
            fi
            ;;
        --name | -n)
            if [[ $# -gt 1 ]]; then
                FLAG_NAME=$2
                shift 2
            else
                echo "Error: --name flag requires a value."
                exit 1
            fi
            ;;
        --org | -o)
            if [[ $# -gt 1 ]]; then
                orgs+=("$2")
                shift 2
            else
                echo "Error: --org flag requires a value."
                exit 1
            fi
            ;;
        --skip | -s)
            SKIP_COLL_GEN=true
            shift 1
            ;;
        --help | -h)
            echo "Usage: ./generateTar.sh [--label <label>] [--name <name>]"
            echo "  --help , -h: Show this help message."
            echo "  --label, -l: Label to be used for the chaincode package. Default is 1.0."
            echo "  --name , -n: Name of the chaincode package. Default is ${FLAG_NAME}."
            echo "  --org  , -o: Include an organization in the collection configuration file"
            echo "      This option can be used multiple times to include multiple organizations"
            echo "      If no organizations are specified, the default is to include any organization found in the readers"
            echo "  --skip , -s: Skip the generation of the collections.json file"
            exit 0
            ;;
        *)
            # Ignore unrecognized arguments
            shift
            ;;
    esac
done

# Generate collection configuration file
if [ "$SKIP_COLL_GEN" = false ]
then
    if [ ${#orgs[@]} -gt 0 ]
    then
        cd ./chaincode; go run . -g --orgs ${orgs[@]}; cd ..
    else
        cd ./chaincode; go run . -g; cd ..
    fi
fi

# Remove previous tar file
rm -f chaincode.tar.gz

# Make sure go mod is up to date
cd chaincode && GOWORK=off go mod vendor && cd ..

# Pack chaincode
export FABRIC_CFG_PATH=fabric/config
peer lifecycle chaincode package chaincode.tar.gz --path chaincode --lang golang --label ${FLAG_NAME}_${FLAG_LABEL}

