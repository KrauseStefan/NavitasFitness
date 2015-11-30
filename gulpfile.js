'use strict';

const buildenvFolder = './buildenv/';

const gulp = require('gulp'),
    config = require( buildenvFolder + 'gulpfile.config'),
    connectServer = require(`${buildenvFolder}connect.server`)(config),
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
  // delete the files
  return del([`${config.outputPath}/**/*`]);
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

function runGoappCmd(cmd, done) {
  const goapp = spawn('bash', ['-c',`goapp ${cmd}`], {
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
  runGoappCmd('serve', done);
}

gulp.task(deploy);
function deploy(done) {
  runGoappCmd('deploy', done);
}

gulp.task(build);
function build(done) {
  return gulp.series(tsLint, compileTs, jade)(done);
}

gulp.task(connect);
function connect() {
  return connectServer.connect(config);
}

gulp.task('default', gulp.parallel(gulp.series(build, gulp.parallel(connect, watch)), serve));

