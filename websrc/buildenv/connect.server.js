'use strict';

const connect = require('connect');
const modRewrite = require('connect-modrewrite');
const serveStatic = require('serve-static');
const fs = require('fs-extra');

const serverPort = 9000;
const libPath = '/libs/';
const serverFolder = __dirname + '/../../webapp';
const outputLibs = serverFolder + libPath;

const libaryFolders = [
  '/node_modules'
].map((lib) => __dirname + '/..' + lib);

connect()
  .use((req, res, next) => {
    if (req.url.startsWith('/rest')) {
      modRewrite(['^/(.*)$ http://localhost:8080/$1 [P]'])(req, res, next);
    } else if (req.url.match(/\./) === null && req.url !== '') {
      console.log(`redirecting ${req.url} to /`)
      req.url = '/'
      next();
    } else {
      next();
    }
  })

  .use(serveStatic(serverFolder))

  .listen(serverPort, () => {
    console.log('Server running on: ' + serverPort);
  });