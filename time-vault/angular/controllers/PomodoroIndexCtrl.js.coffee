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
