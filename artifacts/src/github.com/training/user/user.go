package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"log"
)


// 买卖方用户帐号状态
const (
	UserInfo_State_Valid="Valid"      //账户有效（成功注册/修改，未登录）
	UserInfo_State_Logining="Logining"    //登录中
	UserInfo_State_Frozen="Frozen"      //账户冻结
	UserInfo_State_Canceld="Canceled"    //账户注销
)
const Users_Prefix="Users_"
//交易用户信息
type Users struct {
	UserName string `json:"userName"`       //用户名
	UserId string `json:"userId"`       //用户主键id    0～10W(筑票票)
	Password string `json:"password"`   //密码
	LinkName string `json:"linkName"`    //联系人姓名
	UserType string `json:"userType"`    //用户类型     1.卖方  2.买方
	CompanyName string `json:"companyName"`   //公司名称
	CreditNum string `json:"creditNum"`    //统一社会信用代码
	BusLicenseNum string `json:"busLicenseNum"`   //营业执照注册号
	OrgLicenseNum string `json:"orgLicenseNum"`    //组织机构代码证号
	TaxLicenseNum string `json:"taxLicenseNum"`     //税务登记证号码
	Fund string `json:"fund"`            //注册资本
	CompanyType string `json:"companyType"`     //1.国有企业 2.中央企业 3.上市企业 4.民营企业 99.其他
	Address string `json:"address"`                //公司地址
	LegalName string `json:"legalName"`     //法人姓名
	Idcard string `json:"idcard"`    //法人身份证号
	IsTrinity string `json:"isTrinity"`    //是否三证合一     0.否 1.是
	OrgLicenseUrl string `json:"orgLicenseUrl"`       //组织机构证url地址
	BusLicenseUrl string `json:"busLicenseUrl"`     //营业执照url地址
	TaxLicenseUrl string `json:"taxLicenseUrl"`      //税务登记证url地址
	LogoUrl string `json:"logoUrl"`                   //公司logo  url地址
	CardOnUrl string `json:"cardOnUrl"`               //身份证正面url地址
	CardBackUrl string `json:"cardBackUrl"`         //身份证背面url地址
	AccountName string `json:"accountName"`         //企业户名
	BankAccount string `json:"bankAccount"`         //银行帐号
	BankName string `json:"bankName"`               //开户银行名称
	Operation string `json:"operation"`         //操作说明
	BankNo string `json:"bankNo"`            //开户银行行号
	AccountState string `json:"accountState"`        //State
	History []HistoryItem `json:"History"`        //用户历史
}

type HistoryItem struct {
	TxId string `json:"txId"`
	Users Users `json:"users"`
}

//chaincode response 结构
type chaincodeRet struct {
	Code int // 0表示成功，其他情况为1
	Des string //相关描述
}

type UsersChaincode struct {

}

//response消息格式(byte)
func getRetByte(code int,des string) []byte{
	var r chaincodeRet
	r.Code=code
	r.Des=des
	b,err :=json.Marshal(r)    //将response消息转换成json字符串格式
	if err!=nil {
		fmt.Println("marshal Ret failed") //生成json字符串错误
		return nil
	}
	return b
}

//response 消息日志格式(string)
func getRetString(code int,des string) string{
	var r chaincodeRet
	r.Code=code
	r.Des=des
	b,err:=json.Marshal(r)
	if err!=nil{
		fmt.Println("marshal Ret failed")
		return ""
	}
	//chaincodeLogger.Infof("%s",string(b[:]))
	return string(b[:])
}

//根据userId获取身份信息
func (a *UsersChaincode) getUsers(stub shim.ChaincodeStubInterface,userId string) (Users,bool){
	var users Users
	key := Users_Prefix + userId
	b,err := stub.GetState(key)
	if b==nil{
		return users,false
	}
	err=json.Unmarshal(b,&users)
	if err!=nil{
		return users,false
	}
	return users,true
}

//保存身份信息
func (a *UsersChaincode) putUsers(stub shim.ChaincodeStubInterface, users Users) ([]byte,bool){
	byte,err := json.Marshal(users)
	if err!=nil{
		return nil,false
	}

	err=stub.PutState(Users_Prefix+users.UserId,byte)
	if err!=nil{
		return nil,false
	}
	return byte,true
}

