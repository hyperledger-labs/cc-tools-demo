rm ca-cert.pem ca-key.pem
cat fabric-ca-server-template.yaml | sed "s/<<<COMMONNAME>>>/$FABRIC_CA_SERVER_CA_NAME/g" > fabric-ca-server-config.yaml
fabric-ca-server init -b admin:adminpw
chown $USER:$USER fabric-ca-server.db
fabric-ca-server start -b admin:adminpw
