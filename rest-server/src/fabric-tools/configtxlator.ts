import {protoConversion} from './interfaces';

import _ = require('lodash');
import camelCaseDeep = require('camelcase-keys-deep');
import fabprotos = require('../fabric-tools/bundle');
import {MSPRole, IConfigPolicy} from '../interfaces';

class Configtxlator {
  public createNewConsenter = (hostName: string, certificate: Uint8Array) : fabprotos.etcdraft.Consenter => {
      
    const consenterObj : fabprotos.etcdraft.IConsenter = {clientTlsCert: certificate, host: hostName, port: 7050, serverTlsCert: certificate};
    const errMsg = fabprotos.etcdraft.Consenter.verify(consenterObj);
    if (errMsg) {
      console.debug(`Error creating etcdraft consenter: ${errMsg}`)
      throw errMsg;
    }
    return fabprotos.etcdraft.Consenter.create(consenterObj);
  }


  //The block's json will return a json where keys are snake_case, but proto require they
  //to be camelCase, so this function will return a json encoded ready to proto.create
  public convertBlockToProtoJSON = (block: any) : any => JSON.parse(
      block.encodeJSON().replace(
          /(_\w)\w+":/g,
          (match: String) => match[1].toUpperCase() + match.substring(2)),
      (key, value) =>
          (key === 'sequence' ? parseInt(value) :
                                key === 'version' ? parseInt(value) : value));


  public computeDeltaSet =
      (original: fabprotos.common.Config,
       updated: fabprotos.common
           .Config) => new Promise<fabprotos.common
                                       .ConfigUpdate>((resolve, reject) => {
        // These errors messages are due lack of index.d.ts in fabprotos
        // package, once FABN-1145 were accepted this class types should be
        // rewriten
        if (!original.channelGroup) {
          reject(new Error('no channel group included for original config'));
        }

        if (!updated.channelGroup) {
          reject(new Error('no channel group included for updated config'));
        }

        const [readSet, writeSet, groupUpdated] = this.computeGroupUpdate(
            original.channelGroup, updated.channelGroup);

        if (!groupUpdated) {
          reject(new Error(
              'no differences detected between original and updated config'));
        }

        resolve(fabprotos.common.ConfigUpdate.create({readSet, writeSet}));
      });

  private computePoliciesMapUpdate = (original, updated): [
    {[k: string]: fabprotos.common.IConfigPolicy},
    {[k: string]: fabprotos.common.IConfigPolicy},
    {[k: string]: fabprotos.common.IConfigPolicy}, boolean
  ] => {
    const readSet: {[k: string]: fabprotos.common.IConfigPolicy} = {};
    const writeSet: {[k: string]: fabprotos.common.IConfigPolicy} = {};

    // All modified config goes into the read/write sets, but in case the map
    // membership changes, we retain the config which was the same to add to the
    // read/write sets
    const sameSet: {[k: string]: fabprotos.common.IConfigPolicy} = {};
    let updatedMembers = false;
    protoConversion.forEach(original, (originalPolicy, policyName) => {
      const updatedPolicy = updated[policyName];
      if (!updatedPolicy) {
        updatedMembers = true;
        return;
      }

      const originalValue = Buffer.from(originalPolicy.policy.value, 'base64');

      const updatedValue = Buffer.from(updatedPolicy.policy.value, 'base64');
      if (originalPolicy.modPolicy === updatedPolicy.modPolicy &&
          originalValue.equals(updatedValue)) {
        sameSet[policyName] = fabprotos.common.ConfigPolicy.create(
            {version: originalPolicy.version});
        return;
      }

      const errMsg = fabprotos.common.Policy.verify(updatedPolicy.policy);
      if (errMsg) throw errMsg;
      const policy = fabprotos.common.Policy.create(updatedPolicy.policy);
      writeSet[policyName] = fabprotos.common.ConfigPolicy.create({
        version: +originalPolicy.version + 1,
        modPolicy: updatedPolicy.modPolicy,
        policy
      });
    });

    protoConversion.forEach(updated, (updatedPolicy, policyName) => {
      if (original[policyName]) {
        // If the updatedPolicy is in the original set of policies, it was
        // already handled
        return;
      }
      updatedMembers = true;
      writeSet[policyName] = fabprotos.common.ConfigPolicy.create({
        version: 0,
        modPolicy: updatedPolicy.modPolicy,
        policy: updatedPolicy.policy
      });
    });
    return [readSet, writeSet, sameSet, updatedMembers];
  };

