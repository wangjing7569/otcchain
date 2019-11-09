package main


import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"log"
)

// 票据状态
const (
	//票据自由竞价发布后如果有人报价则不能重新发布
	BillInfo_State_NewPublish = "NewPublish"       //新发布(未报价)
	BillInfo_State_WaitPay = "WaitPay"             //等待买家付款
	BillInfo_State_UnderBid = "UnderBid"          //竞价中,已有人报价
	BillInfo_State_PayWaitEnsure ="PayWaitEnsure"   //待确认收款
	BillInfo_State_WaitEndorse = "WaitEndorse"    //等待背书
	BillInfo_State_EndrWaitSign = "EndrWaitSign"    //已背书待签收
	BillInfo_State_Dealed = "Dealed"           //交易完成
	BillInfo_State_Unshelved = "Unshelved"       //票据下架(需求截止时间已到)
)
// 票据key的前缀
const Bill_Prefix = "Bill_"

// search表的映射名
const IndexName = "holderName~billNo"


const HolderIdDayTimeBillTypeBillNoIndexName = "holderId~dayTime-billType-billNo"

//票据信息
type Bill struct {
	BillId string `json:"billId"`                //票据id 新增时为0
	BillType string `json:"billType"`            //票据类型 1.电子商票
	ImgOnUrl string `json:"imgOnUrl"`              //票据正面(单)
	ImgCodeOnUrl string `json:"imgCodeOnUrl"`              //票据打码正面(单)
	ImgCodeBackUrlList string `json:"imgCodeBackUrlList"`          //票据反面(多) 数组
	ImgBackUrlList string `json:"imgBackUrlList"`
	BillNo string `json:"billNo"`                //票据号
	AcceptorName string `json:"acceptorName"`      //承兑人全称
	BearerName string `json:"bearerName"`         //持票人全称
	UserId string `json:"userId"`             // 当前持票人用户Id
	Amount string `json:"amount"`                 //票面金额
	DrawDate string `json:"drawDate"`             //出票日期
	EndDate string `json:"endDate"`               //到期日期
	FlawState string `json:"flawState"`           //票据瑕疵  多选项
	EndorseTimes string `json:"endorseTimes"`       //背书次数
	TenderType string `json:"tenderType"`            //报价方式  1.一口价  2.自由竞价
	DeductMoney string `json:"deductMoney"`      //每十万扣款 一口价必填 自由竞价为空
	PassTime string `json:"passTime"`            //发布日期
	RequirementEndDate string `json:"requirementEndDate"`    //需求截止时间
	BidPersonType1Id string `json:"bidPersonType1Id"`         //一口价报价人Id
	BidPersonType2Id string `json:"bidPersonType2Id"`            //自由竞价报价人Id
	EnsurePersonId string `json:"ensurePersonId"`              //成交对象用户id
	PaymentVoucherUrl string `json:"paymentVoucherUrl"`         //付款凭证
	EndorsementVoucherUrl string `json:"endorsementVoucherUrl"`    //背书凭证
	Operation string `json:"operation"`
	State string `json:"State"`                        //票据状态
	History []HistoryItem `json:"History"`               //背书历史
}
// 背书历史item结构
type HistoryItem struct {
	TxId  string `json:"txId"`
	Bill Bill `json:"bill"`
}

// chaincode response结构
type chaincodeRet struct {
	Code int // 0 success otherwise 1
	Des  string //description
}

// chaincode
type BillChaincode struct {
}

// response message format
func getRetByte(code int,des string) []byte {
	var r chaincodeRet
	r.Code = code
	r.Des = des

	b,err := json.Marshal(r)

	if err!=nil {
		fmt.Println("marshal Ret failed")
		return nil
	}
	return b
}

// response message format
func getRetString(code int,des string) string {
	var r chaincodeRet
	r.Code = code
	r.Des = des

	b,err := json.Marshal(r)

	if err!=nil {
		fmt.Println("marshal Ret failed")
		return ""
	}
	//chaincodeLogger.Infof("%s",string(b[:]))
	return string(b[:])
}

// 以票号id为key获取票据
func (a *BillChaincode) getBill(stub shim.ChaincodeStubInterface,billId string) (Bill, bool) {
	var bill Bill
	key := Bill_Prefix + billId
	b,err := stub.GetState(key)
	if b==nil {
		return bill, false
	}
	err = json.Unmarshal(b,&bill)
	if err!=nil {
		return bill, false
	}
	return bill, true
}

// 保存票据
func (a *BillChaincode) putBill(stub shim.ChaincodeStubInterface, bill Bill) ([]byte, bool) {

	byte,err := json.Marshal(bill)
	if err!=nil {
		return nil, false
	}

	err = stub.PutState(Bill_Prefix + bill.BillId, byte)
	if err!=nil {
		return nil, false
	}
	return byte, true
}

