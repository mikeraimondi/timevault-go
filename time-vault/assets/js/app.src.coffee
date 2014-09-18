@timevault.controller 'HomeCtrl', ['$scope', ($scope) ->
  $scope.appName = 'Timevault'
]

@timevault.controller 'PomodoroIndexCtrl', [
  '$scope', '$location', '$filter', '$interval', 'Pomodoro',
  ($scope, $location, $filter, $interval, Pomodoro) ->
    $scope.pomodoros = []

    $scope.init = ->
      @pomodorosService = new Pomodoro
      $scope.pomodoros = @pomodorosService.all()

      timeoutId = $interval ->
        $scope.updateRemaining()
      , 1000

    $scope.updateRemaining = ->
      for pomodoro in $scope.pomodoros
        pomodoro.percentageLeft = @pomodorosService.percentageLeft(pomodoro)
        pomodoro.remainingTime = @pomodorosService.remainingTime(pomodoro)
        pomodoro.progressBarType = @pomodorosService.progressBarType(pomodoro)

    $scope.runningPomodoro = ->
      for pomodoro in $scope.pomodoros
        return pomodoro if pomodoro.percentageLeft > 0

    $scope.viewPomodoro = (id) ->
      $location.url "/pomodoros/#{id}"

    $scope.destroyPomodoro = (id) ->
      @pomodorosService.destroy id
      $scope.pomodoros = $scope.pomodoros.filter (pomodoro) ->
        pomodoro.id isnt id

    $scope.addWorkPeriod = ->
      pomodoro = @pomodorosService.create
        set_duration: 1500
        activity: 'work'
      $scope.pomodoros.unshift(pomodoro)

    $scope.addBreakPeriod = ->
      pomodoro = @pomodorosService.create
        set_duration: 300
        activity: 'break'
      $scope.pomodoros.unshift(pomodoro)
]

@timevault.controller 'PomodoroShowCtrl', [
  '$scope', '$http', '$routeParams', 'Pomodoro',
  ($scope, $http, $routeParams, Pomodoro) ->

    $scope.init = ->
      @pomodorosService = new Pomodoro
      $scope.pomodoro = @pomodorosService.find $routeParams.id
]

@timevault.controller 'RegistrationCtrl', [
  '$scope', '$location', 'Auth', '$modalInstance',
  ($scope, $location, Auth, $modalInstance) ->

    $scope.user = {}

    $scope.register = ->
      credentials =
        email: $scope.user.email
        password: $scope.user.password
        password_confirmation: $scope.user.passwordConfirmation
        phone_number: $scope.user.phoneNumber

      Auth.register(credentials).then(
        (registeredUser) ->
          $modalInstance.close registeredUser
        (response) ->
          console.log(
            errorType,
            errorMessage for errorType, errorMessage of response.data.errors))

    $scope.close = ->
      $modalInstance.close()

]

@timevault.directive 'pwCheck', [ ->
  require: 'ngModel',
  link: (scope, element, attrs, ctrl) ->
    firstPassword = '#' + attrs.pwCheck
    element.add(firstPassword).on 'keyup', ->
      scope.$apply ->
        valid = element.val() == $(firstPassword).val()
        ctrl.$setValidity('pwmatch', valid)
]

@timevault.factory 'Pomodoro', ['$resource', ($resource) ->
  class Pomodoro
    constructor: ->
      @service = $resource('/api/pomodoros/:id',
        {id: '@id'},
        {update: {method: 'PATCH'}})

    create: (attrs) ->
      new @service(pomodoro: attrs).$save (pomodoro) ->
        attrs.id = pomodoro.id
        attrs.start = pomodoro.start
        attrs.activity = pomodoro.activity
      attrs

    all: ->
      @service.query()

    find: (id) ->
      @service.get id: id

    destroy: (id) ->
      @service.delete id: id

    remainingSeconds: (pomodoro) ->
      pomodoroStart = new Date(pomodoro.start)
      endSeconds = pomodoroStart.getSeconds() + pomodoro.set_duration
      pomodoroEnd = pomodoroStart.setSeconds(endSeconds)
      now = new Date()
      remaining = pomodoroEnd - now
      if remaining > 0
        remaining / 1000
      else
        0

    remainingTime: (pomodoro) ->
      date = new Date(null)
      date.setSeconds @remainingSeconds(pomodoro)
      utc = date.toUTCString()
      utc.substr(utc.indexOf(':') - 2, 8)

    minutesLeft: (pomodoro) ->
      @remainingSeconds(pomodoro) / 60
      
    percentageLeft: (pomodoro) ->
      Math.floor (@remainingSeconds(pomodoro) / pomodoro.set_duration) * 100

    progressBarType: (pomodoro) ->
      percent = pomodoro.percentageLeft
      switch
        when percent <= 10 then 'danger'
        when percent <= 30 then 'warning'
        else 'success'
]
