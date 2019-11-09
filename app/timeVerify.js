'use strict'
var log4js = require('log4js');
var sd = require('silly-datetime');

exports.timeVerify=function(oldTime){
    var timestamp=Date.parse(new Date());
    var timestamp2=Date.parse(new Date(oldTime));
    var timeSub=(timestamp-timestamp2)/1000;
    return timeSub;
}

