var gulp = require('gulp');

var connect = require('gulp-connect');
var cp = require('child_process');
var deploy = require("gulp-gh-pages");
var gutil = require('gulp-util');
var minifyCSS = require('gulp-minify-css');
var plumber = require('gulp-plumber');
var prefixer = require('gulp-autoprefixer');
var rename = require('gulp-rename');
var sass = require('gulp-sass');
var postcss = require('gulp-postcss');
var customMedia = require("postcss-custom-media");
var imageMin = require('gulp-imagemin');
var gm = require('gulp-gm');
var gulpif = require('gulp-if');
var args   = require('yargs').argv;

// Path Vars
var BASE_PATH = './';
var DEV = BASE_PATH + '_dev/';
var SRC = DEV + 'src/';
var SITE_PATH = BASE_PATH + '_site/'
var DIST = BASE_PATH + 'assets/';

var ASSETS = {
    js: [
        SRC + '/js/**/*.js',
    ],
    css: [
        SRC + '/sass/**/*.scss'
    ],
    img: [
        SRC + '/img/**/*'
    ],
    jekyll: [
        'index.html',
        '_posts/*',
        '_layouts/*',
        '_includes/*' ,
        '_data/*.yml'
    ]
};

var isProd = args.env === 'prod';

// BUILD TASKS

gulp.task('css', function() {

    var processors = [
        customMedia
    ];

    return gulp.src(ASSETS.css)
        .pipe(plumber(function(error) {
            gutil.log(gutil.colors.red(error.message));
            gulp.task('css').emit('end');
        }))
        .pipe(sass())
        .pipe(prefixer('last 3 versions', 'ie 9'))
        .pipe(minifyCSS())
        .pipe(postcss(processors))
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

gulp.task('img', function() {
     return gulp.src(ASSETS.img)
        .pipe(plumber(function(error) {
            gutil.log(gutil.colors.red(error.message));
            gulp.task('img').emit('end');
        }))

        .pipe(gulpif(isProd, gm(function (gmfile) {
          return gmfile
            .resize(750)
            .interlace('Line')
            .quality(70)
        })))
            // imageMin({
            //     progressive: true,
            // })



        .pipe(rename({dirname: DIST + '/img'}))
        .pipe(gulp.dest('./'));
});

gulp.task('jekyll', ['css', 'js', 'img'], function(code) {
    return cp.spawn('jekyll', ['build', '--incremental'], {stdio: 'inherit'})
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
    gulp.watch(ASSETS.img, ['img']);
    gulp.watch(ASSETS.jekyll, ['jekyll']);
});

gulp.task('deploy', ['jekyll'], function () {
    return gulp.src('./_site/**/*')
        .pipe(deploy({
            force: true
        }));
});


gulp.task('default', ['jekyll', 'server', 'watch']);