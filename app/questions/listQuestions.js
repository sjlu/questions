app.controller('listQuestions', function($scope, $http) {
  $http.get('/api/questions').success(function(questions) {
    $scope.questions = questions;
  });
});