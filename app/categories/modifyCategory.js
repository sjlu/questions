app.controller('modifyCategory', function($scope, $http, $location, $routeParams) {

  $scope.category = {};
  $scope.parentCategories = [];

  $http.get('/api/categories').success(function(data) {
    $scope.categories = data;
  });

  if ($routeParams.id) {
    $http.get('/api/categories/' + $routeParams.id).success(function(data) {
      $scope.category = data;
      $scope.parentCategories = data.parent_categories;
    });
  }

  $scope.save = function() {
    $scope.category.parent_ids = _.pluck($scope.parentCategories, "id");

    var q;
    if ($routeParams.id) {
      q = $http.put('/api/categories/' + $routeParams.id, $scope.category);
    } else {
      q = $http.post('/api/categories', $scope.category)
    }

    q.success(function() {
      $location.path('/categories');
    });
  }

});