package main

import (
"fmt"
"github.com/hyperledger/fabric/core/chaincode/shim"
pb "github.com/hyperledger/fabric/protos/peer"
"encoding/json"
)

const (
	ClearInfo_State_New = "NewPublish"   //新成立订单
	ClearInfo_State_WaitBid = "WaitBid"    //等报价
	ClearInfo_State_UnderBid = "UnderBid"    //报价中
	ClearInfo_State_BidAccept = "Accept"    //接受报价
	ClearInfo_State_Hang = "Hang"      //挂起
	ClearInfo_State_Complete = "Complete"   //清算完成
	ClearInfo_State_HalfComplete = "HalfComplete"  //部分清算完成
)

const SaleID = "salesParty~orderID"
const BuyID = "buyParty~orderID"



//清算结构
type Clear struct{
	OrderID string `json:"OrderID"`  //交易订单号
	SaleParty string `json:"SaleParty"`  //卖方
	BuyParty string `json:"BuyParty"`    //买方
	Underlying string `json:Underlying`          //标的
	ConstractType string `json:ConstractType`              //合约类型
	OptionStyle string `json:OptionStyle`                  //期权类型
	StrikePrice string `json:StrikePrice`      //行权价
	QuotePrice string `json:"QuotePrice"`          //成交时的保证金报价
	AccPrice string `json:"AccPrice"`            //接受价
	Angel string `json:"Angel"`            //方向
	ConstractSize string `json:ConstractSize`      //合约规模（手）
	ExpiringDate string `json:ExpiringDate`      //到期日
	Request string `json:"Request"`       //当前请求
	Percent string `json:"Percent"`           //金额阶段成功比例
	BidPrice string `json:"BidPrice"`       //报价 
	State string `json:State`                          //状态
	History []HistoryItem `json:History`               //状态变更历史
}



// 背书历史item结构
type HistoryItem struct {
	TxId  string `json:"txId"`
	Clear Clear `json:"clear"`
}

// chaincode response结构
type chaincodeRet struct {
	Code int // 0 success otherwise 1
	Des  string //description
}

type ClearChaincode struct {

}

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

// 根据订单号取出询价单信息
func (a *ClearChaincode) getClear(stub shim.ChaincodeStubInterface, OrderID string) (Clear, bool) {
	var clear Clear
	key := OrderID
	b,err := stub.GetState(key)
	if b==nil {
		return clear, false
	}
	err = json.Unmarshal(b,&clear)
	if err!=nil {
		return clear, false
	}
	return clear, true
}

// 保存询价单
func (a *ClearChaincode) putClear(stub shim.ChaincodeStubInterface, clear Clear) ([]byte, bool) {

	byte,err := json.Marshal(clear)
	if err!=nil {
		return nil, false
	}

	err = stub.PutState(clear.OrderID, byte)
	if err!=nil {
		return nil, false
	}
	return byte, true
}

// chaincode Init 接口
func (a *ClearChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


// 发布订单
// args: 0 - {Price Object}
func (a *ClearChaincode) publish(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"输入参数数目错误")
		return shim.Error(res)
	}
	var clear Clear
	err := json.Unmarshal([]byte(args[0]), &clear)  //解析Object
	if err!=nil {
		res := getRetString(1,"参数解析失败")
		return shim.Error(res)
	}

	_, existpr := a.getClear(stub, clear.OrderID)
	if existpr {
		res := getRetString(1,"该订单已存在")
		return shim.Error(res)
	}
	clear.State = ClearInfo_State_New
	// 保存当前状态
	_, pr := a.putClear(stub, clear)
	if !pr {
		res := getRetString(1,"订单发布失败")
		return shim.Error(res)
	}
	
	saleIdKey, err := stub.CreateCompositeKey(SaleID, []string{clear.SaleParty, clear.OrderID})
	if err != nil {
		res := getRetString(1,"放入搜索表失败")
		return shim.Error(res)
	}
	stub.PutState(saleIdKey, []byte{0x00})

	//
	buyIdKey, err := stub.CreateCompositeKey(BuyID, []string{clear.BuyParty, clear.OrderID})
	if err != nil {
		res := getRetString(1,"放入搜索表失败BUY")
		return shim.Error(res)
	}
	stub.PutState(buyIdKey, []byte{0x00})

	res := getRetByte(0,"订单生成成功") //返回成功信息
	return shim.Success(res)
}

//部分清算成功 (订单号和比例, 订单请求)
func (a *ClearChaincode) expirationHalf(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)< 2{
		res:= getRetString(1,"参数数目至少等于1")
		return shim.Error(res)
	}
	//根据订单号获取信息
	clear, cl:=a.getClear(stub,args[0])
	if !cl{
		res:=getRetString(1,"获取订单信息失败")
		return shim.Error(res)
	}

	clear.Percent = args[1]
	clear.Request = args[2]
	clear.State = ClearInfo_State_HalfComplete
	_, cl = a.putClear(stub, clear)    //保存
	if !cl {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res:= getRetByte(0,"部分成功清算")
	return shim.Success(res)
}


