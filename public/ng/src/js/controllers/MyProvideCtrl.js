app.controller('MyProvideCtrl', ['$scope','$localStorage','$state','HttpService','REST_URL','$modal','DialogService','$q', function($scope, $localStorage, $state,HttpService,REST_URL,$modal,DialogService,$q) {

    $scope.rowCollectionPage = [];


    //  pagination
    $scope.itemsByPage=1;

    function render() {

        if($localStorage.loginuser) {

        }else {
            $state.go('access.signin');
        }
        console.log($localStorage.loginuser.cmId);
        HttpService.post(REST_URL.query, {fcn: "queryMyProvide", args:[$localStorage.loginuser.cmId]}).then(function (response) {
            $scope.rowCollectionPage = JSON.parse(response.data.message);
            $scope.loginuser = $localStorage.loginuser;

            if (!$scope.$$phase) {
                $scope.$apply();
            }
        });

    }
    render();
    $scope.open = function(row) {
                var modalInstance = $modal.open({  
                     templateUrl: 'myProvideInfo.html',
                     controller: 'MyProvideModalInstanceCtrl', 
                    resolve : {  
                        items : function() {  
                            $scope.items = row;
                            return $scope.items;  
                        }  
                    }  
                });  
                modalInstance.opened.then(function() {// 模态窗口打开之后执行的函数  
                    console.log('modal is opened');  
                });  
                
        };  
}]);

var MyProvideModalInstanceCtrl = function($scope, $modalInstance, items, $localStorage, REST_URL, HttpService, DialogService) {
    $scope.items = items;
    var ask = items.AskNum;
    console.log(ask);
    console.log($localStorage.loginuser.cmId);

    $scope.withdraw = function (){
        HttpService.post(REST_URL.invoke, { fcn: "withdraw", args:[ask, $localStorage.loginuser.cmId]}).then(function (response) {
                //alert("updateuser success");
                DialogService.open('infoDialog', {
                    scope: $scope,
                    title: '提交成功',
                    message: response.data.message,
                    onOk: function (value) {
                        //$state.go('app.table.myask');
                    },
                    onCancel: function (value) {
                        // do nothing
                    }
                });
            }, function (err) {
                DialogService.open('infoDialog', {
                    scope: $scope,
                    title: '提示',
                    message: '发生错误',
                    onOk: function (value) {
                    },
                    onCancel: function (value) {
                        // do nothing
                    }
                });
            });
    }
    $scope.cancel = function (){
        $modalInstance.dismiss('cancel');
    };
}
