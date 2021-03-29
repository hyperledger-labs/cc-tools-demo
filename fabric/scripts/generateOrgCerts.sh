#! /bin/bash

generateOrgCerts.generate() {
    if [[ ! ${ORGANIZATION_DOMAIN} ]]; then
        generateCerts.utils.echoErr "Missing ORGANIZATION_DOMAIN env var, check docker-compose file"
        exit 1
    fi

    local peers=('peer0')
    # Append org domain to each peer
    local peers_full=("${peers[@]/%/.${ORGANIZATION_DOMAIN}}")
    for peer in ${peers_full[@]}; do
        # register it
        local PASSHASH=$(generateCerts.common.registerIdentity ${peer})
        if [[ ! ${PASSHASH} ]]; then
            generateCerts.utils.echoErr "Could not get password from ${peer} registration"
            exit 1
        fi
        # enroll it
        generateCerts.common.enrollIdentity $peer $PASSHASH
        if [[ $? -ne 0 ]]; then
            generateCerts.utils.echoErr "Could not enroll ${peer} certs, check ca.${ORGANIZATION_DOMAIN} logs"
            exit 1
        fi
    done

    local orderers=('orderer0')
    # Append org domain to each orderer
    local orderers_full=("${orderers[@]/%/.${ORGANIZATION_DOMAIN}}")
    for orderer in ${orderers_full[@]}; do
        # register it
        local PASSHASH=$(generateCerts.common.registerIdentity ${orderer})
        if [[ ! ${PASSHASH} ]]; then
            generateCerts.utils.echoErr "Could not get password from ${orderer} registration"
            exit 1
        fi
        # enroll it
        generateCerts.common.enrollIdentity $orderer $PASSHASH
        if [[ $? -ne 0 ]]; then
            generateCerts.utils.echoErr "Could not enroll ${orderer} certs, check ca.${ORGANIZATION_DOMAIN} logs"
            exit 1
        fi
    done

    # IDMAP is an associative array whihc maps created identities to generate MSPTree
    declare -A IDMAP
    IDMAP["users"]="Admin.${ORGANIZATION_DOMAIN}"
    IDMAP["peers"]="${peers_full[@]}"
    IDMAP["orderers"]="${orderers_full[@]}"

    # THIS IS THE ONLY WAY TO PASS AN ASSOCIATIVE ARRAY TO FUNCTION
    ## more info: https://stackoverflow.com/a/8879444
    generateOrgCerts.createMSPTree "$(declare -p IDMAP)"
    return $?

}