//挂起
func (a *ClearChaincode) hang(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)< 2{
		res:= getRetString(1,"参数数目至少等于1")
		return shim.Error(res)
	}
	//根据订单号获取信息
	clear, cl:=a.getClear(stub,args[0])
	if !cl{
		res:=getRetString(1,"获取订单信息失败")
		return shim.Error(res)
	}
	clear.Request = args[1]
	clear.State = ClearInfo_State_Hang
	_, cl = a.putClear(stub, clear)    //保存
	if !cl {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res:= getRetByte(0,"订单挂起")
	return shim.Success(res)
}


//到期行权
// args: 0 - OrderID
func (a *ClearChaincode) expiration(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)< 1{
		res:= getRetString(1,"参数数目至少等于1")
		return shim.Error(res)
	}
	//根据订单号获取信息
	clear, cl:=a.getClear(stub,args[0])
	if !cl{
		res:=getRetString(1,"获取订单信息失败")
		return shim.Error(res)
	}
	clear.State = ClearInfo_State_Complete
	_, cl = a.putClear(stub, clear)    //保存
	if !cl {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res:= getRetByte(0,"成功清算")
	return shim.Success(res)
}

//提前平仓发起
//args: 0 - OrderID
func (a *ClearChaincode) askForOffset(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)< 1{
		res:= getRetString(1,"参数数目至少等于1")
		return shim.Error(res)
	}
	//根据订单号获取信息
	clear, cl:=a.getClear(stub,args[0])
	if !cl{
		res:=getRetString(1,"获取订单信息失败")
		return shim.Error(res)
	}
	//判断当前订单状态
	if clear.State != ClearInfo_State_New{
		res:=getRetString(1,"当前订单不允许发起提前平仓")
		return shim.Error(res)
	}
	//订单转换为等待报价状态
	clear.State=ClearInfo_State_WaitBid
	clear.Request = "请求平仓"
	_, cl = a.putClear(stub, clear)    //保存
	if !cl {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res:= getRetByte(0,"成功发起平仓请求")
	return shim.Success(res)
}


//报价
func (a *ClearChaincode) bid(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)< 2{
		res:= getRetString(1,"参数数目至少等于1")
		return shim.Error(res)
	}
	//根据订单号获取信息
	clear, cl:=a.getClear(stub,args[0])
	if !cl{
		res:=getRetString(1,"获取订单信息失败")
		return shim.Error(res)
	}
	//判断当前订单状态
	if clear.State != ClearInfo_State_WaitBid{
		res:=getRetString(1,"当前订单不允许报价")
		return shim.Error(res)
	}
	//订单转换为等待报价状态
	clear.State=ClearInfo_State_UnderBid
	clear.BidPrice = args[1]
	clear.Request = "对手报价"
	_, cl = a.putClear(stub, clear)    //保存
	if !cl {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res:= getRetByte(0,"成功发起平仓请求")
	return shim.Success(res)
}


//接受报价
func (a *ClearChaincode) accept(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)< 1{
		res:= getRetString(1,"参数数目至少等于1")
		return shim.Error(res)
	}
	//根据订单号获取信息
	clear, cl:=a.getClear(stub,args[0])
	if !cl{
		res:=getRetString(1,"获取订单信息失败")
		return shim.Error(res)
	}
	//判断当前订单状态
	if clear.State != ClearInfo_State_UnderBid{
		res:=getRetString(1,"当前订单不允许报价")
		return shim.Error(res)
	}
	//订单转换为报价成功状态
	clear.State=ClearInfo_State_BidAccept
	clear.AccPrice = clear.BidPrice
	
	_, cl = a.putClear(stub, clear)    //保存
	if !cl {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res:= getRetByte(0,"接受平仓")
	return shim.Success(res)
}



func (a *ClearChaincode) reject(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)< 1{
		res:= getRetString(1,"参数数目至少等于1")
		return shim.Error(res)
	}
	//根据订单号获取信息
	clear, cl:=a.getClear(stub,args[0])
	if !cl{
		res:=getRetString(1,"获取订单信息失败")
		return shim.Error(res)
	}
	//判断当前订单状态
	if clear.State != ClearInfo_State_UnderBid{
		res:=getRetString(1,"当前订单不允许拒绝")
		return shim.Error(res)
	}
	//订单转换为初始状态
	clear.State=ClearInfo_State_New
	
	_, cl = a.putClear(stub, clear)    //保存
	if !cl {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res:= getRetByte(0,"拒绝平仓")
	return shim.Success(res)
}






