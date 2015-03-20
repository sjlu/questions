var gulp = require('gulp');
var stylus = require('gulp-stylus');
var sourcemaps = require('gulp-sourcemaps');
var nib = require('nib');
var watch = require('gulp-watch');
var livereload = require('gulp-livereload');
var less = require('gulp-less');
var minifyCSS = require('gulp-minify-css');
var plumber = require('gulp-plumber');
var templateCache = require('gulp-angular-templatecache');
var jade = require('gulp-jade');
var _ = require('lodash');
var uglify = require('gulp-uglifyjs');
var path = require('path');
var concat = require('gulp-concat');
var fs = require('fs');

var bowerjs = [
  'angular/angular.min.js',
  'angular-bootstrap/ui-bootstrap-tpls.min.js',
  'angular-route/angular-route.min.js',
  'jquery/dist/jquery.min.js',
  'lodash/lodash.min.js',
  'twitter/dist/js/bootstrap.min.js'
]

var bowercopy = [
  'jquery/dist/jquery.min.map',
  'font-awesome/fonts'
]

var prefix = function(b) {
  return './bower_components/' + b;
}

bowerjs = _.map(bowerjs, prefix);
bowercopy = _.map(bowercopy, prefix);

gulp.task('stylus', function() {
  gulp
    .src('./stylesheets/styles.styl')
    .pipe(plumber())
    .pipe(stylus({
      compress: true,
      use: nib(),
      sourcemap: {
        inline: true,
        sourceRoot: '.',
        basePath: 'public/build',
      }
    }))
    .pipe(sourcemaps.init({
      loadMaps: true
    }))
    .pipe(sourcemaps.write('.', {
      includeContent: false,
      sourceRoot: '.'
    }))
    .pipe(gulp.dest('./public/build'))
});

gulp.task('less', function() {
  gulp
    .src('./stylesheets/bootstrap.less')
    .pipe(plumber())
    .pipe(less())
    .pipe(minifyCSS())
    .pipe(gulp.dest('./public/build'))
});

gulp.task('html', function() {
  gulp
    .src('./app/**/*.jade')
    .pipe(plumber())
    .pipe(jade({
      doctype: 'html'
    }))
    .pipe(templateCache({
      filename: 'templates.js',
      standalone: true,
      base: function(file) {
        return path.basename(file.path);
      }
    }))
    .pipe(gulp.dest('./public/build'))
});

gulp.task('js', function() {
  gulp
    .src('./app/**/*.js')
    .pipe(plumber())
    .pipe(uglify('app.js', {
      mangle: false,
      outSourceMap: true
    }))
    .pipe(gulp.dest('./public/build'))
})

gulp.task('concat', function() {
  gulp
    .src(bowerjs)
    .pipe(concat('bower.js'))
    .pipe(gulp.dest('./public/build'))
})

gulp.task('copy', function() {
  _.each(bowercopy, function(f) {
    var dest = './public/build/';
    var src = f;

    if (fs.lstatSync(f).isDirectory()) {
      dest += _.last(f.split('/'));
      src += '/*';
    }

    gulp
      .src(src)
      .pipe(gulp.dest(dest))
  });
})

gulp.task('watch', function() {
  livereload.listen({
    port: 22746
  });

  var changed = function(file) {
    var ext = path.extname(file.path);
    switch(ext) {
      case '.js': gulp.start('js'); break;
      case '.jade': gulp.start('html'); break;
      case '.styl': gulp.start('stylus'); break;
      case '.less': gulp.start('less'); break;
    }
  }

  var paths = [
    './stylesheets/**/*.styl',
    './stylesheets/bootstrap.less',
    './views/*'
  ];

  paths.push('./app/**/*.js');
  paths.push('./app/**/*.jade');
  paths.push('./app/**/*.styl');

  _.each(paths, function(path) {
    watch(path, changed);
  });
});

gulp.task('default', ['concat', 'copy', 'stylus', 'less', 'html', 'js']);

