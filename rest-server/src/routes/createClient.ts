import * as path from 'path';
import Client = require('fabric-client');
import getDefaultCertsPath from './defaultCertsPath';
import fs = require('fs-extra');

// _client is a private instance of this module
let _client: Client = null;

// _setClient() starts a new fabric-client and returns it
const _setClient = (): Client => {
  const isDocker = process.env.DOCKER;
  let configNetwork = '';

  const configSDKFile = process.env.CONFIG_SDK_FILE;
  const isDeploy = process.env.DEPLOY;

  if (configSDKFile) {
    configNetwork = path.resolve(`/rest-server/${configSDKFile}`);
  } else if (isDeploy) {
    configNetwork = path.resolve('/rest-server/configsdk.yaml');
  } else {
    configNetwork = path.resolve('/rest-server/configsdk.yaml');
  }
  _client = Client.loadFromConfig(configNetwork);
  _client.initCredentialStores();

  const [adminPrivKeyPath, adminPubKeyPath, _] = getDefaultCertsPath();

  const adminPubKey = fs.readFileSync(adminPubKeyPath);
  const adminKey = fs.readFileSync(adminPrivKeyPath);
  _client.setAdminSigningIdentity(
    Buffer.from(adminKey).toString(),
    Buffer.from(adminPubKey).toString(),
    _client.getMspid()
  );

  return _client;
};

// ClientStore is an object containing only one function .get()
// if _client is null, then instantiate and return, otherwise, return
const ClientStore = {
  get: (): Client => (_client ? _client : _setClient()),
  resetClient: (): Client => _setClient()
};

Object.freeze(ClientStore);

export default ClientStore;
