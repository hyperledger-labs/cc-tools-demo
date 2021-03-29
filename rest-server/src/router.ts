import { Router } from 'express';

import * as chaincodeHandler from './routes/chaincode';

const ccRouter = (): Router => {
  const router = Router();

  router.post('/invoke/:tx', chaincodeHandler.invokeTx);
  router.put('/invoke/:tx', chaincodeHandler.invokeTx);
  router.delete('/invoke/:tx', chaincodeHandler.invokeTx);

  router.post('/query/:tx', chaincodeHandler.queryTx);
  router.get('/query/:tx', chaincodeHandler.queryTx);

  return router;
};

export default ccRouter();
