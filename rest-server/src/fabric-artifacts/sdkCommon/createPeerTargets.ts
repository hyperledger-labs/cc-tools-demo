import { PeerTarget } from '../../interfaces';
import * as fs from 'fs';
import * as path from 'path';
import Client = require('fabric-client');

const createPeerTargets = (peers: PeerTarget[]) => {
  const isDocker = process.env.DOCKER;
  const tlsPeer = isDocker
    ? path.resolve('/certs/ca.pem')
    : path.resolve('../fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls');
  const pem = fs.readFileSync(tlsPeer).toString();
  const targets: Client.Peer[] = peers.map(peer => {
    const opts: any = {
      'ssl-target-name-override': peer.name,
      pem,
      'request-timeout': 60000
    };
    return new Client.Peer(`grpcs://${peer.ip}:7051`, opts);
  });
  return targets;
};

export default createPeerTargets;
