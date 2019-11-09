app.controller('HalfCtrl', ['$scope','$localStorage','$state','HttpService','REST_URL','$modal','DialogService','$q', function($scope, $localStorage, $state,HttpService,REST_URL,$modal,DialogService,$q) {
    $scope.rowCollectionPage = [];
    $scope.itemsByPage=1;
    $scope.itemsxq = '';
    function render() {

        if($localStorage.loginuser) {

        }else {
            $state.go('access.signin');
        }
       HttpService.post(REST_URL.query1, {fcn: "queryHalf1", args:[$localStorage.loginuser.cmId]}).then(function (response) {
                $scope.rowCollectionPage = JSON.parse(response.data.message);
            
                $scope.loginuser = $localStorage.loginuser;

                if (!$scope.$$phase) {
                    $scope.$apply();
                }
            });

        HttpService.post(REST_URL.query1, {fcn: "queryHalf2", args:[$localStorage.loginuser.cmId]}).then(function (response) {
                $scope.rowCollections = JSON.parse(response.data.message);
                for (var i = 0; i< $scope.rowCollections.length; i++){
                    if ($scope.rowCollections[i].Angel == '买入'){
                        $scope.rowCollections[i].Angel = '卖出';
                    }else if ($scope.rowCollections[i].Angel == '卖出'){
                        $scope.rowCollections[i].Angel = '买入';
                    }
                }
                //console.log($scope.rowCollections[0].Angel);
                if (!$scope.$$phase) {
                    $scope.$apply();
                }
            });
    }
    render();

}]);