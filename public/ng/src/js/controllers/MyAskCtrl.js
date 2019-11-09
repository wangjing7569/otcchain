app.controller('MyAskCtrl', ['$scope','$localStorage','$state','HttpService','REST_URL','$modal','DialogService','$q', function($scope, $localStorage, $state,HttpService,REST_URL,$modal,DialogService,$q) {

    $scope.rowCollectionPage = [];
    //  pagination
    $scope.itemsByPage=1;

    function render() {

        if($localStorage.loginuser) {

        }else {
            $state.go('access.signin');
        }

       HttpService.post(REST_URL.query, {fcn: "queryMyAsk", args:[$localStorage.loginuser.cmId]}).then(function (response) {
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
                      templateUrl: 'myAskInfo.html',
                      controller: 'MyAskModalInstanceCtrl',
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

var MyAskModalInstanceCtrl = function($scope, $modalInstance, items, $localStorage, REST_URL, HttpService, DialogService) {
    $scope.items = items;
    var ask = items.AskNum;
    console.log(ask);
    console.log($scope.items.ProvideID1);
    //console.log($localStorage.loginuser.cmId);
    var myDate = new Date();
    var s = myDate.getMonth().toString()+myDate.getDate().toString()+myDate.getHours().toString()+myDate.getMinutes().toString()+myDate.getSeconds().toString();

    $scope.item = {};
    $scope.item.OrderID = s;
    $scope.item.SaleParty = $scope.items.PublishID;
    $scope.item.BuyParty = $scope.items.ProvideID1;
    $scope.item.Underlying = $scope.items.Biaodi;
    $scope.item.ConstractType = $scope.items.ConstractType;
    $scope.item.StrikePrice = $scope.items.ExecuPrice;
    $scope.item.AccPrice = $scope.items.Price1;
    $scope.item.Angel = $scope.items.Angel;
    $scope.item.ConstractSize = $scope.items.Amount;
    $scope.item.ExpiringDate = $scope.items.EndDate;
    $scope.item.Request = '';

    console.log($scope.item);

    $scope.deal = function (){
        HttpService.post(REST_URL.invoke, { fcn: "deal", args:[ask, $scope.items.ProvideID1]}).then(function (response) {
              
                
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

        HttpService.post(REST_URL.invoke1, {fcn: "publish", args:[JSON.stringify($scope.item)]}).then(function (response) {
                //$scope.rows = JSON.parse(response.data.message);
                //$scope.loginuser = $localStorage.loginuser;
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

                if (!$scope.$$phase) {
                    $scope.$apply();
                }
                });

        //console.log($scope.rows);
    }
    $scope.cancel = function (){
        $modalInstance.dismiss('cancel');
    };
}