  private computeValuesMapUpdate = (original, updated): [
    {[k: string]: fabprotos.common.IConfigValue},
    {[k: string]: fabprotos.common.IConfigValue},
    {[k: string]: fabprotos.common.IConfigValue}, boolean
  ] => {
    const readSet: {[k: string]: fabprotos.common.IConfigValue} = {};
    const writeSet: {[k: string]: fabprotos.common.IConfigValue} = {};

    // All modified config goes into the read/write sets, but in case the map
    // membership changes, we retain the config which was the same to add to the
    // read/write sets
    const sameSet: {[k: string]: fabprotos.common.IConfigValue} = {};
    let updatedMembers = false;

    protoConversion.forEach(original, (originalValue, valueName) => {
      const updatedValue = updated[valueName];
      if (!updatedValue) {
        updatedMembers = true;
        return;
      }
      const valueOriginal = Buffer.from(originalValue.value, 'base64');

      const valueUpdated = Buffer.from(updatedValue.value, 'base64');
      if (originalValue.modPolicy === updatedValue.modPolicy &&
          valueOriginal.equals(valueUpdated)) {
        sameSet[valueName] = fabprotos.common.ConfigValue.create(
            {version: originalValue.version});
        return;
      }

      writeSet[valueName] = fabprotos.common.ConfigValue.create({
        version: +originalValue.version + 1,
        modPolicy: updatedValue.modPolicy,
        value: updatedValue.value
      });
    });

    protoConversion.forEach(updated, (updatedValue, valueName) => {
      if (original[valueName]) {
        // If the updatedValue is in the original set of values, it was already
        // handled
        return;
      }
      updatedMembers = true;
      writeSet[valueName] = fabprotos.common.ConfigValue.create({
        version: 0,
        modPolicy: updatedValue.modPolicy,
        value: updatedValue.value
      });
    });

    return [readSet, writeSet, sameSet, updatedMembers];
  };

  private computeGroupsMapUpdate =
      (original: {[k: string]: fabprotos.common.IConfigGroup},
       updated: {[k: string]: fabprotos.common.IConfigGroup}):
          [
            {[k: string]: fabprotos.common.IConfigGroup},
            {[k: string]: fabprotos.common.IConfigGroup},
            {[k: string]: fabprotos.common.IConfigGroup}, boolean
          ] => {
            const readSet: {[k: string]: fabprotos.common.IConfigGroup} = {};
            const writeSet: {[k: string]: fabprotos.common.IConfigGroup} = {};

            // All modified config goes into the read/write sets, but in case
            // the map membership changes, we retain the config which was the
            // same to add to the read/write sets
            const sameSet: {[k: string]: fabprotos.common.IConfigGroup} = {};
            let updatedMembers = false;
            protoConversion.forEach(original, (originalGroup, groupName) => {
              const updatedGroup = updated[groupName];
              if (!updatedGroup) {
                updatedMembers = true;
                return;
              }

              const [groupReadSet, groupWriteSet, groupUpdated] =
                  this.computeGroupUpdate(originalGroup, updatedGroup);
              if (!groupUpdated) {
                sameSet[groupName] =
                    fabprotos.common.ConfigGroup.create(originalGroup);
                sameSet[groupName].modPolicy = originalGroup.modPolicy;
                return;
              }

              readSet[groupName] = groupReadSet;
              writeSet[groupName] = groupWriteSet;
            });

            protoConversion.forEach(updated, (updatedGroup, groupName) => {
              if (original[groupName]) {
                // If the updatedGroup is in the original set of groups, it was
                // already handled
                return;
              }
              updatedMembers = true;
              const [_, groupWriteSet, __] = this.computeGroupUpdate(
                  fabprotos.common.ConfigGroup.create({}), updatedGroup);

              writeSet[groupName] = fabprotos.common.ConfigGroup.create({
                version: 0,
                modPolicy: updatedGroup.modPolicy,
                policies: groupWriteSet.policies,
                values: groupWriteSet.values,
                groups: groupWriteSet.groups
              });
            });

            return [readSet, writeSet, sameSet, updatedMembers];
          };

