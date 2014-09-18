module.exports = function(grunt) {
  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),

    // sass: {
    //   options: {
    //     includePaths: ['bower_components/foundation/scss']
    //   },
    //   dist: {
    //     options: {
    //       outputStyle: 'compressed'
    //     },
    //     files: {
    //       'css/app.css': 'src/scss/app.scss'
    //     }        
    //   }
    // },

    coffee: {
      dist: {
        options: {
          sourceMap: true
        },
        files: {
          'time-vault/assets/js/app.js': ['time-vault/angular/**/*.coffee']
        }
      }
    },

    uglify: {
      dist: {
        options: {
          sourceMap: true,
          sourceMapIn: 'time-vault/assets/js/app.js.map'
        },
        files: {
          'time-vault/assets/js/app.min.js': ['time-vault/assets/js/app.js']
        }
      }
    },

    copy: {
      dist: {
        files: {
          'time-vault/assets/js/angular.min.js': 'bower_components/angular/angular.min.js',
          'time-vault/assets/js/angular.min.js.map': 'bower_components/angular/angular.min.js.map'
        }
      }
    },

    watch: {
      grunt: { files: ['Gruntfile.js'] },

      // sass: {
      //   files: 'src/scss/**/*.scss',
      //   tasks: ['sass']
      // },

      uglify: {
        files: 'src/js/**/*.js',
        tasks: ['uglify']
      }
    }
  });

  // grunt.loadNpmTasks('grunt-sass');
  // grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-coffee');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-copy');
  // grunt.registerTask('build', ['sass', 'uglify']);
  grunt.registerTask('build', ['coffee', 'copy', 'uglify']);
  grunt.registerTask('default', ['build', 'watch']);
};