//链码初始化 Init接口
func (a *UsersChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
	return shim.Success(nil)
}

//用户注册/修改
//args:0-{Users Object}
func (a *UsersChaincode) register(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args)!=1{
		res := getRetString(1,"zpp Invoke register args!=1")
		return shim.Error(res)
	}
	var users Users
	err:=json.Unmarshal([]byte(args[0]),&users)
	if err!=nil{
		res:=getRetString(1,"zpp Invoke register unmarshal failed")
		return shim.Error(res)
	}
	//更改并存储客服信息，更新客服帐号状态：帐号状态设置为有效"Valid"
	users.AccountState=UserInfo_State_Valid
	//保存客服信息
	_,us :=a.putUsers(stub,users)
	if !us{
		res := getRetString(1,"zpp invoke register put users failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(users)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//登录
//args：0-userId,1-userName,2-userType,3-password
func (a *UsersChaincode) login(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args)!=5{
		res := getRetString(1,"zpp invoke login args error")
		return shim.Error(res)
	}
	//根据userId获取客服身份信息
	users,us := a.getUsers(stub,args[0])
	if !us{
		res := getRetString(1,"zpp invoke login get users error,not exist")
		return shim.Error(res)
	}
	//验证账户状态(冻结或注销不可用)
	if users.AccountState==UserInfo_State_Frozen{
		res:= getRetString(1,"accountState is Frozen")
		return shim.Error(res)
	}
	if users.AccountState==UserInfo_State_Canceld{
		res:= getRetString(1,"accountState is Canceled")
		return shim.Error(res)
	}
	//信息无误
	users.AccountState=UserInfo_State_Logining
	users.Operation=args[4]
	//保存更改信息
	_,us =a.putUsers(stub,users)
	if !us{
		res := getRetString(1,"zpp invoke login put users error")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(users)
	if err!=nil{
		log.Fatal(err)
	}
	res:= getRetByte(0,string(jsonStr))
	return shim.Success(res)
}
//客服登出
//args:0-userId
//前端解码传输
func (a *UsersChaincode) logout(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	//如果传入参数数量不对
	if len(args)!=2{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	//根据userid获取客服身份信息
	users,us := a.getUsers(stub,args[0])
	if !us{
		res := getRetString(1,"zpp invoke logout get users error,not exist")
		return shim.Error(res)
	}
	users.AccountState=UserInfo_State_Valid   //账户变成有效状态
	users.Operation=args[1]
	//保存更改信息
	_,us =a.putUsers(stub,users)
	if !us {
		res := getRetString(1, "zpp invoke logout put users error")
		return shim.Error(res)
	}
	//返回成功信息
	jsonStr,err:=json.Marshal(users)
	if err!=nil{
		log.Fatal(err)
	}
	res:= getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//冻结(违规操作之后冻结帐号)
//args:0-userId
func (a *UsersChaincode) frozen(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	//参数数量不对
	if len(args)!=2{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	//根据userid获取客服身份信息
	users,us := a.getUsers(stub,args[0])
	if !us{
		res := getRetString(1,"zpp invoke frozen get users error,not exist")
		return shim.Error(res)
	}
	//查询状态,非注销用户即可冻结
	if users.AccountState==UserInfo_State_Canceld{
		res := getRetString(1,"wrong accountState error")
		return shim.Error(res)
	}
	users.AccountState=UserInfo_State_Frozen
	users.Operation=args[1]
	//保存更改信息
	_,us =a.putUsers(stub,users)
	if !us{
		res := getRetString(1,"zpp invoke frozen put agent error")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(users)
	if err!=nil{
		log.Fatal(err)
	}
	res:= getRetByte(0,string(jsonStr))
	return shim.Success(res)
}
//帐号解冻
//args:0-userId
func (a *UsersChaincode) unfrozen(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	//参数数量不对
	if len(args)!=2{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	users,us := a.getUsers(stub,args[0])
	if !us{
		res := getRetString(1,"zpp invoke frozen get users error,not exist")
		return shim.Error(res)
	}
	users.AccountState=UserInfo_State_Valid
	users.Operation=args[1]
	//保存更改信息
	_,us =a.putUsers(stub,users)
	if !us{
		res := getRetString(1,"zpp invoke frozen put users error")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(users)
	if err!=nil{
		log.Fatal(err)
	}
	res:= getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//客服帐号注销/删除
//args:0-userId
func (a *UsersChaincode) delete(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	//参数数量不对
	if len(args)!=2{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	users,us := a.getUsers(stub,args[0])
	if !us{
		res := getRetString(1,"zpp invoke frozen get users error,not exist")
		return shim.Error(res)
	}
	//查询状态,已删除用户不可再删除
	if users.AccountState==UserInfo_State_Canceld{
		res := getRetString(1,"wrong accountState error")
		return shim.Error(res)
	}
	users.AccountState=UserInfo_State_Canceld
	users.Operation=args[1]
	//保存更改信息
	_,us =a.putUsers(stub,users)
	if !us{
		res := getRetString(1,"zpp invoke delete put users error")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(users)
	if err!=nil{
		log.Fatal(err)
	}
	res:= getRetByte(0,string(jsonStr))
	return shim.Success(res)
}
//modify password
//args:0-userId,1-password
func (a *UsersChaincode) modifyPassword(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args)!=3{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	user,us := a.getUsers(stub,args[0])
	if !us{
		res := getRetString(1,"zpp invoke frozen get user error,not exist")
		return shim.Error(res)
	}
	user.Password=args[1]
	user.Operation=args[2]
	//保存更改信息
	_,us =a.putUsers(stub,user)
	if !us{
		res := getRetString(1,"zpp invoke modifyPassword put user error")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(user)
	if err!=nil{
		log.Fatal(err)
	}
	res:= getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

func (a *UsersChaincode) queryById(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=2{
		res:=getRetString(1,"zpp queryById args error")
		return shim.Error(res)
	}
	users,us:=a.getUsers(stub,args[0])
	if !us{
		res := getRetString(1,"zpp queryById get users error,not exist")
		return shim.Error(res)
	}
	resultsIterator, err := stub.GetHistoryForKey(Users_Prefix+args[0])
	if err != nil {
		res := getRetString(1,"zpp queryById GetHistoryForKey error")
		return shim.Error(res)
	}
	defer resultsIterator.Close()
	var history []HistoryItem
	var hisAgent Users
	for resultsIterator.HasNext(){
		historyData,err:=resultsIterator.Next()
		if err!=nil{
			res:=getRetString(1,"zpp queryById resultsIterator.Next() error")
			return shim.Error(res)
		}
		var hisItem HistoryItem
		hisItem.TxId=historyData.TxId
		json.Unmarshal(historyData.Value,&hisAgent)
		if historyData.Value==nil{
			var emptyAgent Users
			hisItem.Users=emptyAgent
		} else{
			json.Unmarshal(historyData.Value,&hisAgent)
			hisItem.Users=hisAgent
		}
		history=append(history,hisItem)
	}
	users.History=history
	users.Operation=args[1]
	age,err:=json.Marshal(users)
	if err != nil {
		res := getRetString(1,"zpp Marshal queryById usersList error")
		return shim.Error(res)
	}
	return shim.Success(age)
}


//chaincode invoke 接口
func (a *UsersChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
	function,args:=stub.GetFunctionAndParameters()
	//chaincodeLogger.Info("%s%s","zppChaincode function=",function)
	//chaincodeLogger.Info("%s%s","zppChaincode args=",args)
	//invoke
	if function=="register"{
		return a.register(stub,args)
	} else if function=="login"{
		return a.login(stub,args)
	} else if function=="logout"{
		return a.logout(stub,args)
	} else if function=="frozen"{
		return a.frozen(stub,args)
	} else if function=="unfrozen" {
		return a.unfrozen(stub,args)
	}else if function=="delete"{
		return a.delete(stub,args)
	}else if function=="modifyPassword"{
		return a.modifyPassword(stub,args)
	}else if function=="queryById"{
		return a.queryById(stub,args)
	}

	res:=getRetString(1,"Unknown method")
	//chaincodeLogger.Info("%s",res)
	//chaincodeLogger.Infof("%s",res)
	return shim.Error(res)
}

func main() {
	if err := shim.Start(new(UsersChaincode)); err != nil {
		fmt.Printf("Error starting BillChaincode: %s", err)
	}
}