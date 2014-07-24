var gulp = require('gulp');

var bower          = require('gulp-bower');
var csslint        = require('gulp-csslint');
var cssmin         = require('gulp-cssmin');
var mainBowerFiles = require('main-bower-files');
var plumber        = require('gulp-plumber');
var rename         = require('gulp-rename');
var runSequence    = require('run-sequence');

var path = {
  css: ['public/css/*.css', '!public/css/*.min.css'],
  lib: 'public/lib'
};

gulp.task('default', function () {
  gulp.watch(path.css, ['cssmin']);
});

gulp.task('cssmin', function () {
  gulp.src(path.css)
  .pipe(plumber())
  .pipe(cssmin())
  .pipe(rename({suffix: '.min'}))
  .pipe(gulp.dest('public/css'));
});

gulp.task('csslint', function () {
  gulp.src(path.css)
  .pipe(csslint())
  .pipe(csslint.reporter());
});

gulp.task('bower-install', function () {
  bower();
});

gulp.task('bower-copy', function () {
  gulp.src(mainBowerFiles())
  .pipe(gulp.dest(path.lib));
});

gulp.task('bower', function () {
   runSequence('bower-install', 'bower-copy');
});
