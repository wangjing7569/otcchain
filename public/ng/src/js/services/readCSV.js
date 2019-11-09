'use strict';

myModule.factory('Items', ['$http', function($http){
  var Url   = "csv/2019-10-31.csv";
  var Items = $http.get(Url).then(function(response){
     return csvParser(response.data);
  });
  return Items;
}]);