'use strict';

var crypto = require('crypto');
var log4js = require('log4js');
var logger = log4js.getLogger('SampleWebApp');

exports.aes_algorithm = "aes-128-ecb";
exports.aes_secrect = "HdMW6jdVT5dLDpYTpnWJk2uqSBFRXebt";

exports.encrypt = function (text) {
    var cipher = crypto.createCipher(this.aes_algorithm, this.aes_secrect)
    var crypted = cipher.update(text, 'utf8', 'hex')
    crypted += cipher.final('hex');
    return crypted;
};

exports.decrypt = function (text) {
    var decipher = crypto.createDecipher(this.aes_algorithm, this.aes_secrect)
    var dec = decipher.update(text, 'hex', 'utf8')
    dec += decipher.final('utf8');
    return dec;
};
/*
var json={
    "agentId":"1",
    "agentUserName":"test1",
    "agentCode":"3",
    "agentOrgName":"org1",
    "agentPassword":"12345678",
    "agentName":"wang",
    "agentDep":"产品",
    "agentTel":"13212345678"
}
*/

/*
var json={
    "userId":"1",
    "userName":"wangjing",
    "orgName":"org1",
    "password":"12345678",
    "linkName":"张三",
    "userType":"1",
    "companyName":"筑客网络",
    "creditNum":"318372632932023x",
    "busLicenseNum":"276325243254",
    "orgLicenseNum":"28327325326452",
    "taxLicenseNum":"uyqee7e625",
    "fund":"100W",
    "companyType":"4",
    "province":"上海",
    "city":"上海",
    "area":"黄埔区",
    "address":"蒙自路",
    "legalName":"李四",
    "idcard":"313211198810192345",
    "isTrinity":"1",
    "orgLicenseUrl":"./images/orgUrl/1",
    "busLicenseUrl":"./images/busUrl/1",
    "taxLicenseUrl":"./images/taxUrl/1",
    "logoUrl":"./images/logoUrl/1",
    "cardOnUrl":"./images/cardOnUrl/1",
    "cardBackUrl":"./images/cardBackUrl/1",
    "accountName":"筑客网络",
    "bankAccount":"62122618907652173",
    "bankName":"工商银行",
}

var json={
    "billId":"1",
    "billType":"电子商票",
    "imgOnUrl":"./images/On/1",
    "imgCodeOnUrl":"./images/CodeOn/1",
    "imgCodeBackList":"./images/CodeBack/1,./images/CodeBack/2",
    "imgBackUrlList":"./images/BackUrl/1,./images/BackUrl/2",
    "billNo":"2018ABCD789",
    "acceptorName":"南通三建",
    "bearerName":"筑客网络",
    "bearerId":"1",
    "amount":"100W",
    "drawDate":"2018-06-12",
    "endDate":"2019-06-12",
    "flawState":"1,2,3,4",
    "endorseTimes":"0",
    "deductMoney":"80W",
    "passTime":"2019-01-24",
    "requirementEndDate":"2019-02-01"
}

var json={
    "billId":"2",
    "billType":"电子商票",
    "imgOnUrl":"./images/On/2",
    "imgCodeOnUrl":"./images/CodeOn/2",
    "imgCodeBackList":"./images/CodeBack/3,./images/CodeBack/4",
    "imgBackUrlList":"./images/BackUrl/3,./images/BackUrl/4",
    "billNo":"2018ABCD7hu",
    "acceptorName":"南通三建",
    "bearerName":"筑客网络",
    "bearerId":"1",
    "amount":"120W",
    "drawDate":"2018-06-12",
    "endDate":"2019-06-12",
    "flawState":"1,2,3,4",
    "endorseTimes":"0",
    "passTime":"2019-01-24",
    "requirementEndDate":"2019-02-01"
}

var json={
    "userId":"2",
    "userName":"buybuybuy",
    "orgName":"org1",
    "password":"12345678",
    "linkName":"李四",
    "userType":"2",
    "companyName":"未来科技公司",
    "creditNum":"318372637632023x",
    "busLicenseNum":"2763223873254",
    "orgLicenseNum":"2282644265",
    "taxLicenseNum":"272hhdgwwe",
    "fund":"1000W",
    "companyType":"4",
    "province":"上海",
    "city":"上海",
    "area":"黄埔区",
    "address":"马当路",
    "legalName":"李四",
    "idcard":"313211198810192367",
    "isTrinity":"1",
    "orgLicenseUrl":"./images/orgUrl/2",
    "busLicenseUrl":"./images/busUrl/2",
    "taxLicenseUrl":"./images/taxUrl/2",
    "logoUrl":"./images/logoUrl/2",
    "cardOnUrl":"./images/cardOnUrl/2",
    "cardBackUrl":"./images/cardBackUrl/2",
    "accountName":"未来科技公司",
    "bankAccount":"62122618907652123",
    "bankName":"农业银行",
}*/



