import Client = require('fabric-client');

const submit = async (txRequest: Client.TransactionRequest, orgMSP: string, channel: Client.Channel) => {
  return new Promise<string | Error>(async (resolve, reject) => {
    // set the transaction listener and set a timeout of 30sec
    // if the transaction did not get committed within the timeout period,
    // fail the test
    const eventPromises: Array<Promise<any>> = [];
    const transactionID = txRequest.txId.getTransactionID();
    let eventhubs: Client.ChannelEventHub[];

    // Try to get the event hubs of active peers using discovery service
    try {
      const discRes = await channel.getDiscoveryResults();
      const orgPeer = discRes['peers_by_org'][orgMSP]['peers'];
      eventhubs = orgPeer.map((peer) => channel.getChannelEventHub(peer.name));
    } catch (error) {
      reject(error);
    }

    eventhubs.forEach((eh) => {
      let txPromise = new Promise((resolve, reject) => {
        let eventTimeout = setTimeout(() => {
          console.log('REQUEST_TIMEOUT: ' + eh.getPeerAddr());
          eh.disconnect();
        }, 60000);

        eh.registerTxEvent(
          transactionID,
          (tx: string, code: string, blockNum: number) => {
            clearTimeout(eventTimeout);

            if (code !== 'VALID') {
              reject(new Error('The invoke chaincode transaction was invalid, code: ' + code));
            } else {
              resolve();
            }
          },
          (err) => {
            clearTimeout(eventTimeout);
            reject(err);
          },
          { unregister: true, disconnect: false }
        );

        eh.connect();
      });
      eventPromises.push(txPromise);
    });

    const sendPromise = channel.sendTransaction(txRequest);
    eventPromises.push(sendPromise);

    try {
      const invokeResults = await Promise.all(eventPromises);
      const response = invokeResults.pop();
      let errorMessage: string;

      if (response.status !== 'SUCCESS') {
        reject(new Error('Failed to order the transaction.'));
      } else {
        resolve(txRequest.proposalResponses[0].response.payload.toString('utf-8'));
      }
    } catch (err) {
      reject(err);
    }
  });
};

export default submit;
