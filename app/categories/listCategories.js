app.controller('listCategories', function($scope, $http, $modal) {
  var getCategories = function() {
    $http.get('/api/categories').success(function(categories) {
      $scope.categories = categories;
    })
  }
  getCategories();

  $scope.deleteCategory = function(category) {
    $http.delete('/api/categories/' + category.id)
      .success(function() {
        getCategories();
      })
  }

});