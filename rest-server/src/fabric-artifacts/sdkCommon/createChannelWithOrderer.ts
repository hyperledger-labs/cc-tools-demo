import Client = require('fabric-client');
import { Channel } from 'fabric-client';
import initializeChannel from './initializeChannel';

const createChannelWithOrderer = async (client: Client, channelName: string, discovery: boolean = false) => {
  return new Promise<Channel>((resolve, reject) => {
    const channel = new Channel(channelName, client);
    channel.addOrderer(client.getOrderer('orderer0'));

    if (discovery) {
      initializeChannel(client, channel)
        .then(channel => {
          resolve(channel);
        })
        .catch(err => {
          console.error('Failed to initialize channel with discovery: ', err);
          reject(err);
        });
    }

    resolve(channel);
  });
};

export default createChannelWithOrderer;
