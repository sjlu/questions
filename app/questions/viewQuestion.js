app.controller('viewQuestion', function($scope, $http, $routeParams) {

  // $scope.comment = '';
  // $scope.commentInputDisabled = false;

  $http.get('/api/questions/' + $routeParams.id).success(function(data) {
    $scope.question = data;
  });

  // $scope.comments = [];

  // var getComments = function() {
  //   $http.get('/api/questions/' + $routeParams.id + '/comments').success(function(data) {
  //     $scope.comments = data;
  //   })
  // }
  // getComments();

  // var submitComment = function() {
  //   $scope.commentInputDisabled = true;
  //   $http.post('/api/questions/' + $routeParams.id + '/comments', {
  //     comment: $scope.comment
  //   }).success(function(data) {
  //     $scope.comment = '';
  //     $scope.commentInputDisabled = false;
  //     $scope.comments.push(data);
  //   });
  // }

  // $scope.detectSubmit = function(evt) {
  //   if (evt.keyCode == 13 && evt.shiftKey) {
  //     submitComment();
  //     evt.preventDefault();
  //   }
  // }

});