import fabricClient = require('fabric-client');

const propose = async (ccInvokeRequest: fabricClient.ChaincodeInvokeRequest, channelClient: fabricClient.Channel) => {
  return new Promise<fabricClient.TransactionRequest>(async (resolve, reject) => {
    channelClient
      .sendTransactionProposal(ccInvokeRequest)
      .then(res => {
        const proposalResponses = res[0] as fabricClient.ProposalResponse[];
        const proposal = res[1] as fabricClient.Proposal;

        const txRequest: fabricClient.TransactionRequest = {
          txId: ccInvokeRequest.txId,
          proposalResponses,
          proposal
        };
        return resolve(txRequest);
      })
      .catch(err => {
        const proposalResponses = err.endorsements;
        if (proposalResponses) {
          // Check if all errors are equal
          const firstMessage = proposalResponses[0].message;

          const isSameErr = proposalResponses.reduce((acc, pr) => {
            const sameMsg = pr.message === firstMessage;
            return acc && sameMsg;
          }, true);

          if (isSameErr) {
            return reject(proposalResponses[0]);
          }
        }

        return reject(err);
      });
  });
};

export default propose;
