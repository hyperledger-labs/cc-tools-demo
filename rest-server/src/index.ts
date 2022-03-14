import * as fs from 'fs';
import * as https from 'https';
import * as http from 'http';

import App from './App';

import debug = require('debug');

debug('ts-express:server');

let server;
let port;

if (process.env.HTTPS === 'true') {
  const key = process.env.KEYFILE || `/rest-server/data/certbot/conf/live/${process.env.DOMAIN}/privkey.pem`;
  const cert = process.env.CERTFILE || `/rest-server/data/certbot/conf/live/${process.env.DOMAIN}/fullchain.pem`;

  const serverOptions = {
    key: fs.readFileSync(key),
    cert: fs.readFileSync(cert),
  };

  port = normalizePort(process.env.PORT || 443);
  App.set('port', port);
  server = https.createServer(serverOptions, App);
} else {
  const defaultPort = process.env.DEPLOY ? 80 : (process.env.DOCKER ? 80 : 8000);
  port = normalizePort(process.env.PORT || defaultPort);

  App.set('port', port);
  server = http.createServer(App);
}

server.listen(port);
server.on('error', onError);
server.on('listening', onListening);

function normalizePort(val: number | string): number | string | boolean {
  const port: number = typeof val === 'string' ? parseInt(val, 10) : val;
  if (isNaN(port)) return val;
  if (port >= 0) return port;
  return false;
}

function onError(error: NodeJS.ErrnoException): void {
  if (error.syscall !== 'listen') throw error;
  const bind = typeof port === 'string' ? `Pipe ${port}` : `Port ${port}`;
  switch (error.code) {
    case 'EACCES':
      console.error(`${bind} requires elevated privileges`);
      process.exit(1);
      break;
    case 'EADDRINUSE':
      console.error(`${bind} is already in use`);
      process.exit(1);
      break;
    default:
      throw error;
  }
}

function onListening(): void {
  const addr = server.address();
  const bind = typeof addr === 'string' ? `pipe ${addr}` : `port ${addr.port}`;
  console.log(`Listening on ${bind}`);
}
