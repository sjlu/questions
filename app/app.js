app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {

  $locationProvider.html5Mode(false);

  var routes = {
    '/': 'questions/list.html'
  };

  for (var route in routes) {
    var template = routes[route];
    $routeProvider.when(route, {
      templateUrl: template
    });
  }

  $routeProvider.otherwise({
    redirectTo: '/'
  });

}]);