// 查询我的票据:根据持票人编号 批量查询票据
//  0 - publishId ;
func (a *ClearChaincode) queryMySale(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyAsk args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(SaleID, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMySale get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var clearList = []Clear{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		// 取得持票人名下的票号
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMySale SplitCompositeKey error")
			return shim.Error(res)
		}
		// 根据票号取得票据
		clear, bl := a.getClear(stub, compositeKeyParts[1])
		if !bl {
			res := getRetString(1,"queryMySale get order error")
			return shim.Error(res)
		}
		if (clear.State == ClearInfo_State_New || clear.State == ClearInfo_State_WaitBid || clear.State == ClearInfo_State_UnderBid){
		    clearList = append(clearList, clear)
		}
		
	}
	// 取得并返回票据数组
	b, err := json.Marshal(clearList)
	if err != nil {
		res := getRetString(1,"Marshal queryMySale List error")
		return shim.Error(res)
	}
	return shim.Success(b)
}



func (a *ClearChaincode) queryMyBuy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyBuy args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(BuyID, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMyBuy get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var clearList = []Clear{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMyBuy SplitCompositeKey error")
			return shim.Error(res)
		}
		// 
		clear, bl := a.getClear(stub, compositeKeyParts[1])
		if !bl {
			res := getRetString(1,"queryMyBuy get order error")
			return shim.Error(res)
		}
		if (clear.State == ClearInfo_State_New || clear.State == ClearInfo_State_WaitBid || clear.State == ClearInfo_State_UnderBid){
		    clearList = append(clearList, clear)
		}
	}
	// 取得并返回票据数组
	b, err := json.Marshal(clearList)
	if err != nil {
		res := getRetString(1,"Marshal queryMyBuy List error")
		return shim.Error(res)
	}
	return shim.Success(b)
}



func (a *ClearChaincode) queryHis1(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyHis args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(BuyID, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMyHis get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var clearList = []Clear{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMyBuy SplitCompositeKey error")
			return shim.Error(res)
		}
		// 
		clear, bl := a.getClear(stub, compositeKeyParts[1])
		if !bl {
			res := getRetString(1,"queryMyBuy get order error")
			return shim.Error(res)
		}
		if (clear.State == ClearInfo_State_Complete){
		    clearList = append(clearList, clear)
		}
	}
	// 取得并返回票据数组
	b, err := json.Marshal(clearList)
	if err != nil {
		res := getRetString(1,"Marshal queryHis List error")
		return shim.Error(res)
	}
	return shim.Success(b)
}



func (a *ClearChaincode) queryHis2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyAsk args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(SaleID, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMySale get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var clearList = []Clear{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		// 取得持票人名下的票号
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMySale SplitCompositeKey error")
			return shim.Error(res)
		}
		// 根据票号取得票据
		clear, bl := a.getClear(stub, compositeKeyParts[1])
		if !bl {
			res := getRetString(1,"queryMySale get order error")
			return shim.Error(res)
		}
		if (clear.State == ClearInfo_State_Complete){
		    clearList = append(clearList, clear)
		}
		
	}
	// 取得并返回票据数组
	b, err := json.Marshal(clearList)
	if err != nil {
		res := getRetString(1,"Marshal queryMySale List error")
		return shim.Error(res)
	}
	return shim.Success(b)
}


func (a *ClearChaincode) queryHalf1(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyHis args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(BuyID, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMyHis get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var clearList = []Clear{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMyBuy SplitCompositeKey error")
			return shim.Error(res)
		}
		// 
		clear, bl := a.getClear(stub, compositeKeyParts[1])
		if !bl {
			res := getRetString(1,"queryMyBuy get order error")
			return shim.Error(res)
		}
		if (clear.State == ClearInfo_State_HalfComplete){
		    clearList = append(clearList, clear)
		}
	}
	// 取得并返回票据数组
	b, err := json.Marshal(clearList)
	if err != nil {
		res := getRetString(1,"Marshal queryHis List error")
		return shim.Error(res)
	}
	return shim.Success(b)
}


func (a *ClearChaincode) queryHalf2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyAsk args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(SaleID, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMySale get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var clearList = []Clear{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		// 取得持票人名下的票号
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMySale SplitCompositeKey error")
			return shim.Error(res)
		}
		// 根据票号取得票据
		clear, bl := a.getClear(stub, compositeKeyParts[1])
		if !bl {
			res := getRetString(1,"queryMySale get order error")
			return shim.Error(res)
		}
		if (clear.State == ClearInfo_State_HalfComplete){
		    clearList = append(clearList, clear)
		}
		
	}
	// 取得并返回票据数组
	b, err := json.Marshal(clearList)
	if err != nil {
		res := getRetString(1,"Marshal queryMySale List error")
		return shim.Error(res)
	}
	return shim.Success(b)
}



