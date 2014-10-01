'use strict'

describe 'Controller: NotfoundCtrl', ->

  # load the controller's module
  beforeEach module 'timevaultApp'

  NotfoundCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    NotfoundCtrl = $controller 'NotfoundCtrl', {
      $scope: scope
    }

  it 'should attach a list of awesomeThings to the scope', ->
    expect(scope.awesomeThings.length).toBe 3
