/*jshint camelcase: false */
/*global module:false */
const url = require('url');
module.exports = function(grunt) {
  grunt.initConfig({
    exec: {
      'build_binary': {
        command: 'GOOS=linux GOARCH=amd64 go build'
      },
      'deploy': {
        command: 'scp -r -i ~/.ssh/tweetify.pem ./dist/templates ./dist/assets ./showmyprs.com ubuntu@52.87.181.233:/home/ubuntu/showmyprs'
      }
    },

    useminPrepare: {
      options: {
        dest: 'dist'
      },
      html: [
        'dist/**/*.html',
      ]
    },

    // Performs rewrites based on rev and the useminPrepare configuration
    usemin: {
      options: {
        assetsDirs: [
          'dist',
          'dist/assets/fonts'
        ]
      },
      html: ['dist/**/*.html'],
      css: ['dist/**/*.css']
    },

    // Renames files for browser caching purposes
    filerev: {
      dist: {
        src: [
          'dist/assets/**/*.*',
           /* exclude favicons from fingerprinting */
          '!dist/assets/favicon.ico',
          '!dist/assets/favicon.png'
        ]
      }
    },

    // Empties folders to start fresh
    clean: {
      dist: {
        files: [{
          dot: true,
          src: [
            'dist',
            '.tmp'
          ]
        }]
      }
    },

    copy: {
      static_dist: {
        files: [{
          expand: true,
          src: [
            'assets/**/*',
            'templates/**/*'
          ],
          dest: 'dist'
        }]
      },
    },

    postcss: {
      options: {
        map: true,
        processors: [
          // Add vendor prefixed styles
          require('autoprefixer')({
            browsers: [
              'Android 2.3',
              'Android >= 4',
              'Chrome >= 20',
              'Firefox >= 24',
              'Explorer >= 8',
              'iOS >= 6',
              'Opera >= 12',
              'Safari >= 6'
            ]
          })
        ]
      },
      dist: {
        files: [{
          expand: true,
          cwd: 'dist/assets/styles/',
          src: '{,*/}*.css',
          dest: 'dist/assets/styles/'
        }]
      }
    }
  });

  grunt.loadNpmTasks('grunt-exec');
  grunt.loadNpmTasks('grunt-contrib-copy');
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-cssmin');
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-usemin');
  grunt.loadNpmTasks('grunt-filerev');
  grunt.loadNpmTasks('grunt-postcss');

  grunt.registerTask('build_static', [
    'copy:static_dist',
    'useminPrepare',
    'postcss',
    'concat',
    'cssmin',
    'filerev',
    'usemin',
  ]);

  grunt.registerTask('build', [
    'clean',
    'build_static',
    'exec:build_binary'
  ]);

  grunt.registerTask('deploy', [
    'exec:deploy'
  ]);

};

