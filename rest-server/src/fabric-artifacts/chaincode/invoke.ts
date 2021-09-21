import fabricClient = require('fabric-client');
import { initializeChannel } from '../sdkCommon';
import propose from './propose';
import submit from './submit';

const invoke = (
  client: fabricClient,
  txName: string,
  args: string[],
  transientRequest?: Buffer,
  collectionNames?: string[]
) => {
  const channel = client.getChannel(process.env.CHANNEL);

  return new Promise<string | Error>(async (resolve, reject) => {

    initializeChannel(client, channel)
      .then(async (res: any) => {
        const txId = client.newTransactionID(true);

        const ccInvokeRequest: fabricClient.ChaincodeInvokeRequest = {
          args,
          txId,
          chaincodeId: process.env.CCNAME,
          fcn: txName,
        };

        if (transientRequest) {
          const transientMap: fabricClient.TransientMap = {
            '@request': transientRequest,
          };
          ccInvokeRequest.transientMap = transientMap;
        }

        if (collectionNames) {
          ccInvokeRequest.endorsement_hint = {
            chaincodes: [
              {
                name: process.env.CCNAME,
                collection_names: collectionNames,
              },
            ],
          };
        }

        propose(ccInvokeRequest, channel)
          .then((txRequest) => {
            submit(txRequest, client.getMspid(), channel)
              .then((response) => {
                resolve(response);
              })
              .catch((err) => {
                reject(err);
              });
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

export default invoke;
