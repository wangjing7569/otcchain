package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"time"
)

const (
	PriceInfo_State_NewPublish = "NewPublish"
	PriceInfo_State_UnderBid = "UnderBid"
	PriceInfo_State_WaitCal = "WaitCal"
	PriceInfo_State_Dealed = "Dealed"
)

const IndexName = "publishID~askNum"
const BidName1  = "provideId1~askNum"
const BidName2  = "provideId2~askNum"
const BidName3  = "provideId3~askNum"
const PublishIDDayTimeConstractTypeAskNum = "publishID~dayTime-constractType-askNum"

// 询报价结构
type Price struct {
	AskNum string `json:AskNum`                //询价委托号
	Biaodi string `json:Biaodi`          //标的
	CurrentPrice string `json:CurrentPrice`   //标的现价
	EndDate string `json:EndDate`      //到期日
	Angel string `json:Angel`     //方向
	ConstractType string `json:ConstractType`              //合约类型
	ExecuPrice string `json:ExecuPrice`      //执行价
	Amount string `json:Amount`      //合约规模（手）
	PublishID string `json:PublishID`      //发布人
	ProvideID1 string `json:ProvideID1`    //报价人1
	ProvideNum1 string `json:ProvideNum1`  //报价委托号1
	Price1 string `json:Price1`            //保证金报价1
	ProvideID2 string `json:ProvideID2`    //报价人2
	ProvideNum2 string `json:ProvideNum2`   //报价委托号2
	Price2 string `json:Price2`             //保证金报价2
	ProvideID3  string `json:ProvideID3`    //报价人3
	ProvideNum3 string `json:ProvideNUm3`    //报价委托号3
	Price3 string `json:Price3`               //保证金报价3
	DealerID string `json:DealerID`          //报价成交人
	State string `json:State`                          //状态
	History []HistoryItem `json:History`               //背书历史
}

// 背书历史item结构
type HistoryItem struct {
	TxId  string `json:"txId"`
	Price Price `json:"price"`
}

// chaincode response结构
type chaincodeRet struct {
	Code int // 0 success otherwise 1
	Des  string //description
}

type PriceChaincode struct {

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

// 根据询价委托号取出询价单信息
func (a *PriceChaincode) getPrice(stub shim.ChaincodeStubInterface, AskNum string) (Price, bool) {
	var price Price
	key := AskNum
	b,err := stub.GetState(key)
	if b==nil {
		return price, false
	}
	err = json.Unmarshal(b,&price)
	if err!=nil {
		return price, false
	}
	return price, true
}

// 保存询价单
func (a *PriceChaincode) putPrice(stub shim.ChaincodeStubInterface, price Price) ([]byte, bool) {

	byte,err := json.Marshal(price)
	if err!=nil {
		return nil, false
	}

	err = stub.PutState(price.AskNum, byte)
	if err!=nil {
		return nil, false
	}
	return byte, true
}

// chaincode Init 接口
func (a *PriceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


// 发布询价
// args: 0 - {Price Object}
func (a *PriceChaincode) publish(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"输入参数数目错误")
		return shim.Error(res)
	}

	var price Price
	err := json.Unmarshal([]byte(args[0]), &price)  //解析Object
	if err!=nil {
		res := getRetString(1,"参数解析失败")
		return shim.Error(res)
	}
	// TODO 根据询价委托号 查找是否该询价已存在
	// TODO 如已有同号询价委托号 返回error message
	_, existpr := a.getPrice(stub, price.AskNum)
	if existpr {
		res := getRetString(1,"该询价单已存在")
		return shim.Error(res)
	}
	//创建组合键发布人、合约类型、询价委托号与询价单的对应关系
	publishIdConstractTypeAskNumIndexKey, err := stub.CreateCompositeKey(PublishIDDayTimeConstractTypeAskNum, []string{price.PublishID,  price.ConstractType, price.AskNum})
	if err != nil {
		res := getRetString(1,"放入搜索表失败")
		return shim.Error(res)
	}
	stub.PutState(publishIdConstractTypeAskNumIndexKey, []byte(time.Now().Format("2019-08-29 12:56:56")))
	// 更改询价信息和状态并保存询价单:询价状态设为可报价
	price.State = PriceInfo_State_NewPublish
	// 保存当前状态
	_, pr := a.putPrice(stub, price)
	if !pr {
		res := getRetString(1,"询价发布失败")
		return shim.Error(res)
	}
	// 创建组合键发布人、询价委托号与询价单的对应关系
	publishIdAskNumIndexKey, err := stub.CreateCompositeKey(IndexName, []string{price.PublishID, price.AskNum})
	 if err != nil {
		res := getRetString(1,"放入搜索表失败")
		return shim.Error(res)
	}
	stub.PutState(publishIdAskNumIndexKey, []byte{0x00})

	res := getRetByte(0,"发布询价请求成功") //返回成功信息
	return shim.Success(res)
}

