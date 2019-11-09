app.controller('PublishProvideCtrl', ['$scope','$localStorage','$state','HttpService','REST_URL','$modal','DialogService','$q', function($scope, $localStorage, $state,HttpService,REST_URL,$modal,DialogService,$q) {

    $scope.rowCollectionPage = [];


    //  pagination
    $scope.itemsByPage=1;
    console.log($localStorage.loginuser.cmID);

    function render() {

        if($localStorage.loginuser) {

        }else {
            $state.go('access.signin');
        }

        if($localStorage.loginuser.type != "OTC"){
            alert("您没有报价权限");
        }else{


        

        HttpService.post(REST_URL.query, {fcn: "queryMyAsk", args:['A']}).then(function (response) {
            $scope.rowCollectionA = JSON.parse(response.data.message);
            $scope.loginuser = $localStorage.loginuser;

            if (!$scope.$$phase) {
                $scope.$apply();
            }
        });
        
        

            HttpService.post(REST_URL.query, {fcn: "queryMyAsk", args:['B']}).then(function (response) {
                $scope.rowCollectionB = JSON.parse(response.data.message);
                $scope.loginuser = $localStorage.loginuser;

                if (!$scope.$$phase) {
                    $scope.$apply();
                }
            });
        
        
            HttpService.post(REST_URL.query, {fcn: "queryMyAsk", args:['C']}).then(function (response) {
                $scope.rowCollectionC = JSON.parse(response.data.message);
                $scope.loginuser = $localStorage.loginuser;

                if (!$scope.$$phase) {
                    $scope.$apply();
                }
            });
        }
        
    }
    render();
        $scope.open = function(row) {
               if (row.PublishID==$localStorage.loginuser.cmId){
                alert('自己发布的询价不能报价');
               }//else if ($localStorage.loginuser.type != "OTC"){
                //alert('您没有报价权限');
               //}
               else{
                var modalInstance = $modal.open({  
                    templateUrl : 'myBillInfo.html',  
                    controller : MyBillModalInstanceCtrl,  
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
            }
                
        };  
}]);
var MyBillModalInstanceCtrl = function($scope, $modalInstance, items, $localStorage, REST_URL, HttpService, DialogService) {
    $scope.items = items;
    var ask = items.AskNum;
    console.log($localStorage.loginuser.cmId);
    console.log(items);
    $scope.items.ensurePrice = '';
    $scope.items.Quota = '';

    if ($scope.items.Angel == "买入"){
        //document.getElementById('ensurePrice').setAttribute('disabled', 'disabled');
        //document.getElementById('Quota').setAttribute('disabled', 'disabled');
        jQuery('#ensurePrice').attr('disabled', 'disabled');
        jQuery('#Quota').attr('disabled', 'disabled');
        //jQuery('#ensurePrice').removeAttr('disabled');  // 删除指定element的属性，使用jQuery方法
        //jQuery('#Quota').removeAttr('disabled');  // 删除指定element的属性，使用jQuery方法
    }





    $scope.ok = function (){
         var doupo1912 =  [
        ["3250","4","279"],
        ["3200", "5", "242.5"],
        ["3150", "9.5","196"],
        ["3100", "13.5", "147"],
        ["3050", "22", "111.5"],
        ["3000", "34.5", "79"],
        ["2950", "53","43"],
        ["2900", "79", "21.5"],
        ["2850", "122", "9"],
        ["2800",  "168", "4.5"],
        ["2750", "214.5", "3"],
        ["2700",  "264", "2.5"],
        ["2650", "313.5", "1.5"],
        ["2600", "390", "1"],
        ["2550", "429", "0.5"],
        ["2500", "474", "0.5"],
        ["2450",  "524", "0.5"],
        ["2400",  "588", "0.5"]
        ];

        var doupo2001 = [
        ["3250", "4",   "279"],
        ["3200", "5",   "242.5"],
        ["3150", "9.5", "196"],
        ["3100", "13.5", "147"],
        ["3050","22", "111.5"],
        ["3000", "34.5", "79"],
        ["2950", "53", "43"],
        ["2900", "79", "21.5"],
        ["2850", "122", "9"],
        ["2800", "168", "4.5"],
        ["2750","214.5", "3"],
        ["2700", "264", "2.5"],
        ["2650", "313.5", "1.5"],
        ["2600", "390", "1"],
        ["2550", "429", "0.5"],
        ["2500", "474", "0.5"],
        ["2450", "524", "0.5"],
        ["2400", "588", "0.5"]
        ];
        
        var doupo2003 = [
        ["3200", "13",  "290.5"],
        ["3150", "19.5", "247"],
        ["3100", "28.5",  "206"],
        ["3050", "40.5",  "168.5"],
        ["3000", "56",  "134"],
        ["2950", "75.5", "103.5"],
        ["2900", "99.5", "77.5"],
        ["2850", "128", "56"],
        ["2800", "160.5", "39"],
        ["2750", "197.5", "26"],
        ["2700", "238", "16.5"],
        ["2650", "281", "10"],
        ["2600", "326.5", "5.5"],
        ["2550", "374", "3"],
        ["2500", "422.5", "1.5"],
        ["2450",  "472", "0.5"],
        ["2400", "522", "0.5"],
        ];

        if ($scope.items.Biaodi == "豆粕1912"){
            var tem1 =0;
            var tem2 = 0;
            for(var i=0;i<doupo1912.length;i++){
                if ($scope.items.ExecuPrice>doupo1912[i][0] && i>0){
                    tem1=Math.abs($scope.items.ExecuPrice-doupo1912[i][0]);
                    tem2=Math.abs($scope.items.ExecuPrice-doupo1912[i-1][0]);
                    if (tem1-tem2>0){
                        if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo1912[i-1][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo1912[i-1][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo1912[i-1][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo1912[i-1][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }

                    }else if (tem1-tem2<0){
                        if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo1912[i][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo1912[i][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo1912[i][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo1912[i][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }
                    }
                    break;
                }else if($scope.items.ExecuPrice>doupo1912[i][0] && i==0){
                    tem1=Math.abs($scope.items.ExecuPrice-doupo1912[i][0]);
                    if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo1912[i][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo1912[i][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo1912[i][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo1912[i][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }
                }
            }
        }




        if ($scope.items.Biaodi == "豆粕2001"){
            var tem1 =0;
            var tem2 = 0;
            for(var i=0;i<doupo2001.length;i++){
                if ($scope.items.ExecuPrice>doupo1912[i][0] && i>0){
                    tem1=Math.abs($scope.items.ExecuPrice-doupo2001[i][0]);
                    tem2=Math.abs($scope.items.ExecuPrice-doupo2001[i-1][0]);
                    if (tem1-tem2>0){
                        if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2001[i-1][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2001[i-1][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2001[i-1][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2001[i-1][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }

                    }else if (tem1-tem2<0){
                        if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2001[i][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2001[i][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2001[i][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2001[i][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }
                    }
                    break;
                }else if($scope.items.ExecuPrice>doupo2001[i][0] && i==0){
                    tem1=Math.abs($scope.items.ExecuPrice-doupo2001[i][0]);
                    if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2001[i][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2001[i][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2001[i][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2001[i][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }
                }
            }
        }


        if ($scope.items.Biaodi == "豆粕2003"){
            var tem1 =0;
            var tem2 = 0;
            for(var i=0;i<doupo1912.length;i++){
                if ($scope.items.ExecuPrice>doupo1912[i][0] && i>0){
                    tem1=Math.abs($scope.items.ExecuPrice-doupo2003[i][0]);
                    tem2=Math.abs($scope.items.ExecuPrice-doupo2003[i-1][0]);
                    if (tem1-tem2>0){
                        if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2003[i-1][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2003[i-1][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2003[i-1][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2003[i-1][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }

                    }else if (tem1-tem2<0){
                        if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2003[i][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2003[i][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2003[i][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2003[i][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }
                    }
                    break;
                }else if($scope.items.ExecuPrice>doupo2003[i][0] && i==0){
                    tem1=Math.abs($scope.items.ExecuPrice-doupo2003[i][0]);
                    if ($scope.items.ConstractType == "美式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2003[i][1];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2003[i][1];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "美式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = doupo2003[i][2];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = doupo2003[i][2];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }
                }
            }
        }




        var shangzheng = [
        //StrikePrice Call_202006 Put_202006  Call_202003 Put_202003  Call_201912 Put_201912  Call_201911 Put_201911
        ["2.5", "", "", "", "", "0.5069", "0.0016", "" ,""],
        ["2.55", "", "", "", "", "0.4558", "0.0017", "", ""],
        ["2.6", "",    "",  "0.4212",  "0.0085",  "0.4044",  "0.0021",  "", ""],
        ["2.65", "",  "",  "0.3717",  "0.0115",  "0.3608",  "0.0025",  "",  ""],
        ["2.7", "", "", "0.3279", "0.016", "0.3097",  "0.0038",  "", ""],
        ["2.75", "", "", "0.2855",  "0.0227",  "0.2581",  "0.0051",  "0.2573",  "0.0014"],
        ["2.8", "0.2791",  "0.0508",  "0.2453",  "0.0315",  "0.2156",  "0.0084",  "0.2079",  "0.0021"],
        ["2.85", "0.2445",  "0.0661",  "0.2065",  "0.042",   "0.1717",  "0.0148",  "0.1599",  "0.0045"],
        ["2.9", "0.215",   "0.0832",  "0.1733",  "0.058",   "0.1323",  "0.0251",  "0.1158",  "0.0099"],
        ["2.95", "0.1846",  "0.1046",  "0.1421",  "0.0787",  "0.0969",  "0.0406",  "0.0783",  "0.02"],
        ["3","0.1608", "0.1271", "0.117", "0.1019",  "0.0681",  "0.0624",  "0.048",   "0.0397"],
        ["3.1", "0.1162",  "0.1831",  "0.0754",  "0.1603",  "0.031",   "0.1249",  "0.0139",  "0.106"],
        ["3.2","0.0818",  "0.2484",  "0.0466",  "0.2315",  "0.0129",  "0.2063",  "0.0038",  "0.1982"],
        ["3.3", "0.0567",  "0.3215", "0.0289", "0.3129",  "0.0059",  "0.2991",  "0.0015",  "0.2955"],
        ["3.4", "0.0391",  "0.4013",  "0.0178",  "0.3986",  "0.0037",  "0.3972",  "0.0007",  "0.3958"],
        ["3.5", "",    "",    "0.0115",  "0.4933",  "0.0026",  "0.4949",  "0.0005",  "0.4952"]
        ];


        if ($scope.items.Biaodi == "上证50ETF"){
            var tem1 =0;
            var tem2 = 0;
            for(var i=0; i<shangzheng.length;i++){
                 if ($scope.items.ExecuPrice<shangzheng[i][0] && i>0){
                    tem1=Math.abs($scope.items.ExecuPrice-shangzheng[i][0]);
                    tem2=Math.abs($scope.items.ExecuPrice-shangzheng[i-1][0]);
                    if (tem1-tem2>0){
                        if ($scope.items.ConstractType == "欧式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = shangzheng[i-1][5];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = shangzheng[i-1][5];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "欧式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = shangzheng[i-1][6];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = shangzheng[i-1][6];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }

                    }else if (tem1-tem2<0){
                        if ($scope.items.ConstractType == "欧式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = shangzheng[i][5];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = shangzheng[i][5];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "欧式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = shangzheng[i][6];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = shangzheng[i][6];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }
                    }
                    break;
                 }else if ($scope.items.ExecuPrice<shangzheng[i][0] && i==0){
                    if ($scope.items.ConstractType == "欧式看涨"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = shangzheng[i][5];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = shangzheng[i][5];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }else if($scope.items.ConstractType == "欧式看跌"){
                            if ($scope.items.Angel == "买入"){
                                $scope.items.Price1 = shangzheng[i][6];
                                $scope.items.PricePingPer = (($scope.items.Price1/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c1=true;
                            }else if ($scope.items.Angel == "卖出"){
                                $scope.items.Price1Per = shangzheng[i][6];
                                $scope.items.PricePing = (($scope.items.Price1Per/$scope.items.ExecuPrice)*100).toFixed(2);
                                $scope.items.c2 = true;
                            }
                        }                
                 }                   
            }
        }       
    };

    $scope.clear = function (){
        $scope.items.Price1 = '';
        $scope.items.PricePingPer =''; 
        $scope.items.Price1Per = '';
        $scope.items.PricePing = '';
        $scope.items.c1 = false;
        $scope.items.c2 = false;
    }
    //  args: 0 - AskNum ; 1 - ProvideID ; 2 - price; 3 - provideNum
    //bill.BillInfoID, selected.item.EndrCmID,selected.item.EndrAcct

    $scope.submit = function (){
        console.log($scope.items.AskNum);
        console.log($localStorage.loginuser.cmId);
        if ($scope.items.Angel == "买入" && ($scope.items.ensurePrice != '' || $scope.items.Quota != '')){
            alert('方向为买入时不允许填写保证金和授信额度');
        }else{


        var price;
        if ($scope.items.Price1!= ""){
           price = $scope.items.Price1;
        }else{
            price = $scope.items.Price1Per;
        }


        
        console.log(price);
        var myDate = new Date();
        var s = myDate.getMonth().toString()+myDate.getDate().toString()+myDate.getHours().toString()+myDate.getMinutes().toString()+myDate.getSeconds().toString();
        console.log(s);


        HttpService.post(REST_URL.invoke, { fcn: "bid", args:[$scope.items.AskNum, $localStorage.loginuser.cmId, price, s]}).then(function (response) {
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
    }
    $scope.cancel = function (){
        $modalInstance.dismiss('cancel');
    };
}
