'use strict';

var gulp = require('gulp'),
    debug = require('gulp-debug'),
    inject = require('gulp-inject'),
    tsc = require('gulp-typescript'),
    tslint = require('gulp-tslint'),
    sourcemaps = require('gulp-sourcemaps'),
    del = require('del'),
    Config = require('./gulpfile.config'),
    jade = require('gulp-jade'),
//    browserSync = require('browser-sync'),
//    superstatic = require( 'superstatic' ),
    tsProject = tsc.createProject('tsconfig.json');

var config = new Config();

/**
 * Generates the app.d.ts references file dynamically from all application *.ts files.
 */
 gulp.task('gen-ts-refs', function () {
     var target = gulp.src(config.appTypeScriptReferences);
     var sources = gulp.src([config.allTypeScript], {read: false});
     return target.pipe(inject(sources, {
         starttag: '//{',
         endtag: '//}',
         transform: function (filepath) {
             return '/// <reference path="../..' + filepath + '" />';
         }
     })).pipe(gulp.dest(config.typings));
 });

/**
 * Lint all custom TypeScript files.
 */
gulp.task('ts-lint', function () {
  return gulp.src(config.allTypeScript)
    .pipe(tslint({
      configuration: {
        rules: {
          'class-name': true
        }
      }
    }))
    .pipe(tslint.report('prose'));
});

/**
 * Compile TypeScript and include references to library and app .d.ts files.
 */
gulp.task('compile-ts', function () {
    var sourceTsFiles = [config.allTypeScript,                //path to typescript files
                         config.libraryTypeScriptDefinitions]; //reference to library .d.ts files
                        

    var tsResult = gulp.src(sourceTsFiles)
                       .pipe(sourcemaps.init())
                       .pipe(tsc(tsProject));

        tsResult.dts.pipe(gulp.dest(config.outputPath));
        return tsResult.js
                        .pipe(sourcemaps.write('.'))
                        .pipe(gulp.dest(config.outputPath));
});

/**
 * Remove all generated JavaScript files from TypeScript compilation.
 */
gulp.task('clean-ts', function (cb) {
  var typeScriptGenFiles = [
                              config.outputPath +'/**/*.js',    // path to all JS files auto gen'd by editor
                              config.outputPath +'/**/*.js.map', // path to all sourcemap files auto gen'd by editor
                              '!' + config.outputPath + '/lib'
                           ];

  // delete the files
  del(typeScriptGenFiles, cb);
});

gulp.task('watch', function() {
    gulp.watch([config.allTypeScript], ['build']);
});

//gulp.task('serve', ['compile-ts', 'watch'], function() {
//  process.stdout.write('Starting browserSync and superstatic...\n');
//  browserSync({
//    port: 3000,
//    files: ['index.html', '**/*.js'],
//    injectChanges: true,
//    logFileChanges: false,
//    logLevel: 'silent',
//    logPrefix: 'angularin20typescript',
//    notify: true,
//    reloadDelay: 0,
//    server: {
//      baseDir: './src',
//      middleware: superstatic({ debug: false})
//    }
//  });
//});


gulp.task('jade', function() {

  gulp.src(config.views)
    .pipe(jade({
      pretty: true
    }))
    .pipe(gulp.dest(config.outputPath))
});

gulp.task('default', ['build']);

gulp.task('build', ['ts-lint', 'compile-ts', 'jade'/*, 'gen-ts-refs'*/])
