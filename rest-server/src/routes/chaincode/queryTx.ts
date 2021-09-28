import { Request, Response, NextFunction } from 'express';
import Client from '../createClient';
import query from '../../fabric-artifacts/chaincode/query';

const queryTx = (req: Request, res: Response, next: NextFunction) => {
  const client = Client.get();

  let transientMap: Object = {};
  let txRequest: Object = {};

  // Get transaction name
  const { tx } = req.params;

  let request: Object = {};

  if (req.method === "POST") {
    request = req.body
  } else if (req.method === "GET") {
    // Handle req.query
    let queryRequest
    if (req.query['@request']) {
      queryRequest = req.query['@request']
    } else if ((req.query['request']) ){
      queryRequest = req.query['request']
    }
    if (queryRequest) {
      try {
        const requestJSON = Buffer.from(req.query['@request'] as string, 'base64').toString('utf-8');
        request = JSON.parse(requestJSON);
      } catch (err) {
        console.error(err);
        return res.status(400).send('the @request query parameter must be a valid base64-encoded JSON object');
      }
    }
  }

  // Handle request
  for (const key in request) {
    if (key.charAt(0) == '~') {
      transientMap[key.slice(1, key.length)] = request[key];
    } else {
      txRequest[key] = request[key];
    }
  }

  // Encode tx arguments
  const txArgs = [JSON.stringify(txRequest)];
  const transientRequest = new Buffer(JSON.stringify(transientMap));

  // Query chaincode
  query(client, tx, txArgs, transientRequest)
    .then((response) => {
      return res.send(response);
    })
    .catch((err) => {
      if (err && err.status) {
        return res.status(err.status).send(err.message);
      }
      return res.status(500).send(err);
    });
};

export default queryTx;
