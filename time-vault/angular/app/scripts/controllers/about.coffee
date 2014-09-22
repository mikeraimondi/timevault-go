'use strict'

###*
 # @ngdoc function
 # @name timevaultApp.controller:AboutCtrl
 # @description
 # # AboutCtrl
 # Controller of the timevaultApp
###
angular.module('timevaultApp')
  .controller 'AboutCtrl', ($scope) ->
    $scope.awesomeThings = [
      'HTML5 Boilerplate'
      'AngularJS'
      'Karma'
    ]
