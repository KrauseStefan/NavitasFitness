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
    if(req.url.startsWith('/rest')){
      modRewrite(['^/(.*)$ http://localhost:8080/$1 [P]'])(req, res, next);
    } else {
      next();
    }
  })

  .use(serveStatic(serverFolder))

  .use((req, res, next) => {
    console.log(req.url);
    if (req.url.startsWith(libPath)) {
      copyLibaryFile(next, req);
    } else {
      next();
    }
  })
  .listen(serverPort, () => {
    console.log('Server running on: ' + serverPort);
  });


function copyLibaryFile(next, req) {
  const relativeFile = req.url.substr(libPath.length - 1).split('?')[0];
  const fileTo = outputLibs + relativeFile;
  let notFound;

  libaryFolders.forEach((libPath, i, array) => {
    const fileFrom = libPath + relativeFile;
    if (i === 0) { notFound = array.length; }
    if (notFound === -1) { return; }

    fs.copy(fileFrom, fileTo, (err) => {
      if (err === null) {
        //              console.log('copying: ' + fileFrom + ' to ' + fileTo);
        notFound = -1;
        next();
      } else if (--notFound === 0) {
        console.log('not found: ' + fileFrom);
        console.log('error', err);
        next();
      }
    });
  });
}

// let ConnectServer = function(conf) {

//   function startConnectServer() {



//     function libFinder(req, res, next) {
//       if(req.url.startsWith(libPath)) {
//         copyLibaryFile(next, req);
//       } else {
//         next();
//       }
//     }

//     connect.server({
//       root: ['/'],
//       port: 9000,
//       livereload: true,
//       middleware: function (connect, opt) {
//         let proxy = modRewrite([
//           '^/(.*)$ http://localhost:8080/$1 [P]'
//         ]);

//         return [libFinder, proxy];
//       }
//     });
//   }

//   return {
//     startConnectServer: startConnectServer,
//     reload: connect.reload
//   };
// };

// module.exports = ConnectServer;