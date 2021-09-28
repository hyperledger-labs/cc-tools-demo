import fabricClient = require('fabric-client');
import { initializeChannel } from '../sdkCommon';

const query = (
  client: fabricClient,
  txName: string,
  args: string[],
  transientRequest?: Buffer,
) => {
  const channel = client.getChannel(process.env.CHANNEL);

  return new Promise<string | Error>(async (resolve, reject) => {
    initializeChannel(client, channel)
      .then(async (res: any) => {
        const txId = client.newTransactionID(true);
        const mspOrg = client.getMspid();
        const peers = channel.getPeersForOrg(mspOrg).map((chanPeer => chanPeer.getPeer()));
        const peerIdx = Math.floor(Math.random() * peers.length)
        const queryReq: fabricClient.ChaincodeQueryRequest = {
          args,
          txId,
          chaincodeId: process.env.CCNAME,
          fcn: txName,
          targets: [ peers[peerIdx] ],
        };

        if (transientRequest) {
          const transientMap: fabricClient.TransientMap = {
            '@request': transientRequest,
          };
          queryReq.transientMap = transientMap;
        }

        channel.queryByChaincode(queryReq, true)
          .then((response) => {
            const responseObj = response[0] as any
            if (responseObj.status && responseObj.status != 200) {
              reject(responseObj);
            }
            resolve(response[0].toString('utf-8'));
          })
          .catch((err) => {
            reject(err);
          });
      })
      .catch((err) => {
        reject(err);
      });
  });
};

export default query;