// 报价
//  args: 0 - AskNum ; 1 - ProvideID ; 2 - price; 3 - provideNum
func (a *PriceChaincode) bid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)<4 {
		res := getRetString(1,"参数数目至少等于3")
		return shim.Error(res)
	}
	// 根据询价委托号取得询价单
	price, pr := a.getPrice(stub, args[0])
	if !pr {
		res := getRetString(1,"获取询价信息失败")
		return shim.Error(res)
	}
	//判断当前询价单状态
	if price.State != PriceInfo_State_NewPublish && price.State != PriceInfo_State_UnderBid{
		res := getRetString(1,"当前询价单不是可以报价的状态")
		return shim.Error(res)
	}

	if price.ProvideID1 =="" {
		price.ProvideID1 = args[1]
		price.ProvideNum1 = args[3] //生成报价委托号
		price.Price1 = args[2]
		// 以ProvideID和AskNum构造复合key
		provideIDAskNum1, err := stub.CreateCompositeKey(BidName1, []string{price.ProvideID1, price.AskNum})
		if err != nil {
			res := getRetString(1,"构造失败")
			return shim.Error(res)
		}
		stub.PutState(provideIDAskNum1, []byte{0x00})
	} else if price.ProvideID2 == ""{
		price.ProvideID2 = args[1]
		price.ProvideNum2 = args[3]
		price.Price2 = args[2]
		provideIDAskNum2, err := stub.CreateCompositeKey(BidName2, []string{price.ProvideID2, price.AskNum})
		if err != nil {
			res := getRetString(1,"构造失败")
			return shim.Error(res)
		}
		stub.PutState(provideIDAskNum2, []byte{0x00})
	} else {
		price.ProvideID3 = args[1]
		price.ProvideNum3 = args[3]
		price.Price3 = args[2]
		provideIDAskNum3, err := stub.CreateCompositeKey(BidName3, []string{price.ProvideID3, price.AskNum})
		if err != nil {
			res := getRetString(1,"构造失败")
			return shim.Error(res)
		}
		stub.PutState(provideIDAskNum3, []byte{0x00})
	}
	price.State = PriceInfo_State_UnderBid //报价中
	//保存当前信息
	_, pr = a.putPrice(stub, price)
	if !pr {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
    //报价成功
	res := getRetByte(0,"报价请求成功")
	return shim.Success(res)
}

// 成交
// args: 0 - AskNum ; 1 - ProvideId
func (a *PriceChaincode) deal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)<2 {
		res := getRetString(1,"参数数目至少2个")
		return shim.Error(res)
	}
	// 根据票号取得票据
	price, pr := a.getPrice(stub, args[0])
	if !pr {
		res := getRetString(1,"获取询/报价信息失败")
		return shim.Error(res)
	}
	//判断当前状态
	if price.State != PriceInfo_State_UnderBid{
		res := getRetString(1,"当前询价单状态错误")
		return shim.Error(res)
	}
	//判断成交对象合法性
	if price.ProvideID1 != args[1] && price.ProvideID2 != args[1] && price.ProvideID3 != args[1]{
		res := getRetString(1,"当前输入报价人ID并非有效ID")
		return shim.Error(res)
	}
	price.DealerID = args[1]
    price.State = PriceInfo_State_Dealed //已成交
	// 保存当前信息
	_, pr = a.putPrice(stub, price)
	if !pr {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res := getRetByte(0,"成交请求成功")
	return shim.Success(res)
}


//撤回报价
// 0 - AskNum , 1 - provideId
func (a * PriceChaincode)  withdraw(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)<2{
		res := getRetString(1,"withdraw args <2")
		return shim.Error(res)
	}
	// 根据票号取得票据
	price, pr := a.getPrice(stub, args[0])
	if !pr {
		res := getRetString(1,"获取询/报价信息失败")
		return shim.Error(res)
	}
	//判断当前状态
	if price.State != PriceInfo_State_UnderBid{
		res := getRetString(1,"当前询价单状态错误")
		return shim.Error(res)
	}
	//判断成交对象合法性
	if price.ProvideID1 != args[1] && price.ProvideID2 != args[1] && price.ProvideID3 != args[1]{
		res := getRetString(1,"当前输入报价人ID并非有效ID")
		return shim.Error(res)
	}
	//撤销
	if price.ProvideID1 == args[1]{
		price.ProvideID1 = ""
		price.ProvideNum1 = ""
		price.Price1 = ""

	}else if price.ProvideID2 == args[1]{
		price.ProvideID2 = ""
		price.ProvideNum2 = ""
		price.Price2 = ""
	}else if price.ProvideID3 == args[3] {
		price.ProvideID3 = ""
		price.ProvideNum3 = ""
		price.Price3 = ""
	}
	if price.ProvideID1 == "" && price.ProvideID2 == "" && price.ProvideID3 == ""{
		price.State = PriceInfo_State_NewPublish
	}else{
		price.State = PriceInfo_State_UnderBid
	}
	// 保存当前信息
	_, pr = a.putPrice(stub, price)
	if !pr {
		res := getRetString(1,"当前信息保存失败")
		return shim.Error(res)
	}
	res := getRetByte(0,"成交请求成功")
	return shim.Success(res)
}

