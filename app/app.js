app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {

  $locationProvider.html5Mode(false);

  var routes = {
    '/questions': 'listQuestions',
    '/questions/create': 'modifyQuestion',
    '/questions/view/:id': 'viewQuestion',
    '/questions/edit/:id': 'modifyQuestion',
    '/categories': 'listCategories',
    '/categories/create': 'modifyCategory',
    '/categories/edit/:id': 'modifyCategory'
  };

  for (var route in routes) {
    var controller = routes[route];
    $routeProvider.when(route, {
      templateUrl: controller + '.html',
      controller: controller
    });
  }

  $routeProvider.otherwise({
    redirectTo: '/questions'
  });

}]);