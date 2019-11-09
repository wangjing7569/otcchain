'use strict'
var log4js = require('log4js');

var json1={
    "id":11,
    "username":"wangjing"
}

var json2={
    "id":22,
    "username":"wangjing"
}

var ss=[json1,json2]

for(var i=0;i<ss.length;i++){
    if(ss[i].id==11){
        console.log("test good")
    }
}

var ss="skhg"
var sm="jagd"+ss
console.log(sm)