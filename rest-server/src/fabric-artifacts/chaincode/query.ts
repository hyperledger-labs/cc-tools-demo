import fabricClient = require('fabric-client');
import { initializeChannelPeer } from '../sdkCommon';


const query = (
  client: fabricClient,
  txName: string,
  args: string[],
  transientRequest?: Buffer,
) => {
  const channel = client.getChannel(process.env.CHANNEL);

  return new Promise<string | Error>(async (resolve, reject) => {
    initializeChannelPeer(client, channel)
      .then(async (peer: fabricClient.Peer) => {
        const txId = client.newTransactionID(true);
        const mspOrg = client.getMspid();
        const queryReq: fabricClient.ChaincodeQueryRequest = {
          args,
          txId,
          chaincodeId: process.env.CCNAME,
          fcn: txName,
          // TODO: select peers at random
          targets: [ peer ], //channel.getPeersForOrg(mspOrg).map((chanPeer => chanPeer.getPeer()))[0] ],
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