  private computeGroupUpdate =
      (original: fabprotos.common.IConfigGroup, updated: fabprotos.common.IConfigGroup): [
        fabprotos.common.IConfigGroup, fabprotos.common.IConfigGroup, boolean
      ] => {
        const [readSetPolicies, writeSetPolicies, sameSetPolicies, policiesMembersUpdated] =
            this.computePoliciesMapUpdate(original.policies, updated.policies);
        const [readSetValues, writeSetValues, sameSetValues, valuesMembersUpdated] =
            this.computeValuesMapUpdate(original.values, updated.values);
        const [readSetGroups, writeSetGroups, sameSetGroups, groupsMembersUpdated] =
            this.computeGroupsMapUpdate(original.groups, updated.groups);

        // If the updated group is 'Equal' to the updated group (none of the
        // members nor the mod policy changed)
        if (!(policiesMembersUpdated || valuesMembersUpdated ||
              groupsMembersUpdated ||
              original.modPolicy != updated.modPolicy)) {
          // If there were no modified entries in any of the
          // policies/values/groups maps
          if (protoConversion.size(readSetPolicies) === 0 &&
              protoConversion.size(writeSetPolicies) === 0 &&
              protoConversion.size(readSetValues) === 0 &&
              protoConversion.size(writeSetValues) === 0 &&
              protoConversion.size(readSetGroups) === 0 &&
              protoConversion.size(writeSetGroups) === 0) {
            return [
              fabprotos.common.ConfigGroup.create({version: original.version}),
              fabprotos.common.ConfigGroup.create({version: original.version}),
              false
            ];
          }

          return [
            fabprotos.common.ConfigGroup.create({
              version: original.version,
              policies: readSetPolicies,
              values: readSetValues,
              groups: readSetGroups
            }),
            fabprotos.common.ConfigGroup.create({
              version: original.version,
              policies: writeSetPolicies,
              values: writeSetValues,
              groups: writeSetGroups
            }),
            true
          ];
        }

        protoConversion.forEach(sameSetPolicies, (samePolicy, k) => {
          readSetPolicies[k] = samePolicy;
          writeSetPolicies[k] = samePolicy;
        });

        protoConversion.forEach(sameSetValues, (sameValue, k) => {
          readSetValues[k] = sameValue;
          writeSetValues[k] = sameValue;
        });

        protoConversion.forEach(sameSetGroups, (sameGroup, k) => {
          readSetGroups[k] = sameGroup;
          writeSetGroups[k] = sameGroup;
        });

        return [
          fabprotos.common.ConfigGroup.create({
            version: original.version,
            policies: readSetPolicies,
            values: readSetValues,
            groups: readSetGroups
          }),
          fabprotos.common.ConfigGroup.create({
            version: +original.version + 1,
            policies: writeSetPolicies,
            values: writeSetValues,
            groups: writeSetGroups,
            modPolicy: updated.modPolicy
          }),
          true
        ];
      };

  // This is gambiarra
  // Blame fabric and protobuf.js not me
  private convertFabMSPmsg = (fabMSPmsg: any) => ({
    rootCerts: fabMSPmsg.root_certs,
    intermediateCerts: fabMSPmsg.intermediate_certs,
    admins: fabMSPmsg.admins,
    revocationList: fabMSPmsg.revocation_list,
    organizationalUnitIdentifiers: fabMSPmsg.organizational_unit_identifiers,
    tlsRootCerts: fabMSPmsg.tls_root_certs,
    tlsIntermediateCerts: fabMSPmsg.tls_intermediate_certs,
    name: fabMSPmsg.name
  });

  // SignaturePolicy is a recursive message structure which defines a
  // featherweight DSL for describing policies which are more complicated than
  // 'exactly this signature'.  The NOutOf operator is sufficent to express AND
  // as well as OR, as well as of course N out of the following M policies
  // SignedBy implies that the signature is from a valid certificate which is
  // signed by the trusted authority specified in the bytes.  This will be the
  // certificate itself for a self-signed certificate and will be the CA for
  // more traditional certificates
  private createRuleRecursively =
      (rule: any): fabprotos.common.SignaturePolicy|
      fabprotos.common.SignaturePolicy[] => {
        let newRule = fabprotos.common.SignaturePolicy.create(rule);
        Object.keys(newRule).forEach(key => {
          if (key === 'nOutOf') {
            newRule.nOutOf.rules = this.createRuleRecursively(rule[key]) as
                fabprotos.common.SignaturePolicy[];
          }
          if (key === 'rules' && rule[key] instanceof Array) {
            newRule = rule[key].map(elem => this.createRuleRecursively(elem));
          }
        });
        return newRule;
      };