/*
var json={
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI5OSIsInVzZXJuYW1lIjoidGVzdCIsIm9yZ25hbWUiOiJvcmcxIiwidXNlclR5cGUiOiIzIiwiaWF0IjoxNTQ4ODExODc1LCJleHAiOjE1NDg4OTgyNzV9._zyM8ctBmR8VEfT6BiR_sniloTkSE7ePkF_aW07DM5Y"
}


var json={
    "userId":"2",
    "userName":"buybuybuy",
    "orgName":"org1",
    "password":"12345678",
    "userType":"2"
}
var json={
    "billId":"2",
    "userId":"2",
    "userName":"buybuybuy",
    "companyName":"未来科技公司",
    "deductMoney":"1W"
}

var json= {
    "billId": "2",
    "userName":"buybuybuy",
    "companyName":"未来科技公司",
    "deductMoney":"1W",
    "endorsementVoucherUrl":"/images/endorsementVoucher/1"
}

var json={
    "lastTime":"2019-02-21 17:13:00",
    "orgLicenseUrl":"http://zhujc-oss1.oss-cn-shenzhen.aliyuncs.com/2019-01/G44JXsEdCN.png",
    "accountName":"魔",
    "city":"1",
    "companyName":"打酱油的飘过",
    "bankName":"中国建设银行濮阳采油二厂支行储蓄所",
    "linkName":"陈帅",
    "legalName":"仲华",
    "password":"e10adc3949ba59abbe56e057f20f883e",
    "province":"1",
    "busLicenseUrl":"http://zhujc-oss1.oss-cn-shenzhen.aliyuncs.com/2019-01/G44JXsEdCN.png",
    "id":26,
    "taxLicenseUrl":"http://zhujc-oss1.oss-cn-shenzhen.aliyuncs.com/2019-01/G44JXsEdCN.png",
    "area":"1",
    "bankAccount":"1254612555616554",
    "cardOnUrl":"http://zhujc-oss1.oss-cn-shenzhen.aliyuncs.com/2019-01/yPsmcG3J4X.jpg",
    "cardBackUrl":"http://zhujc-oss1.oss-cn-shenzhen.aliyuncs.com/2019-01/hdWeiFZQbn.jpg",
    "address":"龙子湖大学城",
    "companyType":"1",
    "isTrinity":1,
    "userName":"15865320125",
    "userId":26,
    "logoUrl":"http://zhujc-oss1.oss-cn-shenzhen.aliyuncs.com/2018-12/xHPYN6fGfd.jpg",
    "busLicenseNum":"333333232323",
    "taxLicenseNum":"333333232323",
    "fund":100.0,
    "idcard":"12345678963546",
    "creditNum":"123",
    "orgLicenseNum":"333333232323",
    "userType":1,
    "time":"2019-02-21 17:13:00"
}
*/

var json= {
    "billId":2,
    "userId":5,
    "recordId":254,
    "time":"2019-03-08 14:22:30",
    "lastTime":"2019-03-08 14:22:30"
}


if(!json.bileId){
    console.log("shdf")
}

console.log(json.billId.toString())
var s = JSON.stringify(json);
console.log(s);
var so= JSON.parse(s);
console.log(so);
var hw = this.encrypt(s);
var tt=this.decrypt(hw)
var so= JSON.parse(tt);
console.log(so);
console.log(hw);
console.log(this.decrypt(hw));


