app.controller('listCategories', function($scope, $http, $modal) {
  var getCategories = function() {
    $http.get('/api/categories').success(function(categories) {
      $scope.categories = categories;
    })
  }
  getCategories();

  $scope.addCategoryTo = function(category) {
    var modal = $modal.open({
      templateUrl: 'addCategoryModal.html',
      controller: 'addCategoryModal',
      size: 'small',
      resolve: {
        category: function() {
          return category
        }
      }
    });
    modal.result.then(function(categoryId) {
      $http.post('/api/categories/' + category.id + '/subcategory', {
        category_id: categoryId
      }).success(function() {
        getCategories();
      });
    });
  }

  $scope.removeCategory = function(category, subcategory) {
    $http.delete('/api/categories/' + category.id + '/subcategory/' + subcategory.id)
      .success(function() {
        getCategories();
      })
  }

  $scope.deleteCategory = function(category) {
    $http.delete('/api/categories/' + category.id)
      .success(function() {
        getCategories();
      })
  }

});