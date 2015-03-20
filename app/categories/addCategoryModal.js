app.controller('addCategoryModal', function($scope, $modalInstance, $http, category) {

  $scope.category = category;

  $http.get('/api/categories?parents=1').success(function(data) {
    $scope.categories = data;
  })

  $scope.add = function() {
    $modalInstance.close($scope.selected.id);
  }

});