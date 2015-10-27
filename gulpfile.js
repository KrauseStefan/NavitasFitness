'use strict';

const buildenvFolder = './websrc/buildenv/';

const gulp = require('gulp'),
    config = require( buildenvFolder + 'gulpfile.config'),
    connectServer = require(buildenvFolder + 'connect.server')(config),
    inject = require('gulp-inject'),
    ts = require('gulp-typescript'),
    gulpTslint = require('gulp-tslint'),
    sourcemaps = require('gulp-sourcemaps'),
    del = require('del'),
    gulpJade = require('gulp-jade'),
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

function compileTs() {
  const tsResult = tsProject.src()
    .pipe(sourcemaps.init())
    .pipe(ts(tsProject));

  tsResult.dts.pipe(gulp.dest(config.outputPath));
  return tsResult.js
    .pipe(sourcemaps.write('.'))
    .pipe(gulp.dest(config.outputPath));
}

gulp.task(clean);
function clean(done) {
  const typeScriptGenFiles = [
    config.outputPath +'/**/*.js',    // path to all JS files auto gen'd by editor
    config.outputPath +'/**/*.js.map', // path to all sourcemap files auto gen'd by editor
    '!' + config.outputPath + '/lib'
  ];

  // delete the files
  del(typeScriptGenFiles, done);
}

gulp.task(watch);
function watch() {
  return gulp.watch([config.allTypeScript, config.views], build);
}

function jade() {
  return gulp.src(config.views)
    .pipe(gulpJade({
      pretty: true
    }))
    .pipe(gulp.dest(config.outputPath));
}

gulp.task(serve);
function serve(done) {
  const goapp = spawn('bash', ['-c','goapp serve'], {
    stdio: 'inherit'
  });

  goapp.on('close', function (code) {
    console.log('child process exited with code ' + code);
    done();
  });

  goapp.on('error', function (err) {
    console.log('Failed to start child process.');
  });
}

gulp.task(build);
function build(done) {
  return gulp.series(tsLint, compileTs, jade)(done);
}

gulp.task(connect);
function connect() {
  return connectServer.connect(config);
}

gulp.task('test', gulp.parallel(connect));

gulp.task('default', gulp.series(build, gulp.parallel(connect, serve, watch)));

