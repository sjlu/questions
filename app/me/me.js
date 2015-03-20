app.controller('me', function($scope, $http) {
  $http.get('/api/me').success(function(data) {
    $scope.user = data;
  });

  $scope.save = function() {
    $http.put('/api/me', $scope.user).success(function(data) {
      $scope.user = data;
    });
  }
});