// chaincode Init 接口
func (a *BillChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}
// 票据发布
// args: 0 - {Bill Object}
func (a *BillChaincode) issue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"zpp Invoke issue args!=1")
		return shim.Error(res)
	}
	var bill Bill
	err:=json.Unmarshal([]byte(args[0]),&bill)
	if err!=nil{
		res:=getRetString(1,"zpp Invoke issue unmarshal failed")
		return shim.Error(res)
	}
	// 更改票据信息和状态并保存票据:票据状态设为新发布
	bill.State = BillInfo_State_NewPublish
	// 保存票据
	_, bl := a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke issue put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//一口价票据买方报价
//args：token,0-billId,1-userId
func (a *BillChaincode) bidType1(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=3{
		res := getRetString(1,"zpp invoke bid args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res:=getRetString(1,"zpp invoke bid get bill error")
		return shim.Error(res)
	}
	bill.BidPersonType1Id = args[1]
	bill.EnsurePersonId=args[1]
	bill.Operation=args[2]
	bill.State = BillInfo_State_WaitPay
	// 保存票据
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke endorse put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)

}

//自由竞价票据买方报价
//args：0-billId,1-userId,2-deductMoney
func (a *BillChaincode) bidType2(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=4{
		res := getRetString(1,"zpp invoke bid args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res:=getRetString(1,"zpp invoke bid get bill error")
		return shim.Error(res)
	}

	//判断是否是第一个报价
    if bill.BidPersonType2Id==""{
    	bill.BidPersonType2Id=args[1]
	}else {
		bill.BidPersonType2Id= bill.BidPersonType2Id+","+args[1]
	}
	if bill.DeductMoney==""{
		bill.DeductMoney=args[2]
	}else {
		bill.DeductMoney=bill.DeductMoney+","+args[2]
	}
	bill.State = BillInfo_State_UnderBid
	bill.Operation=args[3]
	// 保存票据
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke endorse put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//确认成交对象
//args:0-billId,1-userId
func (a *BillChaincode) ensure(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=3{
		res := getRetString(1,"zpp invoke ensure args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res := getRetString(1,"zpp invoke get bill error")
		return shim.Error(res)
	}
	//判断是否一口价报价
	if bill.TenderType=="1"{
		bill.BidPersonType1Id=args[1]
	}
	
	bill.EnsurePersonId=args[1]
	bill.Operation=args[2]
	bill.State=BillInfo_State_WaitPay
	//保存修改信息
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke endorse put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}
//买方付款
//args：0-billId,1-userId,2-paymentVoucherUrl
func (a *BillChaincode) pay(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=4{
		res := getRetString(1,"zpp invoke ensure args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res := getRetString(1,"zpp invoke get bill error")
		return shim.Error(res)
	}

	if bill.EnsurePersonId==""&&bill.BidPersonType1Id==args[1]{
		bill.EnsurePersonId=args[1]
	}

	bill.PaymentVoucherUrl=args[2]
	bill.Operation=args[3]
	bill.State=BillInfo_State_PayWaitEnsure
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke pay put bill failed")
		return shim.Error(res)
	}

	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//卖方确认收款
//args：0-billId,1-userId
func (a *BillChaincode) getMoney(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args)!=3{
		res := getRetString(1,"zpp invoke ensure args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res := getRetString(1,"zpp invoke get bill error")
		return shim.Error(res)
	}
	if bill.UserId!=args[1]{
		res:=getRetString(1,"卖方身份不正确")
		return shim.Error(res)
	}
	bill.Operation=args[2]
	bill.State=BillInfo_State_WaitEndorse
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke getMoney put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//卖方背书
//args:0-billId,1-userId,2-endorsementVoucherUrl
func (a *BillChaincode) endorse(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=4{
		res := getRetString(1,"zpp invoke ensure args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res := getRetString(1,"zpp invoke get bill error")
		return shim.Error(res)
	}

	//验证信息
	if bill.UserId!=args[1]{
		res:=getRetString(1,"卖方身份不正确")
		return shim.Error(res)
	}
	bill.EndorsementVoucherUrl=args[2]
	bill.Operation=args[3]
	bill.State=BillInfo_State_EndrWaitSign
	//保存修改信息
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke endorse put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

func add(x int, y string) string {

	i, err := strconv.Atoi(y)
	if err != nil {
		panic(err)
	}
	i = i + x
	var res string
	res=strconv.Itoa(i)
	return res
}

//票据签收
//args:0-billId,1-userId
func (a *BillChaincode) sign (stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=3{
		res := getRetString(1,"zpp invoke ensure args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res := getRetString(1,"zpp invoke get bill error")
		return shim.Error(res)
	}

	bill.State=BillInfo_State_Dealed
	bill.UserId=args[1]
	bill.Operation=args[2]
	bill.EndorseTimes=add(1,bill.EndorseTimes)
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke sign put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//票据下架
//args：0-billId,1-userId
func (a *BillChaincode) unShelve (stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=3{
		res := getRetString(1,"zpp invoke ensure args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res := getRetString(1,"zpp invoke get bill error")
		return shim.Error(res)
	}
	bill.Operation=args[2]
	bill.State=BillInfo_State_Unshelved
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke unShelve put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}

//票据重新发布
//args:0-billId,1-tenderType,2-deductMoney,3-requirementEndDate
func (a *BillChaincode) pushBillAgain(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args)!=5 {
		res := getRetString(1," args error")
		return shim.Error(res)
	}
	bill,bl:=a.getBill(stub,args[0])
	if !bl{
		res := getRetString(1,"zpp invoke get bill error")
		return shim.Error(res)
	}
	bill.TenderType=args[1]
	bill.DeductMoney=args[2]
	bill.RequirementEndDate=args[3]
	bill.Operation=args[4]
	bill.State=BillInfo_State_NewPublish
	_, bl = a.putBill(stub, bill)
	if !bl {
		res := getRetString(1,"zpp Invoke put bill failed")
		return shim.Error(res)
	}
	jsonStr,err:=json.Marshal(bill)
	if err!=nil{
		log.Fatal(err)
	}
	res:=getRetByte(0,string(jsonStr))
	return shim.Success(res)
}


// 根据票号取得票据 以及该票据背书历史
//  0 - Bill_Id ;
func (a *BillChaincode) queryByBillId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=2 {
		res := getRetString(1,"queryByBillId args error")
		return shim.Error(res)
	}
	// 取得该票据
	bill, bl := a.getBill(stub, args[0])
	if !bl {
		res := getRetString(1,"queryByBillId get bill error")
		return shim.Error(res)
	}

	// 取得背书历史: 通过fabric api取得该票据的变更历史
	resultsIterator, err := stub.GetHistoryForKey(Bill_Prefix+args[0])
	if err != nil {
		res := getRetString(1,"queryByBillId GetHistoryForKey error")
		return shim.Error(res)
	}
	defer resultsIterator.Close()
	var history []HistoryItem
	var hisBill Bill
	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			res := getRetString(1,"queryByBillId resultsIterator.Next() error")
			return shim.Error(res)
		}

		var hisItem HistoryItem
		hisItem.TxId = historyData.TxId //copy transaction id over
		json.Unmarshal(historyData.Value, &hisBill) //un stringify it aka JSON.parse()
		if historyData.Value == nil {              //bill has been deleted
			var emptyBill Bill
			hisItem.Bill = emptyBill //copy nil marble
		} else {
			json.Unmarshal(historyData.Value, &hisBill) //un stringify it aka JSON.parse()
			hisItem.Bill = hisBill                          //copy bill over
		}
		history = append(history, hisItem) //add this tx to the list
	}
	// 将背书历史做为票据的一个属性 一同返回
	bill.History = history
	bill.Operation=args[1]

	b, err := json.Marshal(bill)
	if err != nil {
		res := getRetString(1,"Marshal queryByBillId billList error")
		return shim.Error(res)
	}
	return shim.Success(b)
}
// chaincode Invoke 接口
func (a *BillChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function,args := stub.GetFunctionAndParameters()
	//chaincodeLogger.Info("%s%s","zpp function=",function)
	//chaincodeLogger.Info("%s%s","zpp args=",args)
	// invoke
	if function == "issue" {
		return a.issue(stub, args)
	}else if function=="bidType1"{
		return a.bidType1(stub,args)
	}else if function=="bidType2"{
		return a.bidType2(stub,args)
	}else if function=="ensure"{
		return a.ensure(stub,args)
	}else if function=="pay"{
		return a.pay(stub,args)
	}else if function=="getMoney"{
		return a.getMoney(stub,args)
	}else if function=="endorse"{
		return a.endorse(stub,args)
	}else if function=="sign"{
		return a.sign(stub,args)
	}else if function=="unShelve"{
		return a.unShelve(stub,args)
	}else if function == "queryByBillId" {
		return a.queryByBillId(stub, args)
	}else if function == "pushBillAgain" {
		return a.pushBillAgain(stub, args)
	}
	res := getRetString(1,"zpp Unknown method!")
	//chaincodeLogger.Info("%s",res)
	//chaincodeLogger.Infof("%s",res)
	return shim.Error(res)
}

func main() {
	if err := shim.Start(new(BillChaincode)); err != nil {
		fmt.Printf("Error starting BillChaincode: %s", err)
	}
}