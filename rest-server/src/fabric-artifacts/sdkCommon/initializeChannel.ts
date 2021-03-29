import Client = require('fabric-client');
import { getAllPeersSDK } from '.';
import { Channel } from 'fabric-client';

/**
 * Initialize channel with discovery, if it fails, try with another peer
 * Args:
 *  peerSrc - Where the initialize get the peers info: networkCC, configSDK
 */
const initializeChannel = async (client: Client, channel: Channel) => {
  const peers = channel.getPeers() as Client.ChannelPeer[];
  let initErr: Error;

  for (const peer of peers) {
    try {
      await channel.initialize({ discover: true, asLocalhost: false, target: peer });
      return Promise.resolve(channel);
    } catch (err) {
      initErr = err;
    }
  }
  return Promise.reject(initErr);
};

export default initializeChannel;
