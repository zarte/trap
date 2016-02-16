module.exports = function(grunt) {
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        uglify: {
            options: {
                banner: '/*! <%= pkg.name %> See README.md */\n'
            },
            dist: {
                files: {
                    'temp/html5.min.js': [
                        'node_modules/html5shiv/dist/html5shiv.min.js'
                    ],
                    'temp/script.min.js': [
                        'node_modules/chartist/dist/chartist.min.js',
                        'node_modules/jqu/dist/jqu.min.js',
                        'node_modules/jquery/dist/jquery.min.js',
                        'node_modules/vue/dist/vue.min.js',
                        'src/app.js',
                        'src/font.js',
                        'src/load.js'
                    ]
                }
            }
        },
        cssmin: {
            options: {
                shorthandCompacting: false,
                roundingPrecision: -1
            },
            dist: {
                files: {
                    'temp/style.min.css': [
                        'node_modules/normalize-css/normalize.css',
                        'node_modules/chartist/dist/chartist.min.css',
                        'src/common.css'
                    ]
                }
            }
        },
        copy: {
            dist: {
                src: 'src/index.html',
                dest: 'temp/index.html',
            }
        },
        processhtml: {
            dist: {
                files: {
                    'temp/index.packed.html': ['temp/index.html']
                }
            }
        },
        htmlmin: {
            dist: {
                options: {
                    removeComments: false,
                    collapseWhitespace: true
                },
                files: {
                    'temp/index.min.html': 'temp/index.packed.html'
                }
            }
        },
        compress: {
            dist: {
                options: {
                    mode: 'gzip'
                },
                files: [
                    {
                        src: ['temp/index.min.html'],
                        dest: 'temp/index.gz.html'
                    }
                ]
            }
        },
        base64: {
            dist: {
                files: {
                    'dist/index.b64': 'temp/index.gz.html'
                }
            }
        },
        clean: ["temp"]
    });

    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-htmlmin');
    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-contrib-compress');
    grunt.loadNpmTasks('grunt-processhtml');
    grunt.loadNpmTasks('grunt-base64');

    grunt.registerTask('default', [
        'uglify',
        'cssmin',
        'copy',
        'processhtml',
        'htmlmin',
        'compress',
        'base64',
        'clean'
    ]);
};