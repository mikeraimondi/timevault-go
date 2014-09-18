(function() {
  this.timevault.controller('HomeCtrl', [
    '$scope', function($scope) {
      return $scope.appName = 'Timevault';
    }
  ]);

  this.timevault.controller('PomodoroIndexCtrl', [
    '$scope', '$location', '$filter', '$interval', 'Pomodoro', function($scope, $location, $filter, $interval, Pomodoro) {
      $scope.pomodoros = [];
      $scope.init = function() {
        var timeoutId;
        this.pomodorosService = new Pomodoro;
        $scope.pomodoros = this.pomodorosService.all();
        return timeoutId = $interval(function() {
          return $scope.updateRemaining();
        }, 1000);
      };
      $scope.updateRemaining = function() {
        var pomodoro, _i, _len, _ref, _results;
        _ref = $scope.pomodoros;
        _results = [];
        for (_i = 0, _len = _ref.length; _i < _len; _i++) {
          pomodoro = _ref[_i];
          pomodoro.percentageLeft = this.pomodorosService.percentageLeft(pomodoro);
          pomodoro.remainingTime = this.pomodorosService.remainingTime(pomodoro);
          _results.push(pomodoro.progressBarType = this.pomodorosService.progressBarType(pomodoro));
        }
        return _results;
      };
      $scope.runningPomodoro = function() {
        var pomodoro, _i, _len, _ref;
        _ref = $scope.pomodoros;
        for (_i = 0, _len = _ref.length; _i < _len; _i++) {
          pomodoro = _ref[_i];
          if (pomodoro.percentageLeft > 0) {
            return pomodoro;
          }
        }
      };
      $scope.viewPomodoro = function(id) {
        return $location.url("/pomodoros/" + id);
      };
      $scope.destroyPomodoro = function(id) {
        this.pomodorosService.destroy(id);
        return $scope.pomodoros = $scope.pomodoros.filter(function(pomodoro) {
          return pomodoro.id !== id;
        });
      };
      $scope.addWorkPeriod = function() {
        var pomodoro;
        pomodoro = this.pomodorosService.create({
          set_duration: 1500,
          activity: 'work'
        });
        return $scope.pomodoros.unshift(pomodoro);
      };
      return $scope.addBreakPeriod = function() {
        var pomodoro;
        pomodoro = this.pomodorosService.create({
          set_duration: 300,
          activity: 'break'
        });
        return $scope.pomodoros.unshift(pomodoro);
      };
    }
  ]);

  this.timevault.controller('PomodoroShowCtrl', [
    '$scope', '$http', '$routeParams', 'Pomodoro', function($scope, $http, $routeParams, Pomodoro) {
      return $scope.init = function() {
        this.pomodorosService = new Pomodoro;
        return $scope.pomodoro = this.pomodorosService.find($routeParams.id);
      };
    }
  ]);

  this.timevault.controller('RegistrationCtrl', [
    '$scope', '$location', 'Auth', '$modalInstance', function($scope, $location, Auth, $modalInstance) {
      $scope.user = {};
      $scope.register = function() {
        var credentials;
        credentials = {
          email: $scope.user.email,
          password: $scope.user.password,
          password_confirmation: $scope.user.passwordConfirmation,
          phone_number: $scope.user.phoneNumber
        };
        return Auth.register(credentials).then(function(registeredUser) {
          return $modalInstance.close(registeredUser);
        }, function(response) {
          var errorMessage, errorType;
          return console.log(errorType, (function() {
            var _ref, _results;
            _ref = response.data.errors;
            _results = [];
            for (errorType in _ref) {
              errorMessage = _ref[errorType];
              _results.push(errorMessage);
            }
            return _results;
          })());
        });
      };
      return $scope.close = function() {
        return $modalInstance.close();
      };
    }
  ]);

  this.timevault.directive('pwCheck', [
    function() {
      return {
        require: 'ngModel',
        link: function(scope, element, attrs, ctrl) {
          var firstPassword;
          firstPassword = '#' + attrs.pwCheck;
          return element.add(firstPassword).on('keyup', function() {
            return scope.$apply(function() {
              var valid;
              valid = element.val() === $(firstPassword).val();
              return ctrl.$setValidity('pwmatch', valid);
            });
          });
        }
      };
    }
  ]);

  this.timevault.factory('Pomodoro', [
    '$resource', function($resource) {
      var Pomodoro;
      return Pomodoro = (function() {
        function Pomodoro() {
          this.service = $resource('/api/pomodoros/:id', {
            id: '@id'
          }, {
            update: {
              method: 'PATCH'
            }
          });
        }

        Pomodoro.prototype.create = function(attrs) {
          new this.service({
            pomodoro: attrs
          }).$save(function(pomodoro) {
            attrs.id = pomodoro.id;
            attrs.start = pomodoro.start;
            return attrs.activity = pomodoro.activity;
          });
          return attrs;
        };

        Pomodoro.prototype.all = function() {
          return this.service.query();
        };

        Pomodoro.prototype.find = function(id) {
          return this.service.get({
            id: id
          });
        };

        Pomodoro.prototype.destroy = function(id) {
          return this.service["delete"]({
            id: id
          });
        };

        Pomodoro.prototype.remainingSeconds = function(pomodoro) {
          var endSeconds, now, pomodoroEnd, pomodoroStart, remaining;
          pomodoroStart = new Date(pomodoro.start);
          endSeconds = pomodoroStart.getSeconds() + pomodoro.set_duration;
          pomodoroEnd = pomodoroStart.setSeconds(endSeconds);
          now = new Date();
          remaining = pomodoroEnd - now;
          if (remaining > 0) {
            return remaining / 1000;
          } else {
            return 0;
          }
        };

        Pomodoro.prototype.remainingTime = function(pomodoro) {
          var date, utc;
          date = new Date(null);
          date.setSeconds(this.remainingSeconds(pomodoro));
          utc = date.toUTCString();
          return utc.substr(utc.indexOf(':') - 2, 8);
        };

        Pomodoro.prototype.minutesLeft = function(pomodoro) {
          return this.remainingSeconds(pomodoro) / 60;
        };

        Pomodoro.prototype.percentageLeft = function(pomodoro) {
          return Math.floor((this.remainingSeconds(pomodoro) / pomodoro.set_duration) * 100);
        };

        Pomodoro.prototype.progressBarType = function(pomodoro) {
          var percent;
          percent = pomodoro.percentageLeft;
          switch (false) {
            case !(percent <= 10):
              return 'danger';
            case !(percent <= 30):
              return 'warning';
            default:
              return 'success';
          }
        };

        return Pomodoro;

      })();
    }
  ]);

}).call(this);

//# sourceMappingURL=app.js.map
