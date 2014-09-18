@timevault.directive 'pwCheck', [ ->
  require: 'ngModel',
  link: (scope, element, attrs, ctrl) ->
    firstPassword = '#' + attrs.pwCheck
    element.add(firstPassword).on 'keyup', ->
      scope.$apply ->
        valid = element.val() == $(firstPassword).val()
        ctrl.$setValidity('pwmatch', valid)
]