// 查询我的票据:根据持票人编号 批量查询票据
//  0 - publishId ;
func (a *PriceChaincode) queryMyAsk(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyAsk args!=1")
		return shim.Error(res)
	}
	// 以持票人ID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(IndexName, []string{args[0]})
	if err != nil {
		res := getRetString(1,"queryMyAsk get list error")
		return shim.Error(res)
	}
	defer asksIterator.Close()

	var priceList = []Price{}

	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		// 取得持票人名下的票号
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMyAsk SplitCompositeKey error")
			return shim.Error(res)
		}
		// 根据askNum取得票据
		price, pr := a.getPrice(stub, compositeKeyParts[1])
		if !pr {
			res := getRetString(1,"queryMyAsk get price error")
			return shim.Error(res)
		}
		priceList = append(priceList, price)
	}
	// 取得并返回数组
	p, err := json.Marshal(priceList)
	if err != nil {
		res := getRetString(1,"Marshal queryMyAsk priceList error")
		return shim.Error(res)
	}
	return shim.Success(p)
}


// 查询我的报价
//  0 - provideId ;
func (a *PriceChaincode) queryMyProvide(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryMyProvide args!=1")
		return shim.Error(res)
	}
	// 以provideID从search表中批量查询所持有的票号
	asksIterator, err := stub.GetStateByPartialCompositeKey(BidName1, []string{args[0]})
	//asksIterator2, err2 := stub.GetStateByPartialCompositeKey(BidName2, []string{args[0]})
	//asksIterator3, err3 := stub.GetStateByPartialCompositeKey(BidName3, []string{args[0]})

	if err != nil {
		res := getRetString(1,"queryMyAsk get list error")
		return shim.Error(res)
	}
    /*
	if asksIterator1 != nil{
		asksIterator := asksIterator1
	}else if asksIterator2 != nil{
		asksIterator := asksIterator2
	}else if asksIterator3 != nil{
		asksIterator := asksIterator3
	}*/
	defer asksIterator.Close()
	var priceList = []Price{}
	for asksIterator.HasNext() {
		kv, _ := asksIterator.Next()
		// 取得报价人名下的票号
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1," queryMyProvide SplitCompositeKey error")
			return shim.Error(res)
		}
		// 根据askNum取得票据
		price, pr := a.getPrice(stub, compositeKeyParts[1])
		if !pr {
			res := getRetString(1,"queryMyProvide get price error")
			return shim.Error(res)
		}
		priceList = append(priceList, price)
	}
	// 取得并返回数组
	p, err := json.Marshal(priceList)
	if err != nil {
		res := getRetString(1,"Marshal queryMyProvide priceList error")
		return shim.Error(res)
	}
	return shim.Success(p)
}


func (a *PriceChaincode) queryByAskNum(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"queryByAskNum args!=1")
		return shim.Error(res)
	}
	// 取得该票据
	price, bl := a.getPrice(stub, args[0])
	if !bl {
		res := getRetString(1,"queryByAskNum get ask error")
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
	var hisPrice Price
	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			res := getRetString(1,"queryByAskNum resultsIterator.Next() error")
			return shim.Error(res)
		}

		var hisItem HistoryItem
		hisItem.TxId = historyData.TxId //copy transaction id over
		json.Unmarshal(historyData.Value, &hisPrice) //un stringify it aka JSON.parse()
		if historyData.Value == nil {              //bill has been deleted
			var emptyPrice Price
			hisItem.Price = emptyPrice //copy nil marble
		} else {
			json.Unmarshal(historyData.Value, &hisPrice) //un stringify it aka JSON.parse()
			hisItem.Price = hisPrice                          //copy bill over
		}
		history = append(history, hisItem) //add this tx to the list
	}
	// 将历史记录做为票据的一个属性 一同返回
	price.History = history

	b, err := json.Marshal(price)
	if err != nil {
		res := getRetString(1,"queryByAskNum getList error")
		return shim.Error(res)
	}
	return shim.Success(b)
}

// chaincode Invoke 接口
func (a *PriceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function,args := stub.GetFunctionAndParameters()

	// invoke
	if function == "publish" {
		return a.publish(stub, args)
	}else if function == "bid" {
		return  a.bid(stub, args)
	}else if function == "deal" {
		return  a.deal(stub, args)
	}else if function == "withdraw" {
		return  a.withdraw(stub, args)
	}else if function == "queryMyAsk" {
		return  a.queryMyAsk(stub, args)
	}else if function == "queryMyProvide" {
		return  a.queryMyProvide(stub, args)
	}else if function == "queryByAskNum" {
		return  a.queryByAskNum(stub, args)
	}
	
	res := getRetString(1,"Chaincode Unknown method!")
	//chaincodeLogger.Info("%s",res)
	//chaincodeLogger.Infof("%s",res)
	return shim.Error(res)
}

func main() {
	if err := shim.Start(new(PriceChaincode)); err != nil {
		fmt.Printf("Error starting PriceChaincode: %s", err)
	}
}
