import { common } from 'fabric-protos';

export namespace protoConversion {
  export interface IConfigPolicyMap {
    [k: string]: common.IConfigPolicy;
  }

  export interface IConfigValueMap {
    [k: string]: common.IConfigValue;
  }

  export interface IConfigGroupMap {
    [k: string]: common.IConfigGroup;
  }

  export enum MSPRoleType {
    MEMBER, // Represents an MSP Member
    ADMIN, // Represents an MSP Admin
    CLIENT, // Represents an MSP Client
    PEER, // Represents an MSP Peer
  }

  export enum Classification {
    ROLE, // Represents the one of the dedicated MSP roles, the
    // one of a member of MSP network, and the one of an
    // administrator of an MSP network
    ORGANIZATION_UNIT, // Denotes a finer grained (affiliation-based)
    // groupping of entities, per MSP affiliation
    // E.g., this can well be represented by an MSP's
    // Organization unit
    IDENTITY, // Denotes a principal that consists of a single
    // identity
    ANONYMITY, // Denotes a principal that can be used to enforce
    // an identity to be anonymous or nominal.
    COMBINED, // Denotes a combined principal
  }
  export const forEach = (
    collection: IConfigPolicyMap | IConfigValueMap | IConfigGroupMap,
    callback: (value, key?: string, object?) => any,
    scope?: IConfigPolicyMap | IConfigValueMap | IConfigGroupMap,
  ) => {
    if (Object.prototype.toString.call(collection) === '[object Object]') {
      for (const prop in collection) {
        if (Object.prototype.hasOwnProperty.call(collection, prop)) {
          callback.call(scope, collection[prop], prop, collection);
        }
      }
    }
  };

  export const size = (obj: object): Number => Object.keys(obj).length;
}
