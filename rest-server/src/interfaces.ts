import fabprotos = require('./fabric-tools/bundle');

export interface PeerResponseError {
  message: string;
  status: number;
  payload: Buffer;
  isProposalResponse: Boolean;
}

export interface PeerJSON {
  peerName: {
    publicIP: string;
    endorsingPeer: boolean;
    chaincodeQuery: boolean;
    ledgerQuery: boolean;
    eventSource: boolean;
  };
}

export interface PeerTarget {
  name: string;
  ip: string;
}

export interface MSPRole extends fabprotos.common.MSPRole {
  msp_identifier?: string;
}

export interface IConfigPolicy extends fabprotos.common.IConfigPolicy {
  mod_policy?: string;
}
