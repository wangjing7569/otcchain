package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"log"
)

// logger
//var chaincodeLogger = flogging.MustGetLogger("ChainnovaChaincode")
// 客服帐号状态
const (
	AgentInfo_State_Valid="Valid"      //账户有效（成功注册/修改，未登录）
	AgentInfo_State_Logining="Logining"    //登录中
	AgentInfo_State_Frozen="Frozen"      //账户冻结
	AgentInfo_State_Canceld="Canceled"    //账户注销
)
//客服key的前缀
const Agent_Prefix="Agent_"
//客服Object结构
type Agent struct {
	Id string `json:"id"`    //用户id,key
	Username string `json:"username"`   //用户名
	Password string `json:"password"`   //客服密码
	ServerName string `json:"serverName"`           //客服人员姓名
	Department string `json:"department"`             //客服人员所属部门
	Phone string `json:"phone"`             //手机号码
	Operation string `json:"operation"`
	AccountState string `json:"accountState"`      //账户状态
	History []HistoryItem `json:"History"`        //客服用户历史
}

//客服用户item结构
type HistoryItem struct {
	TxId string `json:"txId"`
	Agent Agent `json:"agent"`
}

//chaincode response 结构
type chaincodeRet struct {
	Code int // 0表示成功，其他情况为1
	Des string //相关描述
}

