import Client = require('fabric-client');

const getAllPeersSDK = (client: Client) => {
  const orgName: string = client.getMspid();
  const targets: Client.Peer[] = client.getPeersForOrg(orgName);
  return targets;
};

export default getAllPeersSDK;
