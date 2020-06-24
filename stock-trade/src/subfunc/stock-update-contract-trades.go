package subfunc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	// "net/http"
	"os"
	"time"

	// for gorm
	_ "github.com/go-sql-driver/mysql"

	"github.com/aws/aws-xray-sdk-go/xray"

	"./common"
)

var (
	BalanceServiceURL = os.Getenv("BALANCE_SERVICE_URL")
	BalanceResource   = os.Getenv("BALANCE_RESOURCE")
)

// ContractTrade table
type ContractTrade struct {
	TradeNo          int     `gorm:"column:trade_no"  json:"tradeNo"`
	ContractDate     int     `gorm:"column:contract_date"  json:"contractDate"`
	ContractTime     string  `gorm:"column:contract_time"  json:"contractTime"`
	StockCode        string  `gorm:"column:stock_code"  json:"stockCode"`
	TradeType        string  `gorm:"column:trade_type"  json:"tradeType"`
	OrderType        string  `gorm:"column:order_type"  json:"orderType"`
	ContractPrice    float64 `gorm:"column:contract_price"  json:"contractPrice"`
	ContractQuantity int     `gorm:"column:contract_quantity"  json:"contractQuantity"`
	Fee              float64 `gorm:"column:fee"  json:"fee"`
	Tax              float64 `gorm:"column:tax"  json:"tax"`
	ContractValue    int     `gorm:"column:contract_value"  json:"contractValue"`
	SettlementAmount int     `gorm:"column:settlement_amount"  json:"settlementAmount"`
	UpdateTime       string  `gorm:"column:update_time"  json:"updateTime"`
}

type SQSContractTrade struct {
	TradeNo          int     `json:"tradeNo"`
	OrderStatus      bool    `json:"orderStatus"`
	ContractDate     int     `json:"contractDate"`
	ContractTime     string  `json:"contractTime"`
	StockCode        string  `json:"stockCode"`
	TradeType        string  `json:"tradeType"`
	OrderType        string  `json:"orderType"`
	ContractPrice    float64 `json:"contractPrice"`
	ContractQuantity int     `json:"contractQuantity"`
}

func MakeInsertStruct(sqsct *SQSContractTrade, rfcTime string) *ContractTrade {

	// get fee and tax from MySQL
	fee, tax := GetFeeTaxFromDatabase(sqsct.TradeNo)

	// calculate ContractValue
	contractValue := CalculateContractValue(sqsct, fee)

	// calculate settlementAmount
	settlementAmount := CalculateSettlementAmount(contractValue, tax)

	ct := ContractTrade{
		TradeNo:          sqsct.TradeNo,
		ContractDate:     sqsct.ContractDate,
		ContractTime:     sqsct.ContractTime,
		StockCode:        sqsct.StockCode,
		TradeType:        sqsct.TradeType,
		OrderType:        sqsct.OrderType,
		ContractPrice:    sqsct.ContractPrice,
		ContractQuantity: sqsct.ContractQuantity,
		Fee:              fee,
		Tax:              tax,
		ContractValue:    contractValue,
		SettlementAmount: settlementAmount,
		UpdateTime:       rfcTime,
	}

	return &ct
}

func GetFeeTaxFromDatabase(tradeNo int) (fee, tax float64) {

	var trade Trade
	// MySQLからselect
	db := common.ConnectDB("read")
	db.Debug().Table("trades").Where("trade_no = ?", tradeNo).Find(&trade)

	fee = trade.Fee
	tax = trade.Tax

	return fee, tax

}

func CalculateContractValue(sqsct *SQSContractTrade, fee float64) int {
	ContractValue := int(sqsct.ContractPrice)*sqsct.ContractQuantity + int(fee)
	return ContractValue
}

func CalculateSettlementAmount(contractValue int, tax float64) int {
	SettlementAmount := contractValue + int(tax)
	return SettlementAmount
}

func CheckSameTradeNoAlreadyExist(tradeNo int) bool {

	isResisted := false

	var ContractTrade ContractTrade
	db := common.ConnectDB("read")
	db.Debug().Table("contract_trades").Where("trade_no = ?", tradeNo).Find(&ContractTrade)

	if ContractTrade.TradeNo == 0 {
		isResisted = true
		fmt.Println("This trade is failed!!")
	}

	return isResisted
}

// 約定済みならテーブルinsert
func InsertContractTrade(ContractTrade ContractTrade) {

	isResisted := CheckSameTradeNoAlreadyExist(ContractTrade.TradeNo)

	if isResisted {
		db := common.ConnectDB("write")
		defer db.Close()

		db.NewRecord(ContractTrade)
		db.Create(&ContractTrade)
	} else {
		fmt.Println("This message has already registed!!")
	}
}

// 取引のステータスを更新
func UpdateTradeStatus(tradeNo int, status, rfcTime string) {

	var trade Trade

	// MySQLでupdate
	db := common.ConnectDB("write")
	db.Debug().Model(&trade).Where("trade_no = ?", tradeNo).Update(map[string]interface{}{"trade_status": status, "update_time": rfcTime})
	defer db.Close()
}

// 成立済み取引をStock-BalanceにPOST
func PostContractTradeToBalanceService(ctx context.Context, ContractTrade ContractTrade) {
	RequestUrl := BalanceServiceURL + BalanceResource

	postBody, _ := json.Marshal(ContractTrade)
	fmt.Printf("[+] %s\n", string(postBody))

	// request
	res, err := ctxhttp.Post(ctx, xray.Client(nil), RequestUrl, "application/json", bytes.NewBuffer(postBody))
	defer res.Body.Close()

	if err != nil {
		fmt.Println("[!] " + err.Error())
	} else {
		fmt.Println("[*] " + res.Status)
	}

}

func UpdateContractTrades(ctx context.Context, r string) {

	var sqsct SQSContractTrade

	if err := json.Unmarshal([]byte(r), &sqsct); err != nil {
		fmt.Println(err)
	}

	// get today date and time in JST
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Now().UTC().In(jst)
	rfcTime := t.Format(time.RFC3339)

	var status string

	// 約定なら成立済みテーブルにON DUPLICATE KEY UPDATEでinsert
	if sqsct.OrderStatus {
		status = "1" // trueなら1約定
		ContractTradeAddress := MakeInsertStruct(&sqsct, rfcTime)
		InsertContractTrade(*ContractTradeAddress)
		PostContractTradeToBalanceService(ctx, *ContractTradeAddress)
	} else {
		status = "2" // falseなら2fail
	}
	// 取引のステータスを更新
	UpdateTradeStatus(sqsct.TradeNo, status, rfcTime)

}
