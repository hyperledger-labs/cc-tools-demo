import ccRouter from './router';

import logger = require('morgan');
import bodyParser = require('body-parser');
import * as cors from 'cors';
import * as swaggerUi from 'swagger-ui-express';
import fs = require('fs')
import yamljs = require('yamljs');

import express = require('express');

const swaggerFile = process.env.DOCKER ? '/rest-server/swagger.yaml' : './swagger.yaml';
const swaggerDoc = yamljs.load(swaggerFile);

// Creates and configures an ExpressJS web server.
class App {
  // ref to Express instance
  public express: express.Application;

  // Run configuration methods on the Express instance.
  constructor() {
    this.express = express();
    this.middleware();
    this.routes();
  }

  // Configure Express middleware.
  private middleware(): void {
    this.express.use(logger('dev'));
    this.express.use(bodyParser.json());
    this.express.use(bodyParser.urlencoded({ extended: true }));
  }

  // Configure API endpoints.
  private routes(): void {
    const options: cors.CorsOptions = {
      allowedHeaders: ['Origin', 'X-Requested-With', 'Content-Type', 'Accept', 'X-Access-Token', 'Authorization'],
      credentials: true,
      methods: 'GET,POST,PUT,DELETE',
      origin: '*',
      preflightContinue: false,
    };

    // use cors middleware
    this.express.use(cors(options));

    // enable pre-flight
    this.express.options('*', cors(options));

    /* This is just to get up and running, and to make sure what we've got is
     * working so far. This function will change when we start to add more
     * API endpoints */
    const router = express.Router();
    // placeholder route handler
    router.get('/', (req, res, next) => {
      res.json({
        message: 'online',
      });
    });
    this.express.use('/', router);
    this.express.use('/api/', ccRouter);
    this.express.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerDoc));
  }
}

export default new App().express;
