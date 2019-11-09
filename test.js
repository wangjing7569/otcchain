//'use strict';
/*
var fs = require('fs');
var params = {
    "username": "Ann",
    "name": "D公司",
    "passwd": "123456",
    "cmID": "DCMID",
    "Acct": "D公司"
}//在真实的开发中id肯定是随机生成的而且不会重复的，下一篇写如何生成随机切不会重复的随机数，现在就模拟一下假数据
//写入json文件选项
function writeJson(params){
    //现将json文件读出来
    fs.readFile('./config.json',function(err,data){
        if(err){
            return console.error(err);
        }
        var person = data.toString()//将二进制的数据转换为字符串
        person = JSON.parse(person);//将字符串转换为json对象
        person.mockupUsers.push(params);//将传来的对象push进数组对象中
        console.log(person.mockupUsers);
        var str = JSON.stringify(person);//因为nodejs的写入文件只认识字符串或者二进制数，所以把json对象转换成字符串重新写入json文件中
        fs.writeFile('./config.json',str,function(err){
            if(err){
                console.error(err);
            }
            console.log('----------新增成功-------------');
        })
    })
}
writeJson(params)//执行一下;*/

var fs=require('fs');
var args = [{
    "BillInfoID":"POC10000988",
    "BillInfoAmt":"10001",
    "BillInfoType":"AC01",
    "BillInfoIsseDate":"20161001",
    "BillInfoDueDate":"20161012",
    "DrwrCmID":"ChupiaoId",
    "DrwrAcct":"C11111111",
    "AccptrCmID":"ChengduiId",
    "AccptrAcct":"C11111111",
    "PyeeCmID":"ShoukuanId",
    "PyeeAcct":"S11111111",
    "HodrCmID":"BCMID",
    "HodrAcct":"B公司"}]

args=JSON.stringify(args);
args=JSON.stringify(args);
console.log(args)

console.log((new Date()).Format("yyyy-MM-dd hh:mm:ss"));
//args=JSON.parse(args);
//console.log(args)
