import * as path from "path";
import Client = require("fabric-client");

const fab_modules = path.join("../../node_modules/fabric-client/lib");
const sdk_utils = require(path.join(fab_modules, "./utils.js"));
const client_utils = require(path.join(fab_modules, "./client-utils.js"));
const TransactionID = require(path.join(fab_modules, "./TransactionID.js"));
import ProtoLoader from "fabric-client/lib/ProtoLoader";
const logger = sdk_utils.getLogger("Channel.js");

const protos_folder = path.join(process.cwd(), "clientprotos");
const _commonProto = ProtoLoader.load(path.join(protos_folder + "/common/common.proto")).common;
const _configtxProto = ProtoLoader.load(protos_folder + "/common/configtx.proto").common;
const _abProto = ProtoLoader.load(protos_folder + "/orderer/ab.proto").orderer;

const getChannelConfigFromOrderer = async (channel: any) => {
  //Use with caution, channel MUST be a fabric Channel instance. It is as 'any' due to need to access
  //private attributes not exposed by Channel interface (typescript), although it is accessible through js
  const method = "getChannelConfigFromOrderer";

  const self = channel;
  const orderer = channel._clientContext.getTargetOrderer(null, channel.getOrderers(), channel._name);

  const signer = channel._clientContext._getSigningIdentity(true);
  let txId = new TransactionID(signer, true);

  // seek the latest block
  let seekSpecifiedStart = new _abProto.SeekNewest();
  let seekStart = new _abProto.SeekPosition();
  seekStart.setNewest(seekSpecifiedStart);

  let seekSpecifiedStop = new _abProto.SeekNewest();
  let seekStop = new _abProto.SeekPosition();
  seekStop.setNewest(seekSpecifiedStop);

  // seek info with all parts
  let seekInfo = new _abProto.SeekInfo();
  seekInfo.setStart(seekStart);
  seekInfo.setStop(seekStop);
  seekInfo.setBehavior(_abProto.SeekInfo.SeekBehavior.BLOCK_UNTIL_READY);

  // build the header for use with the seekInfo payload
  let seekInfoHeader = client_utils.buildChannelHeader(
    _commonProto.HeaderType.DELIVER_SEEK_INFO,
    self._name,
    txId.getTransactionID(),
    self._initial_epoch,
    null,
    client_utils.buildCurrentTimestamp(),
    channel._clientContext.getClientCertHash()
  );

  let seekHeader = client_utils.buildHeader(signer, seekInfoHeader, txId.getNonce());
  let seekPayload = new _commonProto.Payload();
  seekPayload.setHeader(seekHeader);
  seekPayload.setData(seekInfo.toBuffer());

  // building manually or will get protobuf errors on send
  let envelope = client_utils.toEnvelope(client_utils.signProposal(signer, seekPayload));
  // client will return us a block
  let block = await orderer.sendDeliver(envelope);
  logger.debug("%s - good results from seek block ", method); // :: %j',results);
  // verify that we have the genesis block
  if (block) {
    logger.debug("%s - found latest block", method);
  } else {
    logger.error("%s - did not find latest block", method);
    throw new Error(`Failed to retrieve latest block ${method}`);
  }

  logger.debug("%s - latest block is block number %s", block.header.number);
  // get the last config block number
  const metadata = _commonProto.Metadata.decode(block.metadata.metadata[_commonProto.BlockMetadataIndex.LAST_CONFIG]);
  const last_config = _commonProto.LastConfig.decode(metadata.value);
  logger.debug("%s - latest block has config block of %s", method, last_config.index);

  txId = new TransactionID(signer);

  // now build the seek info to get the block called out
  // as the latest config block
  seekSpecifiedStart = new _abProto.SeekSpecified();
  seekSpecifiedStart.setNumber(last_config.index);
  seekStart = new _abProto.SeekPosition();
  seekStart.setSpecified(seekSpecifiedStart);

  //   build stop
  seekSpecifiedStop = new _abProto.SeekSpecified();
  seekSpecifiedStop.setNumber(last_config.index);
  seekStop = new _abProto.SeekPosition();
  seekStop.setSpecified(seekSpecifiedStop);

  // seek info with all parts
  seekInfo = new _abProto.SeekInfo();
  seekInfo.setStart(seekStart);
  seekInfo.setStop(seekStop);
  seekInfo.setBehavior(_abProto.SeekInfo.SeekBehavior.BLOCK_UNTIL_READY);
  // logger.debug('initializeChannel - seekInfo ::' + JSON.stringify(seekInfo));

  // build the header for use with the seekInfo payload
  seekInfoHeader = client_utils.buildChannelHeader(
    _commonProto.HeaderType.DELIVER_SEEK_INFO,
    self._name,
    txId.getTransactionID(),
    self._initial_epoch,
    null,
    client_utils.buildCurrentTimestamp(),
    self._clientContext.getClientCertHash()
  );

  seekHeader = client_utils.buildHeader(signer, seekInfoHeader, txId.getNonce());
  seekPayload = new _commonProto.Payload();
  seekPayload.setHeader(seekHeader);
  seekPayload.setData(seekInfo.toBuffer());

  // building manually or will get protobuf errors on send
  envelope = client_utils.toEnvelope(client_utils.signProposal(signer, seekPayload));
  // client will return us a block
  block = await orderer.sendDeliver(envelope);
  if (!block) {
    throw new Error("Config block was not found");
  }
  // lets have a look at the block
  logger.debug(
    "%s -  config block number ::%s  -- numberof tx :: %s",
    method,
    block.header.number,
    block.data.data.length
  );
  if (block.data.data.length !== 1) {
    throw new Error("Config block must only contain one transaction");
  }

  return block;
};

export default getChannelConfigFromOrderer;
