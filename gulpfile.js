'use strict';

const config = require( './buildenv/gulpfile.config'),
    connectServer = require('./buildenv/connect.server')(config),

    gulp = require('gulp'),
    inject = require('gulp-inject'),
    ts = require('gulp-typescript'),
    gulpSass = require('gulp-sass'),
    gulpTslint = require('gulp-tslint'),
    sourcemaps = require('gulp-sourcemaps'),
    gulpJade = require('gulp-jade'),
    webpack = require('webpack'),

    del = require('del'),
    spawn = require('child_process').spawn,

    tsProject = ts.createProject( './websrc/tsconfig.json', {
      typescript: require('typescript')
    });

/**
 * Generates the app.d.ts references file dynamically from all application *.ts files.
 */
gulp.task(genTsRefs);
function genTsRefs() {
  const target = gulp.src(config.appTypeScriptReferences);
  const sources = gulp.src([config.allTypeScript], {read: false});

  return target.pipe(inject(sources, {
    starttag: '//{',
    endtag: '//}',
    transform: function (filepath) {
      return '/// <reference path="../..' + filepath + '" />';
    }
  })).pipe(gulp.dest(config.typings));
}

/**
 * Lint all custom TypeScript files.
 */
function tsLint() {
  return gulp.src(config.allTypeScript)
    .pipe(gulpTslint({
      configuration: {
        rules: {
          'class-name': true
        }
      }
    }))
    .pipe(gulpTslint.report('prose'));
}

/**
 * Compile TypeScript and include references to library and app .d.ts files.
 */
function compileTs(callback) {

  webpack(require('./websrc/webpack.config'), function(err, stats) {
    callback(err);
  });

  // const tsResult = tsProject.src()
  //   .pipe(sourcemaps.init())
  //   .pipe(ts(tsProject));

  // tsResult.dts.pipe(gulp.dest(config.outputPath));
  // return tsResult.js
  //   .pipe(sourcemaps.write('.'))
  //   .pipe(gulp.dest(config.outputPath));
}

gulp.task(clean);
function clean(done) {
  // delete the files
  return del([`${config.outputPath}/**/*`]);
}

gulp.task(watch);
function watch() {
  // const config = {
  //   config.allTypeScript: buildTs
  // }
  let conf = new Map();
  conf.set(config.allTypeScript, gulp.series(buildTs, liveReload));
  conf.set(config.views, gulp.series(jade, liveReload));
  conf.set(config.styles, sass);

  conf.forEach((value, key) => {
    gulp.watch(key, value);
  });

  // return gulp.watch([config.allTypeScript, config.views, config.styles], build);
}

function liveReload(done) {
  connectServer.reload();
  done();
}

function jade() {
  return gulp.src(config.views)
    .pipe(gulpJade({
      pretty: true
    }))
    .pipe(gulp.dest(config.outputPath));
}

function runGoappCmd(cmd, done) {

  const goapp = spawn('bash', ['-c',`${cmd} `], {
    stdio: 'inherit',
    cwd: './app-engine'
  });

  goapp.on('close', function (exitCode) {
    if(exitCode !== 0){
      console.log('child process exited with code ' + exitCode);
    }
    done();
  });

  goapp.on('error', function (err) {
    console.log('Failed to start child process: ', err);
    done();
  });
}

gulp.task(serve);
function serve(done) {
  runGoappCmd('dev_appserver.py --dev_appserver_log_level error .', done);
}

gulp.task(deploy);
function deploy(done) {
  runGoappCmd('goapp deploy', done);
}

gulp.task(buildTs);
function buildTs(done) {
  return gulp.series(tsLint, compileTs)(done);
}

gulp.task(build);
function build(done) {
  return gulp.parallel(
    buildTs,
    sass,
    jade)(done);
}

gulp.task(connect);
function connect() {
  return connectServer.startConnectServer();
}

gulp.task(sass);
function sass() {
  return gulp.src(config.styles)
    .pipe(gulpSass().on('error', gulpSass.logError))
    .pipe(gulp.dest(`${config.outputPath}/styles`));
}

gulp.task('buildAndWatch', gulp.parallel(gulp.series(build, gulp.parallel(connect, watch))));
gulp.task('default', gulp.parallel(gulp.series(build, gulp.parallel(connect, watch)), serve));

