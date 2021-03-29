const gulp = require('gulp');
const ts = require('gulp-typescript');
const sourcemaps = require('gulp-sourcemaps');
const path = require('path');
const nodemon = require('gulp-nodemon');

const tsProject = ts.createProject('tsconfig.json');
const JS_FILES = ['src/*.js', 'src/**/*.js'];
gulp.task('assets', done => {
  gulp.src(JS_FILES).pipe(gulp.dest('dist'));
  done();
});

gulp.task('default', done => {
  gulp.series('assets')();
  const tsResult = tsProject
    .src()
    .pipe(sourcemaps.init())
    .pipe(tsProject())
    .on('error', err => console.log(err));
  tsResult.js
    .pipe(
      sourcemaps.write({
        // Return relative source map root directories per file.
        sourceRoot(file) {
          const sourceFile = path.join(file.cwd, file.sourceMap.file);
          return path.relative(path.dirname(sourceFile), file.cwd);
        }
      })
    )
    .pipe(gulp.dest('dist'));
  done();
});

gulp.task('start', done => {
  const node = nodemon({
    script: 'dist/',
    watch: 'src/',
    execMap: {
      js: 'node --inpect'
    },
    done
  })
    .on('start', () => {})
    .on('restart', ['default']);

  gulp.watch('src/**/*.ts', () => {
    console.log('Detected changes, restarting');
    node.emit('restart');
  });
});
