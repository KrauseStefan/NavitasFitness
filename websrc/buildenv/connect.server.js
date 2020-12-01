'use strict';

const connect = require('connect');
const modRewrite = require('connect-modrewrite');
const serveStatic = require('serve-static');

const serverPort = 9000;

// below paths are relative to this script
const appServerFolder = __dirname + '/../../src/NavitasFitness/webapp';
const srcServeFolder = __dirname + '/../app';

connect()
  .use('/src', serveStatic(srcServeFolder))
  .use((req, res, next) => {
    if (req.url.startsWith('/rest')) {
      modRewrite(['^/(.*)$ http://localhost:8080/$1 [P]'])(req, res, next);
    } else if (req.url.match(/\./) === null && req.url !== '') {
      console.log(`redirecting ${ req.url } to /`)
      req.url = '/'
      next();
    } else {
      next();
    }
  })

  .use(serveStatic(appServerFolder))

  .listen(serverPort, () => {
    console.log('Server running on: ' + serverPort);
  });