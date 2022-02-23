import { Router } from 'express';

import * as chaincodeHandler from './routes/chaincode';

const ccRouter = (): Router => {
  const router = Router();

  if (process.env.USEAUTH === "true") {
      router.use((req, res, next) => {
        if (!req.headers.authorization || req.headers.authorization.indexOf('Basic ') === -1) {
        return res.status(401).json({ message: 'Missing Authorization Header' });
        }

        // verify auth credentials
        const base64Credentials =  req.headers.authorization.split(' ')[1];
        const credentials = Buffer.from(base64Credentials, 'base64').toString('ascii');
        const [username, password] = credentials.split(':');
        
        if (username != process.env.USERNAME) {
          return res.status(401).json({ message: 'Credentials are invalid' });
        }
        
        if (password != process.env.PASSWORD) {
          return res.status(401).json({ message: 'Credentials are invalid' });
        }

        next()
    })
  }

  router.post('/invoke/:tx', chaincodeHandler.invokeTx);
  router.put('/invoke/:tx', chaincodeHandler.invokeTx);
  router.delete('/invoke/:tx', chaincodeHandler.invokeTx);

  router.post('/query/:tx', chaincodeHandler.queryTx);
  router.get('/query/:tx', chaincodeHandler.queryTx);

  return router;
};

export default ccRouter();
