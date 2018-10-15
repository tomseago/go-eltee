var gulp = require('gulp');
var pug = require('gulp-pug');
var stylus = require('gulp-stylus');
var minifyCSS = require('gulp-csso');
var concat = require('gulp-concat');
var sourcemaps = require('gulp-sourcemaps');
//var watch = require("gulp-watch");

var DEST_BASE = "./root";

gulp.task('lib', function(){
  return gulp.src("lib/*")
    .pipe(gulp.dest(DEST_BASE+"/js"))
});

gulp.task('html', function(){
  return gulp.src('html/*.pug')
    .pipe(pug({
        verbose: true
    }))
    .pipe(gulp.dest(DEST_BASE))
});

gulp.task('css', function(){
  // return gulp.src('css/*.less')
  //   .pipe(less())
  //   .pipe(minifyCSS())
  //   .pipe(gulp.dest('build/css'))
  return gulp.src('css/*.styl')
    .pipe(stylus({
      'include css': true
    }))
    .pipe(minifyCSS())
    .pipe(gulp.dest(DEST_BASE)); 
});

gulp.task('js', function(){
  return gulp.src('javascript/*.js')
    .pipe(sourcemaps.init())
    .pipe(concat('app.min.js'))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(DEST_BASE+"/js"))
});

gulp.task('default', [ 'lib', 'html', 'css', 'js' ]);
          