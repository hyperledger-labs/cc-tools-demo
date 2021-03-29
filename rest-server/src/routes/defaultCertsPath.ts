const getDefaultCertsPath = () => {
  const isDocker = process.env.DOCKER;
  const isDeploy = process.env.DEPLOY;
  let adminPrivKeyPath = "",
    adminPubKeyPath = "",
    rootCACertPath = "";
  if (isDeploy) {
    adminPrivKeyPath = "/rest-server/certs/admin.key";
    adminPubKeyPath = "/rest-server/certs/admin.cert";
    rootCACertPath = "/rest-server/certs/ca.pem";
  } else if (isDocker) {
    adminPrivKeyPath = "/certs/admin.key";
    adminPubKeyPath = "/certs/admin.cert";
    rootCACertPath = "/certs/ca.pem";
  } else {
    adminPrivKeyPath = "../fabric/crypto-config/rest-certs/admin.key";
    adminPubKeyPath = "../fabric/crypto-config/rest-certs/admin.cert";
    rootCACertPath = "../fabric/crypto-config/rest-certs/ca.pem";
  }
  return [adminPrivKeyPath, adminPubKeyPath, rootCACertPath];
};

export default getDefaultCertsPath;
