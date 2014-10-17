'use strict'

###*
 # @ngdoc service
 # @name timevaultApp.pomodoro
 # @description
 # # pomodoro
 # Factory in the timevaultApp.
###
angular.module('timevaultApp')
  .factory 'Pomodoro', ['$resource', 'config', ($resource, config) ->
    class Pomodoro
      constructor: ->
        @service = $resource(config.url,
          {id: '@id'},
          {update: {method: 'PATCH'}})

      create: (attrs) ->
        new @service(pomodoro: attrs).$save (pomodoro) ->
          attrs.id = pomodoro.id
          attrs.createdAt = pomodoro.createdAt
          attrs.activity = pomodoro.activity
        attrs

      all: ->
        @service.query()

      find: (id) ->
        @service.get id: id

      destroy: (id) ->
        @service.delete id: id

      remainingSeconds: (pomodoro) ->
        pomodoroStart = new Date(pomodoro.createdAt)
        endSeconds = pomodoroStart.getSeconds() + (pomodoro.duration / 1000000000)
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
        Math.floor (@remainingSeconds(pomodoro) / (pomodoro.duration / 1000000000)) * 100

      progressBarType: (pomodoro) ->
        percent = pomodoro.percentageLeft
        switch
          when percent <= 10 then 'danger'
          when percent <= 30 then 'warning'
          else 'success'
  ]
