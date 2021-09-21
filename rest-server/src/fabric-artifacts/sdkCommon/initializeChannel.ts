import Client = require('fabric-client');
import { getAllPeersSDK } from '.';
import { Channel } from 'fabric-client';


type LastPeers = {
  [channel: string]: Client.ChannelPeer;
}


let lastPeers: LastPeers = {};

/**
 * Initialize channel with discovery, if it fails, try with another peer
 * Args:
 *  peerSrc - Where the initialize get the peers info: networkCC, configSDK
 */
const initializeChannel = async (client: Client, channel: Channel) => {
  try {
    const channelName = channel.getName();
    if (lastPeers[channelName]) {
      await channel.initialize({ discover: true, asLocalhost: false, target: lastPeers[channelName] });
      return Promise.resolve(channel);
    } 
  } catch (err) {
    console.warn("Last peer not responding, trying others")
    console.warn(err)
  }

  const peers = channel.getPeers() as Client.ChannelPeer[];

  for (const peer of peers) {
    try {
      await channel.initialize({ discover: true, asLocalhost: false, target: peer });
      lastPeers[channel.getName()] = peer;
      return Promise.resolve(channel);
    } catch (err) {
      return Promise.reject(err);
    }
  }
};

export default initializeChannel;
