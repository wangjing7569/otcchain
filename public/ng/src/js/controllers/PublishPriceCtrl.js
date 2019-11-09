
app.controller('PublishPriceCtrl', ['$scope', '$rootScope', '$http', '$modal', '$log', 'REST_URL', 'HttpService', 'DialogService', '$localStorage', '$state', '$stateParams',
    function ($scope, $rootScope, $http, $modal, $log, REST_URL, HttpService, DialogService, $localStorage, $state, $stateParams) {
        var myDate = new Date();
        var s = myDate.getMonth().toString()+myDate.getDate().toString()+myDate.getHours().toString()+myDate.getMinutes().toString()+myDate.getSeconds().toString();

        $scope.item = {};
        $scope.item.AskNum = s;
        $scope.item.Biaodi = '';
        //$scope.item.EnsurePrice = '';


        $scope.item.CurrentPrice = '';

        $scope.item.EndDate = '';
        $scope.item.Angel = '';
        $scope.item.ConstractType = '';
        $scope.item.ExecuPrice = '';
        $scope.item.Amount = '';
        //console.log($localStorage.loginuser);
        $scope.item.PublishID = $localStorage.loginuser.cmId;
        $scope.item.HodrAcct = $localStorage.loginuser.Acct;

        /*$scope.dataController = function($http,$scope){
            $http.get("json/config.json").success(function(freetrial) { alert(freetrial);$scope.data = freetrial;});

        };

        $scope.dataController($http,$scope);

        console.log($scope.data);*/

       // function dataController($http,$scope) {
       
 

        $scope.save = function () {
            //alert($scope.item.AskNum);
            var str;
            if ($scope.item.Biaodi == "上证50ETF"){
                str = '3.003';
            }
            if ($scope.item.Biaodi == "豆粕1912"){
                str='2994';
            }
            if ($scope.item.Biaodi == "豆粕2001"){
                str= '2956';
            }
            if($scope.item.Biaodi == "豆粕2003"){
                str= '2900';
            }
            if($scope.item.Biaodi == "五花肉"){
                str='';
            }
            $scope.item.CurrentPrice = str;

            if ($scope.item.Biaodi == "上证50ETF"){
                if ($scope.item.ConstractType == "美式看涨" || $scope.item.ConstractType == "美式看跌"){
                    alert('该标的期权只能是欧式');
                }
                else{
                    HttpService.post(REST_URL.invoke, { fcn: "publish", args:[JSON.stringify($scope.item)]}).then(function (response) {
                //alert("updateuser success");
                DialogService.open('infoDialog', {
                    scope: $scope,
                    title: '提交成功',
                    message: response.data.message,
                    onOk: function (value) {
                        $state.go('app.table.myask');
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
            }else if ($scope.item.Biaodi == "豆粕1912" || $scope.item.Biaodi == "豆粕2001" || $scope.item.Biaodi == "豆粕2003"){
                if ($scope.item.ConstractType == "欧式看涨" || $scope.item.ConstractType == "欧式看跌"){
                    alert('该标的期权只能是美式');
                }
                else{
                    HttpService.post(REST_URL.invoke, { fcn: "publish", args:[JSON.stringify($scope.item)]}).then(function (response) {
                //alert("updateuser success");
                DialogService.open('infoDialog', {
                    scope: $scope,
                    title: '提交成功',
                    message: response.data.message,
                    onOk: function (value) {
                        $state.go('app.table.myask');
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
            }             
        };

        $scope.cancel = function () {
            console.log($state);
            $state.go('app.table.myask');
        };





    }]);