  public convertOrgJsonToConfigGroup = (orgJson: any) => {
    const eachRecursive = (obj: any) => {
      // each recursive iterate recursively down to the object changing its
      // necessary keys to proto expected ones
      for (const k in obj) {
        if (k === 'version' && typeof obj[k] === 'string') {
          obj[k] = parseInt(obj[k]);
        } else if (typeof obj[k] === 'object' && k === 'values') {
          eachRecursive(obj[k]);
        } else if (k === 'msp' && typeof obj[k] === 'object') {
          if (obj[k].value.config) {
            obj[k].version = parseInt(obj[k].version);

            const errMsg =
                fabprotos.msp.FabricMSPConfig.verify(obj[k].value.config);
            if (errMsg) throw new Error('Invalid MSP config');

            const fabMSPConfig =
                fabprotos.msp.FabricMSPConfig.create(obj[k].value.config);
            const mspConfig = fabprotos.msp.MSPConfig.create({
              config:
                  fabprotos.msp.FabricMSPConfig.encode(fabMSPConfig).finish()
            });

            obj[k].value = fabprotos.msp.MSPConfig.encode(mspConfig).finish();
          }

          obj.MSP = fabprotos.common.ConfigValue.create(obj[k]);
          delete obj.msp;
          continue;
        } else if (k === 'policy' && typeof obj[k] === 'object') {
          const sigPolEnvelope =
              fabprotos.common.SignaturePolicyEnvelope.create(obj[k].value);
          sigPolEnvelope
              .identities = sigPolEnvelope.identities.map(identity => {
            const castedPrincipal = (identity.principal as unknown) as MSPRole;

            const principalEnum =
                (fabprotos.common.MSPRole.MSPRoleType[castedPrincipal.role] as
                 unknown) as fabprotos.common.MSPRole.MSPRoleType;
            const principal = fabprotos.common.MSPRole.create({
              mspIdentifier: castedPrincipal.mspIdentifier,
              role: principalEnum
            });

            const errMsg = fabprotos.common.MSPRole.verify(principal);
            if (errMsg)
              throw new Error(
                  `Error verifying MSPRole on translate newOrg ${errMsg}`);

            const principalBytes =
                fabprotos.common.MSPRole.encode(principal).finish();
            const principalClassification =
                (fabprotos.common.MSPPrincipal
                     .Classification[identity.principalClassification] as
                 unknown) as fabprotos.common.MSPPrincipal.Classification;

            return fabprotos.common.MSPPrincipal.create(
                {principal: principalBytes, principalClassification});
          });

          sigPolEnvelope.rule =
              this.createRuleRecursively(sigPolEnvelope.rule) as
              fabprotos.common.SignaturePolicy;

          const errMsg =
              fabprotos.common.SignaturePolicyEnvelope.verify(sigPolEnvelope);
          if (errMsg) throw errMsg;

          const sigPolEnvelopeBytes =
              fabprotos.common.SignaturePolicyEnvelope.encode(sigPolEnvelope)
                  .finish();
          const policy = fabprotos.common.Policy.create(
              {type: obj[k].type, value: sigPolEnvelopeBytes});

          obj[k] = policy;
          continue;
        }
        if (typeof obj[k] === 'object' && obj[k] !== null)
          eachRecursive(obj[k]);
      }
    };
    const convertedValueGroup = camelCaseDeep(orgJson.values);
    const desiredFromJSON =
        _.pick(orgJson, ['groups', 'mod_policy', 'policies']);
    const modPolicy = desiredFromJSON.mod_policy;
    delete desiredFromJSON.mod_policy;
    const translated = fabprotos.common.ConfigGroup.create(
        {...desiredFromJSON, values: convertedValueGroup, modPolicy});
    for (const k in translated.policies) {
      let policy = translated.policies[k] as IConfigPolicy;
      policy = camelCaseDeep(policy);
      translated.policies[k] = fabprotos.common.ConfigPolicy.create(policy);
    }
    eachRecursive(translated);

    const errMsg = fabprotos.common.ConfigGroup.verify(translated);
    if (errMsg)
      throw new Error(
          `Could not convert new OrgJson failed with error ${errMsg}`);
    return translated;
  };
}

export default Configtxlator;
