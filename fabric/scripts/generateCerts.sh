#! /bin/bash

### @generateCerts
## This will generate certificates both for orderer and org1/org2 certs
## orderer.example.com, org1.example.com, org2.example.com are hard-coded
## (this is a development script and should not be used for production purpose)
## all called scripts are under same this same folder called scripts but these
## are called inside proper dockers
generateCerts() {
    docker-compose -f docker-compose-ca.yaml rm -f

    if [ "$1" == 1 ]; then
        local orgs=(org.example.com)
    else
        local orgs=(org1.example.com org2.example.com org3.example.com)
    fi

    for ORG in ${orgs[@]}; do
        docker-compose -f docker-compose-ca.yaml up -d ca.$ORG
    done
    sleep 5
    
    local SCRIPTS_FOLDER="/etc/hyperledger/fabric-ca-server/scripts"

    for ORG in ${orgs[@]}; do
        ca="ca.$ORG"

        docker exec -i $ca bash "${SCRIPTS_FOLDER}/generateOrgCerts.sh"
        retVal=$?
        if [[ $? -ne 0 ]]; then
            return $?
        fi

        docker cp $ca:/etc/hyperledger/fabric-ca-client/crypto-config/peerOrganizations ./crypto-config/
        retVal=$?
        if [[ $? -ne 0 ]]; then
            return $?
        fi
    done

    return 0
}

##############################
## COMMON FUNCTIONS SECTION ##
##############################
## This section is a kind of "util.sh" related to certificate generation
## All these fuctions bellow depends on env variables exported on docker files
## if any of those aren't exported then theses functions will throw error and
## exit with non-zero code

### @generateAdminCerts
## generates Admin certs for a organization, this cert will be inside msp admincerts
## generateAdminCerts saves msp to folder under $FABRIC_CA_CLIENT_HOME/msp/Admin.${ORGANIZATION_DOMAIN}
## it doesn't organize cert files in msp structure. This is work for module that is calling
## this function
generateCerts.common.generateAdminCerts() {
    if [[ ! ${ORGANIZATION_DOMAIN} ]]; then
        generateCerts.utils.echoErr "Missing ORGANIZATION_DOMAIN on system env, check docker-compose file"
        exit 1
    fi
    generateCerts.utils.prepareClientEnvironment

    local IDENTITY=Admin.${ORGANIZATION_DOMAIN}
    local PASSHASH=$(generateCerts.common.registerIdentity $IDENTITY)
    if [[ ! ${PASSHASH} ]]; then
        generateCerts.utils.echoErr "Passhash var is void, something went wrong at id ${IDENTITY} registration"
        exit 1
    fi

    generateCerts.common.enrollIdentity $IDENTITY $PASSHASH
    if [[ $retVal -ne 0 ]]; then
        generateCerts.utils.echoErr "Could not enroll ${IDENTITY} certs, check ca.${ORGANIZATION_DOMAIN} logs"
        exit 1
    fi

    return $?
}

### @registerIdentity
## this func takes only one parameter: identity name and registers it
## returns password in stdout so when called needs to save it through command:
## passwd=$(generateCerts.common.registerIdentity ${id})
generateCerts.common.registerIdentity() {
    if [[ ! ${ORGANIZATION_DOMAIN} || ! ${FABRIC_CA_HOME} ]]; then
        generateCerts.utils.echoErr "Missing ORGANIZATION_DOMAIN or FABRIC_CA_HOME on system env, check docker-compose file"
        exit 1
    elif [[ ! ${FABRIC_CA_CLIENT_HOME} ]]; then
        echo "Missing FABRIC_CA_CLIENT_HOME var, did you call generateCerts.utils.prepareClientEnvironment?"
        exit 1
    fi

    # Register an identity
    local REG_OUTPUT=$(fabric-ca-client register --id.name $1 --id.type client --tls.certfiles ${FABRIC_CA_HOME}/tls-cert.pem)
    retVal=$?
    if [[ retVal -ne 0 ]]; then
        generateCerts.utils.echoErr "Could not register admin, check ca.${ORGANIZATION_DOMAIN} logs"
        exit 1
    fi

    # Last line from register command contains user password, to extract this we ask for
    # output tail -1
    local PASSWD=$(echo $REG_OUTPUT | tail -1)

    # PASSWD="Password: somehash", once we just want somehash, then we need to do following
    # substitution -> "removes everything before first ': '"
    local PASSHASH=${PASSWD#*: }
    echo $PASSHASH
}

### @enrollIdentity
## this function takes two parameters and enrolls a identity
## - $1: identity name
## - $2: passhash
generateCerts.common.enrollIdentity() {
    # Enroll org admin id
    if [[ ! ${FABRIC_CA_CLIENT_HOME} ]]; then
        echo "Missing FABRIC_CA_CLIENT_HOME var, did you call generateCerts.utils.prepareClientEnvironment?"
        exit 1
    fi

    mkdir -p ${FABRIC_CA_CLIENT_HOME}/msp

    fabric-ca-client enroll -u https://$1:$2@ca.${ORGANIZATION_DOMAIN}:7054 -M $FABRIC_CA_CLIENT_HOME/msp/$1 --tls.certfiles $FABRIC_CA_HOME/tls-cert.pem --csr.hosts "$1"
    retVal=$?
    if [[ $retVal -ne 0 ]]; then
        generateCerts.utils.echoErr "Could not enroll $1 certs, check ca.${ORGANIZATION_DOMAIN} logs"
        exit 1
    fi

    fabric-ca-client enroll -u https://$1:$2@ca.${ORGANIZATION_DOMAIN}:7054 -M $FABRIC_CA_CLIENT_HOME/msp/$1/tls --tls.certfiles $FABRIC_CA_HOME/tls-cert.pem --enrollment.profile tls --csr.hosts "$1"
    retVal=$?
    if [[ $retVal -ne 0 ]]; then
        generateCerts.utils.echoErr "Could not enroll $1 certs, check ca.${ORGANIZATION_DOMAIN} logs"
        exit 1
    fi

    # return enroll's response so we can catch up this value upstream
    return $retVal
}


### @prepareEnvironment
## export necessary variables for certificate generation work
## aswell creates client folder and enroll ca's admin
generateCerts.utils.prepareClientEnvironment() {
    if [[ ! ${ORGANIZATION_DOMAIN} || ! ${FABRIC_CA_HOME} ]]; then
        generateCerts.utils.echoErr "Missing ORGANIZATION_DOMAIN on system env, check docker-compose file"
        exit 1
    fi

    export FABRIC_CA_CLIENT_HOME=/etc/hyperledger/fabric-ca-client
    mkdir -p $FABRIC_CA_CLIENT_HOME
    rm -rf ${FABRIC_CA_CLIENT_HOME}/crypto-config

    fabric-ca-client enroll -u https://admin:adminpw@ca.${ORGANIZATION_DOMAIN}:7054 --tls.certfiles $FABRIC_CA_HOME/tls-cert.pem
    if [[ $? -ne 0 ]]; then
        generateCerts.utils.echoErr "Could not enroll ca.${ORGANIZATION_DOMAIN}'s admin check logs"
        exit 1
    fi
}


generateCerts.utils.echoErr() {
    # send received message to stderr
    printf "\e[31m$@\e[39m\n" 1>&2;
    exit 1
}

## this calls main only if script is executed not sourced
unset BASH_SOURCE 2>/dev/null
test ".$0" != ".$BASH_SOURCE" || generateCerts "$@"
