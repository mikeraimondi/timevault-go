'use strict'

###*
 # @ngdoc function
 # @name timevaultApp.controller:MainCtrl
 # @description
 # # MainCtrl
 # Controller of the timevaultApp
###
angular.module('timevaultApp')
  .controller 'MainCtrl', ($scope) ->
    $scope.awesomeThings = [
      'HTML5 Boilerplate'
      'AngularJS'
      'Karma'
    ]