func (a *ClearChaincode) queryHang1(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyHis args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(BuyID, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMyHis get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var clearList = []Clear{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMyBuy SplitCompositeKey error")
			return shim.Error(res)
		}
		// 
		clear, bl := a.getClear(stub, compositeKeyParts[1])
		if !bl {
			res := getRetString(1,"queryMyBuy get order error")
			return shim.Error(res)
		}
		if (clear.State == ClearInfo_State_HalfComplete){
		    clearList = append(clearList, clear)
		}
	}
	// 取得并返回票据数组
	b, err := json.Marshal(clearList)
	if err != nil {
		res := getRetString(1,"Marshal queryHis List error")
		return shim.Error(res)
	}
	return shim.Success(b)
}





func (a *ClearChaincode) queryHang2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyAsk args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(SaleID, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMySale get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var clearList = []Clear{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		// 取得持票人名下的票号
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMySale SplitCompositeKey error")
			return shim.Error(res)
		}
		// 根据票号取得票据
		clear, bl := a.getClear(stub, compositeKeyParts[1])
		if !bl {
			res := getRetString(1,"queryMySale get order error")
			return shim.Error(res)
		}
		if (clear.State == ClearInfo_State_Hang){
		    clearList = append(clearList, clear)
		}
		
	}
	// 取得并返回票据数组
	b, err := json.Marshal(clearList)
	if err != nil {
		res := getRetString(1,"Marshal queryMySale List error")
		return shim.Error(res)
	}
	return shim.Success(b)
}


func (a *ClearChaincode) queryByOrderID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1," args!=1")
		return shim.Error(res)
	}
	// 取得该票据
	clear, bl := a.getClear(stub, args[0])
	if !bl {
		res := getRetString(1," get ask error")
		return shim.Error(res)
	}

	// 取得背书历史: 通过fabric api取得该票据的变更历史
	resultsIterator, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		res := getRetString(1,"queryByAskNum GetHistoryForKey error")
		return shim.Error(res)
	}
	defer resultsIterator.Close()

	var history []HistoryItem
	var hisPrice Clear
	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			res := getRetString(1," resultsIterator.Next() error")
			return shim.Error(res)
		}

		var hisItem HistoryItem
		hisItem.TxId = historyData.TxId //copy transaction id over
		json.Unmarshal(historyData.Value, &hisPrice) //un stringify it aka JSON.parse()
		if historyData.Value == nil {              //bill has been deleted
			var emptyPrice Clear
			hisItem.Clear = emptyPrice //copy nil marble
		} else {
			json.Unmarshal(historyData.Value, &hisPrice) //un stringify it aka JSON.parse()
			hisItem.Clear = hisPrice                          //copy bill over
		}
		history = append(history, hisItem) //add this tx to the list
	}
	// 将历史记录做为票据的一个属性 一同返回
	clear.History = history

	b, err := json.Marshal(clear)
	if err != nil {
		res := getRetString(1,"queryByAskNum getList error")
		return shim.Error(res)
	}
	return shim.Success(b)
}

// chaincode Invoke 接口
func (a *ClearChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function,args := stub.GetFunctionAndParameters()

	// invoke expiration
	if function == "publish" {
		return a.publish(stub, args)
	}else if function == "queryMySale" {
		return  a.queryMySale(stub, args)
	}else if function == "queryByOrderID" {
		return  a.queryByOrderID(stub, args)
	}else if function == "queryMyBuy" {
		return  a.queryMyBuy(stub, args)
	}else if function == "expiration" {
		return  a.expiration(stub, args)
	}else if function == "askForOffset" {
		return  a.askForOffset(stub, args)
	}else if function == "expirationHalf" {
		return  a.expirationHalf(stub, args)
	}else if function == "bid"{
		return a.bid(stub,args)
	}else if function =="accept"{
		return a.accept(stub,args)
	}else if function == "reject"{
		return a.reject(stub, args)
	}else if function == "queryHis1"{
		return a.queryHis1(stub, args)
	}else if function == "queryHis2"{
		return a.queryHis2(stub, args)
	}else if function == "queryHalf1"{
		return a.queryHalf1(stub, args)
	}else if  function == "queryHalf2" {
		return a.queryHalf2(stub, args)
	}else if function == "queryHang1" {
		return a.queryHang1(stub, args)
	}else if function == "queryHang2" {
		return a.queryHang2(stub, args)
	}else if function == "hang"{
		return a.hang(stub, args)
	}

	res := getRetString(1,"Chaincode Unknown method!")
	return shim.Error(res)
}

func main() {
	if err := shim.Start(new(ClearChaincode)); err != nil {
		fmt.Printf("Error starting PriceChaincode: %s", err)
	}
}

