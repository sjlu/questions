app.controller('modifyQuestion', function($scope, $http, $location, $routeParams) {

  // modifyQuestion is misleading
  // this can either create a new one
  // or edit an exiting question

  $scope.question = {};
  $scope.questionCategories = [];

  if ($routeParams.id) {
    $http.get('/api/questions/' + $routeParams.id).success(function(data) {
      $scope.question = data;
      $scope.questionCategories = data.categories;
    });
  }

  $http.get('/api/categories').success(function(data) {
    $scope.categories = data;
  });

  $scope.save = function() {
    $scope.question.category_ids = _.pluck($scope.questionCategories, "id");

    var q;
    if ($routeParams.id) {
      q = $http.put('/api/questions/' + $routeParams.id, $scope.question);
    } else {
      q = $http.post('/api/questions', $scope.question)
    }
    q.success(function(data) {
      $location.path('/questions/view/' + data.id);
    });
  }

  $scope.addCategory = function() {

    if (!_.contains(_.pluck($scope.questionCategories, "id"), $scope.selectedCategory.id)) {
      $scope.questionCategories.push($scope.selectedCategory);
    }
    $scope.selectedCategory = null;

  }

  $scope.removeCategory = function(category) {

    $scope.questionCategories = _.reject($scope.questionCategories, function(c) {
      return c.id === category.id;
    });


  }


});