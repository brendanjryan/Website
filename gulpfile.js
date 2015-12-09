var gulp = require('gulp');

var sass = require('gulp-sass');
var gutil = require('gulp-util');
var plumber = require('gulp-plumber');
var rename = require('gulp-rename');
var minifyCSS = require('gulp-minify-css');
var prefixer = require('gulp-autoprefixer');
var connect = require('gulp-connect');
var cp = require('child_process');

// Path Vars
var BASE_PATH = './';
var DEV = BASE_PATH + '_dev/';
var SRC = DEV + 'src/';
var SITE_PATH = BASE_PATH + '_site/'
var DIST = SITE_PATH + 'assets/';

var ASSETS = {
    js: [
        SRC + '/js/**/*.js',
    ],
    css: [
        SRC + '/sass/**/*.scss'
    ],
    jekyll: [
        'index.html',
        '_posts/*',
        '_layouts/*',
        '_includes/*' ,
        'assets/**/*'
    ]
};


// BUILD TASKS

gulp.task('css', function() {
    return gulp.src(ASSETS.css)
        .pipe(plumber(function(error) {
            gutil.log(gutil.colors.red(error.message));
            gulp.task('css').emit('end');
        }))
        .pipe(sass())
        .pipe(prefixer('last 3 versions', 'ie 9'))
        .pipe(minifyCSS())
        .pipe(rename({
            dirname: DIST + '/css',
            basename: 'styles'
        }))
        .pipe(gulp.dest('./'));
});

gulp.task('js', function() {
     return gulp.src(ASSETS.js)
        .pipe(plumber(function(error) {
            gutil.log(gutil.colors.red(error.message));
            gulp.task('js').emit('end');
        }))

        .pipe(rename({dirname: DIST + '/js'}))
        .pipe(gulp.dest('./'));
});


gulp.task('jekyll', ['css', 'js'], function(code) {
    return cp.spawn('jekyll', ['build'], {stdio: 'inherit'})
        .on('error', function(error) {
            gutil.log(gutil.colors.red(error.message))
        })
        .on('close', code);
});

gulp.task('server', function() {
    connect.server({
        root: ['_site'],
        port: 4000
    });
});

gulp.task('watch', function() {
    gulp.watch(ASSETS.css, ['css']);
    gulp.watch(ASSETS.js, ['js']);
    gulp.watch(ASSETS.jekyll, ['jekyll']);
});



gulp.task('default', ['jekyll', 'server', 'watch']);