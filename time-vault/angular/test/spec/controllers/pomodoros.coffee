'use strict'

describe 'Controller: PomodorosCtrl', ->

  # load the controller's module
  beforeEach module 'timevaultApp'

  PomodorosCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    PomodorosCtrl = $controller 'PomodorosCtrl', {
      $scope: scope
    }

  # it 'should attach a list of awesomeThings to the scope', ->
  #   expect(scope.awesomeThings.length).toBe 3