//chaincode
type AgentChaincode struct {

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
//根据agentId获取客服身份信息
func (a *AgentChaincode) getAgent(stub shim.ChaincodeStubInterface,Id string) (Agent,bool){
	var agent Agent
	key := Agent_Prefix + Id
	b,err := stub.GetState(key)
	if b==nil{
		return agent,false
	}
	err=json.Unmarshal(b,&agent)
	if err!=nil{
		return agent,false
	}
	return agent,true
}

//保存客服身份信息
func (a *AgentChaincode) putAgent(stub shim.ChaincodeStubInterface, agent Agent) ([]byte,bool){
	byte,err := json.Marshal(agent)
	if err!=nil{
		return nil,false
	}


	err=stub.PutState(Agent_Prefix+agent.Id,byte)
	if err!=nil{
		return nil,false
	}
	return byte,true
}


//链码初始化 Init接口
func (a *AgentChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
	return shim.Success(nil)
}

//客服注册/修改
//args:0-{Agent Object}
func (a *AgentChaincode) register(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args)!=1{
		res := getRetString(1,"zpp Invoke register args!=1")
		return shim.Error(res)
	}
	var agent Agent
	err:=json.Unmarshal([]byte(args[0]),&agent)
	if err!=nil{
		res:=getRetString(1,"zpp Invoke issue unmarshal failed")
		return shim.Error(res)
	}
	//更改并存储客服信息，更新客服帐号状态：帐号状态设置为有效"Valid"
	agent.AccountState=AgentInfo_State_Valid
	//保存客服信息
	_,ag :=a.putAgent(stub,agent)
	if !ag{
		res := getRetString(1,"zpp invoke register put agent failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(agent)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}
//客服登录
//args：0-Id,1-Username,2-password,3-operation
func (a *AgentChaincode) login(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args)!=4{
		res := getRetString(1,"zpp invoke login args error")
		return shim.Error(res)
	}
	//根据Id获取客服身份信息
	agent,ag := a.getAgent(stub,args[0])
	if !ag{
		res := getRetString(1,"zpp invoke login get agent error,not exist")
		return shim.Error(res)
	}
	//验证账户状态(冻结或注销不可用)
	if agent.AccountState==AgentInfo_State_Frozen{
		res:= getRetString(1,"accountState is Frozen")
		return shim.Error(res)
	}
	if agent.AccountState==AgentInfo_State_Canceld{
		res:= getRetString(1,"accountState is Canceled")
		return shim.Error(res)
	}
	agent.AccountState=AgentInfo_State_Logining
	agent.Operation=args[3]
	//保存更改信息
	_,ag =a.putAgent(stub,agent)
	if !ag{
		res := getRetString(1,"zpp invoke login put agent error")
		return shim.Error(res)
	}

	jsonStr,err:=json.Marshal(agent)
	if err!=nil{
		log.Fatal(err)
	}
	res:= getRetByte(0,string(jsonStr))
	return shim.Success(res)
}
//客服登出
//args:0-agentId,1-operation
func (a *AgentChaincode) logout(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	//如果传入参数数量不对
	if len(args)!=2{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	//根据agentid获取客服身份信息
	agent,ag := a.getAgent(stub,args[0])
	if !ag{
		res := getRetString(1,"zpp invoke logout get agent error,not exist")
		return shim.Error(res)
	}

	//信息无误
	agent.AccountState=AgentInfo_State_Valid   //账户变成有效状态
	agent.Operation=args[1]
	//保存更改信息
	_,ag =a.putAgent(stub,agent)
	if !ag {
		res := getRetString(1, "zpp invoke logout put agent error")
		return shim.Error(res)
	}
	//返回成功信息

	jsonStr,err:=json.Marshal(agent)
	if err!=nil{
		log.Fatal(err)
	}
	res:= getRetByte(0,string(jsonStr))
	return shim.Success(res)
}
//客服冻结(违规操作之后冻结帐号)
//args:0-agentId
func (a *AgentChaincode) frozen(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	//参数数量不对
	if len(args)!=2{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	//根据agentid获取客服身份信息
	agent,ag := a.getAgent(stub,args[0])
	if !ag{
		res := getRetString(1,"zpp invoke frozen get agent error,not exist")
		return shim.Error(res)
	}
	//查询状态,非注销用户即可冻结
	if agent.AccountState==AgentInfo_State_Canceld{
		res := getRetString(1,"wrong accountState error")
		return shim.Error(res)
	}
	agent.AccountState=AgentInfo_State_Frozen
	agent.Operation=args[1]
	//保存更改信息
	_,ag =a.putAgent(stub,agent)
	if !ag{
		res := getRetString(1,"zpp invoke frozen put agent error")
		return shim.Error(res)
	}

	jsonStr,err:=json.Marshal(agent)
	if err!=nil{
		log.Fatal(err)
	}
	//返回成功信息
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}
//客服帐号解冻
//args:0-Id
func (a *AgentChaincode) unfrozen(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	//参数数量不对
	if len(args)!=2{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	//根据agentid获取客服身份信息
	agent,ag := a.getAgent(stub,args[0])
	if !ag{
		res := getRetString(1,"zpp invoke frozen get agent error,not exist")
		return shim.Error(res)
	}
	agent.AccountState=AgentInfo_State_Valid
	agent.Operation=args[1]
	//保存更改信息
	_,ag =a.putAgent(stub,agent)
	if !ag{
		res := getRetString(1,"zpp invoke frozen put agent error")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(agent)
	if err!=nil{
		log.Fatal(err)
	}
	//返回成功信息
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//客服帐号注销/删除
//args:0-Id
func (a *AgentChaincode) delete(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	//参数数量不对
	if len(args)!=2{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	//根据agentid获取客服身份信息
	agent,ag := a.getAgent(stub,args[0])
	if !ag{
		res := getRetString(1,"zpp invoke frozen get agent error,not exist")
		return shim.Error(res)
	}
	//查询状态,已删除用户不可再删除
	if agent.AccountState==AgentInfo_State_Canceld{
		res := getRetString(1,"wrong accountState error")
		return shim.Error(res)
	}
	agent.AccountState=AgentInfo_State_Canceld
	agent.Operation=args[1]
	//保存更改信息
	_,ag =a.putAgent(stub,agent)
	if !ag{
		res := getRetString(1,"zpp invoke delete put agent error")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(agent)
	if err!=nil{
		log.Fatal(err)
	}
	//返回成功信息
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//modify password
//args:0-ID,1-password
func (a *AgentChaincode) modifyPassword(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args)!=3{
		res:=getRetString(1,"zpp invoke args not correct")
		return shim.Error(res)
	}
	//根据agentid获取客服身份信息
	agent,ag := a.getAgent(stub,args[0])
	if !ag{
		res := getRetString(1,"zpp invoke frozen get agent error,not exist")
		return shim.Error(res)
	}
	agent.Password=args[1]
	agent.AccountState=AgentInfo_State_Valid
	agent.Operation=args[2]
	//保存更改信息
	_,ag =a.putAgent(stub,agent)
	if !ag{
		res := getRetString(1,"zpp invoke delete put agent error")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(agent)
	if err!=nil{
		log.Fatal(err)
	}
	//返回成功信息
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//根据id取得该账户历史
func (a *AgentChaincode) queryById(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=2{
		res:=getRetString(1,"zpp queryById args error")
		return shim.Error(res)
	}
	agent,ag:=a.getAgent(stub,args[0])
	if !ag{
		res := getRetString(1,"zpp queryById get agent error,not exist")
		return shim.Error(res)
	}
	resultsIterator, err := stub.GetHistoryForKey(Agent_Prefix+args[0])
	if err != nil {
		res := getRetString(1,"zpp queryById GetHistoryForKey error")
		return shim.Error(res)
	}
	defer resultsIterator.Close()
	var history []HistoryItem
	var hisAgent Agent
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
			var emptyAgent Agent
			hisItem.Agent=emptyAgent
		} else{
			json.Unmarshal(historyData.Value,&hisAgent)
			hisItem.Agent=hisAgent
		}
		history=append(history,hisItem)
	}
	agent.History=history
	agent.Operation=args[1]
	age,err:=json.Marshal(agent)
	if err != nil {
		res := getRetString(1,"zpp Marshal queryById agentList error")
		return shim.Error(res)
	}
	return shim.Success(age)
}



//chaincode invoke 接口
func (a *AgentChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
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
	if err := shim.Start(new(AgentChaincode)); err != nil {
		fmt.Printf("Error starting BillChaincode: %s", err)
	}
}
