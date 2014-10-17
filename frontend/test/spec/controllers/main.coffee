'use strict'

describe 'Controller: MainCtrl', ->

  # load the controller's module
  beforeEach module 'timevaultApp'

  MainCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    MainCtrl = $controller 'MainCtrl', {
      $scope: scope
    }

  it 'should have the app name', ->
    expect(scope.appName).toBe 'Timevault'