### @createMSPTree
## Takes an associative array containaing key "users" and "peerss" which are
## identities names, this function assumes your admin is Admin.${ORG_DOMAIN}
## Given those names create a new msp tree to block creation
generateOrgCerts.createMSPTree() {
    if [[ ! ${FABRIC_CA_CLIENT_HOME} ]]; then
        generateCerts.utils.echoErr "no FABRIC_CA_CLIENT_HOME set, have you called generateCerts.utils.prepareClientEnvironment?"
        exit 1
    fi

    if [[  ! ${FABRIC_CA_HOME} ]]; then
        generateCerts.utils.echoErr "no FABRIC_CA_HOME set, check docker-compose file"
        exit 1
    fi

    cd ${FABRIC_CA_CLIENT_HOME}
    if [[ ! -d "msp" ]]; then
        generateCerts.utils.echoErr "No msp folder, make sure you've enrolled your identities"
        exit 1
    fi

    local ADMIN_DIRECTORY="Admin.${ORGANIZATION_DOMAIN}"
    if [[ ! -d msp/$ADMIN_DIRECTORY ]]; then
        generateCerts.utils.echoErr "Identity $ADMIN_DIRECTORY folder is missing. Make sure you've enrolled org admin"
        exit 1
    fi

    local CRYPTOFOLDER=crypto-config/peerOrganizations/${ORGANIZATION_DOMAIN}
    mkdir -p $CRYPTOFOLDER/
    eval "declare -A MAP="${1#*=}
    for folder in "${!MAP[@]}"; do
        mkdir -p $CRYPTOFOLDER/${folder}
        for id in ${MAP[$folder]}; do
            local original_certs_folder=msp/${id}
            if [[ ! -d ${original_certs_folder} ]]; then
                generateCerts.utils.echoErr "No such folder ${id}, did you enroll this identity?"
            fi
            local id_directory=$CRYPTOFOLDER/${folder}/$id

            # Create msp and tls folder under user dir
            mkdir -p $id_directory/msp $id_directory/tls

            # Copy admin certs to users folder
            mkdir -p $id_directory/msp/signcerts $id_directory/msp/admincerts $id_directory/msp/tlscacerts
            cp ${original_certs_folder}/signcerts/cert.pem $id_directory/msp/signcerts/${id}-cert.pem
            cp msp/$ADMIN_DIRECTORY/signcerts/cert.pem $id_directory/msp/admincerts/${ADMIN_DIRECTORY}-cert.pem

            # copy keystore and ca certs just as it is
            cp -r ${original_certs_folder}/keystore $id_directory/msp/
            cp -r ${original_certs_folder}/cacerts $id_directory/msp/

            cp $FABRIC_CA_HOME/ca-cert.pem $id_directory/msp/tlscacerts
            # These two Issuer files are new on fabric so I'm copying just to make sure
            cp $original_certs_folder/Issuer* $id_directory/msp/
            if [[ $folder == "users" ]]; then
                local tls_type="client"
            else
                local tls_type="server"
            fi
            cp $original_certs_folder/tls/signcerts/cert.pem $id_directory/tls/$tls_type.crt
            cp $original_certs_folder/tls/keystore/* $id_directory/tls/$tls_type.key
            cp $original_certs_folder/tls/tlscacerts/* $id_directory/tls/ca.crt
        done
    done



    if [[ ! ${FABRIC_CA_SERVER_CA_NAME} ]]; then
        generateCerts.utils.echoErr "Missing env variable FABRIC_CA_SERVER_CA_NAME. Please check docker-compose file"
        exit 1
    fi

    local MSP_FOLDER=$CRYPTOFOLDER/msp
    mkdir -p $MSP_FOLDER/admincerts $MSP_FOLDER/cacerts $MSP_FOLDER/tlscacerts

    cp msp/$ADMIN_DIRECTORY/signcerts/cert.pem $MSP_FOLDER/admincerts/$ADMIN_DIRECTORY-cert.pem
    if [[ $? -ne 0 ]]; then
        return $?
    fi

    cp $FABRIC_CA_HOME/ca-cert.pem $MSP_FOLDER/cacerts/$FABRIC_CA_SERVER_CA_NAME-cert.pem
    if [[ $? -ne 0 ]]; then
        return $?
    fi

    # cp $FABRIC_CA_HOME/tls-cert.pem $MSP_FOLDER/tlscacerts/tlsca.${ORGANIZATION_DOMAIN}-cert.pem
    cp $FABRIC_CA_HOME/ca-cert.pem $MSP_FOLDER/tlscacerts/tlsca.${ORGANIZATION_DOMAIN}-cert.pem

    return $?
}

create_node_ous () {
    node_file_path=crypto-config/peerOrganizations/${ORGANIZATION_DOMAIN}/msp
    cat > $node_file_path/config.yaml << EOL
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/ca.${ORGANIZATION_DOMAIN}-cert.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/ca.${ORGANIZATION_DOMAIN}.pem
    OrganizationalUnitIdentifier: peer
EOL

}


##########################################################
################## Script starts here ####################
##########################################################
main () {
    ## Generate peers certs register/enroll certs for peers
    . ${FABRIC_CA_HOME}/scripts/generateCerts.sh

    # generateAdminCerts is imported from generateCerts.sh
    generateCerts.common.generateAdminCerts
    if [[ $? -ne 0 ]]; then
        generateCerts.utils.echoErr "Something went wrong with admin cert generation, check ca logs"
        exit 1
    fi
    printf "\e[92mAdmin certs generated\e[39m\n"

    generateOrgCerts.generate
    if [[ $? -ne 0 ]]; then
        generateCerts.utils.echoErr "Something went wrong with peers cert generation, check ca logs"
        exit 1
    fi
    printf "\e[92mpeers certs generated\e[39m\n"

    create_node_ous
}




## this calls main only if script is executed not sourced
unset BASH_SOURCE 2>/dev/null
test ".$0" != ".$BASH_SOURCE" || main "$@"
