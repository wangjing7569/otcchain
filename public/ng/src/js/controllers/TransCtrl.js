app.controller('TransCtrl', ['$scope','$localStorage','$state','HttpService','REST_URL','$modal','DialogService','$q', function($scope, $localStorage, $state,HttpService,REST_URL,$modal,DialogService,$q) {
    $scope.rowCollectionPage = [];
    $scope.itemsByPage=1;
    $scope.itemsxq = '';
    function render() {

        if($localStorage.loginuser) {

        }else {
            $state.go('access.signin');
        }
       HttpService.post(REST_URL.query1, {fcn: "queryMySale", args:[$localStorage.loginuser.cmId]}).then(function (response) {
                $scope.rowCollectionPage = JSON.parse(response.data.message);
            
                $scope.loginuser = $localStorage.loginuser;

                if (!$scope.$$phase) {
                    $scope.$apply();
                }
            });

        HttpService.post(REST_URL.query1, {fcn: "queryMyBuy", args:[$localStorage.loginuser.cmId]}).then(function (response) {
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
    $scope.open = function(row) {
                var modalInstance = $modal.open({  
                      templateUrl: 'myTransInfo.html',
                      controller: 'MyOrderModalInstanceCtrl',
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
                modalInstance.result.then(function(result){
                    $scope.items = result;
                    console.log($scope.items.flag);
                    if ($scope.items.flag== 'xq'){ 
                        $scope.openxq(result);
                    }
                    if ($scope.items.flag == 'end'){
                        $scope.openend(result);
                    }
                    if ($scope.items.flag == 'p1'){
                        $scope.ping1(result);
                    }
                    if ($scope.items.flag == 'p2'){
                        $scope.ping2(result);
                    }
                    if ($scope.items.flag == 'p3'){
                        $scope.ping3(result);
                    }
                    
                    console.log(result);

                }, function(reason){
                    console.log(reason);
                });     
    };

    $scope.openxq = function(itemsxq){
        var modalInstance = $modal.open({  
                      templateUrl: 'myXqInfo.html',
                      controller: 'MyXqModalInstanceCtrl',
                    resolve : {  
                        items : function() {  
                            $scope.items = itemsxq;
                            return $scope.items;  
                        }  
                    }  
                });  
                modalInstance.opened.then(function() {// 模态窗口打开之后执行的函数  
                    console.log('modal is opened');  
        });     
    };

     $scope.openend = function(itemsend){
        var modalInstance = $modal.open({  
                      templateUrl: 'myEndInfo.html',
                      controller: 'MyEndModalInstanceCtrl',
                    resolve : {  
                        items : function() {  
                            $scope.items = itemsend;
                            return $scope.items;  
                        }  
                    }  
                });  
                modalInstance.opened.then(function() {// 模态窗口打开之后执行的函数  
                    console.log('modal is opened');  
        });     
    };

    $scope.ping1 = function(itemsping){
        var modalInstance = $modal.open({  
                      templateUrl: 'myPing1Info.html',
                      controller: 'MyPing1ModalInstanceCtrl',
                    resolve : {  
                        items : function() {  
                            $scope.items = itemsping;
                            return $scope.items;  
                        }  
                    }  
                });  
                modalInstance.opened.then(function() {// 模态窗口打开之后执行的函数  
                    console.log('modal is opened');  
        });     
    };


    $scope.ping2 = function(itemsping){
        var modalInstance = $modal.open({  
                      templateUrl: 'myPing2Info.html',
                      controller: 'MyPing2ModalInstanceCtrl',
                    resolve : {  
                        items : function() {  
                            $scope.items = itemsping;
                            return $scope.items;  
                        }  
                    }  
                });  
                modalInstance.opened.then(function() {// 模态窗口打开之后执行的函数  
                    console.log('modal is opened');  
        });     
    };


    $scope.ping3 = function(itemsping){
        var modalInstance = $modal.open({  
                      templateUrl: 'myPing3Info.html',
                      controller: 'MyPing3ModalInstanceCtrl',
                    resolve : {  
                        items : function() {  
                            $scope.items = itemsping;
                            return $scope.items;  
                        }  
                    }  
                });  
                modalInstance.opened.then(function() {// 模态窗口打开之后执行的函数  
                    console.log('modal is opened');  
        });     
    };


}]);


var MyXqModalInstanceCtrl = function($scope, $modalInstance, items, $localStorage, REST_URL, HttpService, DialogService) {
    $scope.items = items;
    //console.log(items);
    console.log($scope.items.status);
    //按照随机数取结算结果，0~7为结算成功，8为结算10%，9为结算30%， 10为结算50%
    $scope.random = function(){
        let a = Math.ceil(Math.random()*10);
        if (a<=7){
            $scope.items.random = 1;
        }else if(a==8){
            $scope.items.random = 0.1;
        }else if(a==9){
            $scope.items.random = 0.3;
        }else{
            $scope.items.random = 0.5;
        }
    }
    $scope.random();


    var r = $scope.items.random;
    console.log(r);
    console.log($scope.items.Percent);
    console.log($scope.items.Percent == "");
    if ($scope.items.Percent == ""){
        $scope.items.Percent = r + "";
    }else{
        $scope.items.Percent = (r + parseFloat($scope.items.Percent)).toFixed(2);
    }
    console.log($scope.items.Percent);


    $scope.submit = function(){
        if($scope.items.status == '可行权'){   
            $scope.items.Percent = '0.3';        

            if ($scope.items.Percent >= 1){
                HttpService.post(REST_URL.invoke1, { fcn: "expiration", args:[$scope.items.OrderID]}).then(function (response) {
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
             }else{
                $scope.items.Request = '结算交收';
                HttpService.post(REST_URL.invoke1, { fcn: "expirationHalf", args:[$scope.items.OrderID, $scope.items.Percent, $scope.items.Request]}).then(function (response) {
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
        }else{
            $scope.items.status = '不可行权';
            alert('当前交易不可行权');
        }

    };


    $scope.cancel = function (){
        $modalInstance.dismiss('cancel');
    };
}


var MyPing1ModalInstanceCtrl = function($scope, $modalInstance, items, $localStorage, REST_URL, HttpService, DialogService) {
    $scope.items = items;
    console.log(items);

    console.log($scope.items.status);
    $scope.submit = function(){
        //if($scope.items.earnper < 0){
            //$scope.items.Request = '请求平仓';
            HttpService.post(REST_URL.invoke1, { fcn: "askForOffset", args:[$scope.items.OrderID]}).then(function (response) {
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

       // }else{
        //    alert('当前处于盈利状态，不允许执行此操作');
        //}
    }

    $scope.cancel = function (){
        $modalInstance.dismiss('cancel');
    };
}


var MyPing2ModalInstanceCtrl = function($scope, $modalInstance, items, $localStorage, REST_URL, HttpService, DialogService) {
    $scope.items = items;
    console.log(items);

    console.log($scope.items.status);

    $scope.getPrice = function (){
        $scope.items.baojia = '50.0';
        $scope.items.baojiaPer = '50.0%';
    }


    $scope.submit = function(){
        HttpService.post(REST_URL.invoke1, { fcn: "bid", args:[$scope.items.OrderID, $scope.items.baojia]}).then(function (response) {
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
    };

    $scope.cancel = function (){
        $modalInstance.dismiss('cancel');
    };
}



var MyPing3ModalInstanceCtrl = function($scope, $modalInstance, items, $localStorage, REST_URL, HttpService, DialogService) {
    $scope.items = items;
    console.log(items);

    console.log($scope.items.status);

    //获取随机比例
    $scope.random = function () {
        let a = Math.ceil(Math.random() * 10);
        if (a <= 7) {
            $scope.items.random = 1;
        } else if (a == 8) {
            $scope.items.random = 0.1;
        } else if (a == 9) {
            $scope.items.random = 0.3;
        } else {
            $scope.items.random = 0.5;
        }
    };
    $scope.random();


    var r = $scope.items.random;
    console.log(r);
    //console.log($scope.items.Percent);
    //console.log($scope.items.Percent == "");
    if ($scope.items.Percent == "") {
        $scope.items.Percent = r + "";
    } else {
        $scope.items.Percent = (r + parseFloat($scope.items.Percent)).toFixed(2);
    }
    console.log($scope.items.Percent);


    $scope.accept = function(){
        HttpService.post(REST_URL.invoke1, {fcn: "accept", args:[$scope.items.OrderID]}).then(function (response) {
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
            if ($scope.items.Percent >=1){
                HttpService.post(REST_URL.invoke1, { fcn: "expiration", args:[$scope.items.OrderID]}).then(function (response) {
                   
                });

            }else{
                 $scope.items.Request = '平仓交收';
                 HttpService.post(REST_URL.invoke1, { fcn: "expirationHalf", args:[$scope.items.OrderID, $scope.items.Percent, $scope.items.Request]}).then(function (response) {
                   
                });
            }
        });

    };


    

        $scope.reject = function () {
            HttpService.post(REST_URL.invoke1, {fcn: "reject", args: [$scope.items.OrderID]}).then(function (response) {
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
        };

        $scope.cancel = function () {
            $modalInstance.dismiss('cancel');
        };
       
}

 var MyEndModalInstanceCtrl = function ($scope, $modalInstance, items, $localStorage, REST_URL, HttpService, DialogService) {
        $scope.items = items;
        console.log(items);
        console.log($scope.items.status);
        $scope.items.event = '';
        $scope.items.option = '';
        $scope.items.des = '';


        $scope.submit = function () {
            if ($scope.items.Angel == '买入') {
                if ($scope.items.option == '继续交易') {
                    if ($scope.items.des == '') {
                        alert('描述不能为空');
                    }

                    //扣除对方信用分
                }

                if ($scope.items.option == '挂起交易') {
                    if ($scope.items.earnper > 0) {
                        $scope.items.Request = $scope.items.event;
                        HttpService.post(REST_URL.invoke1, {
                            fcn: "hang",
                            args: [$scope.items.OrderID, $scope.items.Request]
                        }).then(function (response) {
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


                    } else {
                        alert('当前处于亏损状态，不允许执行此操作');
                    }
                }


                if ($scope.items.earnPer > 0) {
                    if ($scope.items.option == '挂起交易') {

                    }

                } else {
                    alert('当前处于亏损状态，不允许执行此操作');
                }

            } else {
                alert('卖出方不允许执行此操作');
            }
        }


        $scope.cancel = function () {
            $modalInstance.dismiss('cancel');
        };
    }

    var MyOrderModalInstanceCtrl = function ($scope, $modalInstance, $modal, items, $localStorage, REST_URL, HttpService, DialogService) {
        $scope.items = items;
        if ($scope.items.Underlying == "上证50ETF"){
            $scope.items.nowPrice = 3.003;
        }else if ($scope.items.Underlying == "豆粕1912"){
            $scope.items.nowPrice = 2994;
        }else if ($scope.items.Underlying == "豆粕2001"){
            $scope.items.nowPrice = 2956;
        }else if($scope.items.Underlying == "豆粕2003"){
            $scope.items.nowPrice = 2900;
        }



        //$scope.items.nowPrice = 4200;

        if ($scope.items.Angel == "买入"){
            $scope.items.QuotePrice = $scope.items.ConstractSize * $scope.items.StrikePrice * 0.05;
            $scope.items.PayPrice = 0;
        }else{
            $scope.items.QuotePrice = 0;
            $scope.items.PayPrice = $scope.items.ConstractSize * $scope.items.StrikePrice * 0.05;
        }

        

        var now = new Date();
        var atime = Date.parse($scope.items.ExpiringDate);
        var year = now.getFullYear();
        var month = now.getMonth();
        var date = now.getDate();
        month = month + 1;
        var time = year + "-" + month + "-" + date;
        var btime = Date.parse(time);
        //console.log($scope.items);
        console.log($localStorage.loginuser);
        if (atime >= btime) {
            var c = atime - btime;
            var miao = c / 1000;
            var fen = miao / 60;
            var shi = fen / 60;
            var tian = shi / 24;
            $scope.items.vartime = tian.toFixed(0);
        }
        if (btime > atime) {
            $scope.items.vartime = '已过期';
        }

        var earnper = $scope.items.nowPrice - $scope.items.StrikePrice;
        console.log(earnper);
        if ($scope.items.Angel == "买入") {
            $scope.items.earnPer = (earnper - $scope.items.AccPrice).toFixed(2);
            $scope.items.earnTotal = ($scope.items.earnPer * $scope.items.ConstractSize).toFixed(2);
        }
        if ($scope.items.Angel == "卖出") {
            $scope.items.earnPer = -(earnper - $scope.items.AccPrice).toFixed(2);
            $scope.items.earnTotal = ($scope.items.earnPer * $scope.items.ConstractSize).toFixed(2);
        }


        $scope.items.CCR = '';
        $scope.items.EAD = '';
        $scope.items.PD = '';
        $scope.items.LGD = '';
        var itemsxq = $scope.items;
        $scope.items.flag = '';

        //只有买方可以行权
        if ($scope.items.Angel == '买入') {
            if ($scope.items.ConstractType == '美式看涨') {
                if ($scope.items.vartime >= 0 && ($scope.items.nowPrice - $scope.items.StrikePrice) >= 0) {
                    $scope.items.status = '可行权';
                }
                else {
                    $scope.items.status = '不可行权';
                }
            }

            if ($scope.items.ConstractType == '美式看跌') {
                if ($scope.items.vartime >= 0 && $scope.items.StrikePrice - $scope.items.nowPrice > 0) {
                    $scope.items.status = '可行权';
                }
                else {
                    $scope.items.status = '不可行权';
                }
            }

            if ($scope.items.ConstractType == '欧式看涨') {
                if ($scope.items.vartime <= 0 && ($scope.items.nowPrice - $scope.items.StrikePrice) > 0) {
                    $scope.items.status = '可行权';
                }
                else {
                    $scope.items.status = '不可行权';
                }
            }

            if ($scope.items.ConstractType == '欧式看跌') {
                if ($scope.items.vartime <= 0 && $scope.items.StrikePrice - $scope.items.nowPrice > 0) {
                    $scope.items.status = '可行权';
                }
                else {
                    $scope.items.status = '不可行权';
                }
            }
        } else {
            $scope.items.status = '不可行权';
        }

        $scope.randoms = function(){
            let a = Math.ceil(Math.random()*10);
            if (a<=5){
                $scope.items.rr = 0.45;
            }else{
                $scope.items.rr = 0.75;
            }
        }


       

        $scope.ccr = function () {
            $scope.randoms();
            $scope.items.LGD = 1-$scope.items.rr;
            if ($scope.items.Underlying == "上证50ETF"){
                $scope.items.PD = 0.1;
            }else{
                $scope.items.PD = 0;
            }

            $scope.items.V = 1;
            $scope.items.C = $scope.items.ConstractSize * $scope.items.StrikePrice * 0.05;
            if ($scope.items.V - $scope.items.C < 0){
                $scope.items.RC = 0;
            }else{
                $scope.items.RC = $scope.items.V - $scope.items.C;
            }
            if ($scope.items.vartime < 10){
                $scope.items.PFE = $scope.items.V * (1+0.05);
            }else if($scope.items.vartime < 30 ){
                $scope.items.PFE = $scope.items.V * (1+0.1);
            }else if($scope.items.vartime < 90){
                $scope.items.PFE = $scope.items.V * (1+0.15);
            }else{
                $scope.items.PFE = $scope.items.V * (1+0.2);
            }



            $scope.items.EAD = (1.4 *( $scope.items.RC + $scope.items.PFE)).toFixed(2);

            $scope.items.CCR = ($scope.items.EAD * $scope.items.PD * $scope.items.LGD).toFixed(2);

        };

        $scope.pingcang = function () {
            if ($scope.items.Angel == '卖出') {
                if ($localStorage.loginuser.type == "OTC"){
                    alert('OTC机构不能发起平仓');
                }else{
                    if ($scope.items.State == 'NewPublish') {
                        $scope.items.flag = 'p1';
                    } else {
                        $scope.items.flag = 'p3';
                    }
                }                
            }
            else {
                if ($scope.items.Request == '') {
                    alert('当前订单没有平仓请求，无法报价')
                } else {
                    $scope.items.flag = 'p2';
                }
            }
            console.log($scope.items.flag);
            $scope.itemsping = $scope.items;
            $modalInstance.close($scope.itemsping);
        };

        $scope.xingquan = function () {
            $scope.items.flag = 'xq';
            $scope.itemsxq = $scope.items;
            $modalInstance.close($scope.itemsxq);
        };

        $scope.end = function () {
            $scope.items.flag = 'end';
            $scope.itemsend = $scope.items;
            $modalInstance.close($scope.itemsend);
        };

        $scope.cancel = function () {
            $modalInstance.dismiss('cancel');
        };
    }
