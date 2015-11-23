'use strict';
let ConnectServer = function(conf) {

  function connectFn() {
    const connect = require('gulp-connect'),
      modRewrite = require('connect-modrewrite'),
      fs = require('fs-extra');

    const libPath = '/libs/';

    function copyLibaryFile(next, req) {
      const relativeFile = req.url.substr(libPath.length - 1).split('?')[0];
      const fileTo = conf.outputLibs + relativeFile;
      let notFound;

      conf.libaryFolders.forEach((libPath, i, array) => {
          const fileFrom = libPath + relativeFile;
          if(i === 0) { notFound = array.length; }
          if(notFound === -1) { return; }

          fs.copy(fileFrom, fileTo, (err) => {
            if (err === null) {
              console.log('copying: ' + fileFrom + ' to ' + fileTo);
              notFound = -1;
              next();
            } else if(--notFound === 0) {
              console.log('not found: ' + fileFrom);
              console.log('error', err);
              next();
            }
          });
        });
    }

    function libFinder(req, res, next) {
      if(req.url.startsWith(libPath)) {
        copyLibaryFile(next, req);
      } else {
        next();
      }
    }

    connect.server({
      root: ['/'],
      port: 9000,
      // livereload: true,
      middleware: function (connect, opt) {
        let proxy = modRewrite([
          '^/(.*)$ http://localhost:8080/$1 [P]'
        ]);

        return [libFinder, proxy];
      }
    });
  }

  return {
    connect: connectFn
  };
};

module.exports = ConnectServer;