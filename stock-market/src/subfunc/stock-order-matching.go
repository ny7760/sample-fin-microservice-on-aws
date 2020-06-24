package subfunc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"fin-micro/stock-market/subfunc/common"
)

var (
	ServiceName          = os.Getenv("SERVICE_NAME")
	TopicARN             = os.Getenv("CONTRACT_TRADE_TOPIC_ARN")
	TopicName            = os.Getenv("CONTRACT_TRADE_TOPIC_NAME")
	ContractTradeQueURL  = os.Getenv("CONTRACT_TRADE_QUE_URL")
	ContractTradeQueName = os.Getenv("CONTRACT_TRADE_QUE_NAME")
	OrderProtocol        = os.Getenv("CONTRACT_TRADE_MESSAGE_MODE")
)

type RecieveOrderTrade struct {
	TradeNo       int     `json:"tradeNo"`
	OrderDate     int     `json:"orderDate"`
	OrderTime     string  `json:"orderTime"`
	StockCode     string  `json:"stockCode"`
	TradeType     string  `json:"tradeType"`
	OrderType     string  `json:"orderType"`
	OrderPrice    float64 `json:"orderPrice"`
	OrderQuantity int     `json:"orderQuantity"`
}

type ContractTrade struct {
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

// 成行注文
// 100%取引成立。-50〜+50円幅の金額で約定
func ValidateMarketOrder(rot *RecieveOrderTrade) *ContractTrade {

	var orderStatus bool
	var contractPrice float64
	var contractQuantity int

	orderStatus = true

	// set seed
	t := time.Now()
	rand.Seed(t.UnixNano())
	// random Price
	if rand.Intn(2) == 0 { // Intn(2)は50%の確率
		contractPrice = rot.OrderPrice + float64(rand.Intn(50))
	} else {
		contractPrice = rot.OrderPrice + float64((-1 * rand.Intn(50)))
	}

	contractQuantity = rot.OrderQuantity

	return EditContractTradeStruct(rot, orderStatus, contractPrice, contractQuantity)
}

// 指値注文
// 50%で取引成立
func ValidateLimitOrder(rot *RecieveOrderTrade) *ContractTrade {

	var orderStatus bool
	var contractPrice float64
	var contractQuantity int

	// set seed
	t := time.Now()
	rand.Seed(t.UnixNano())
	if rand.Intn(2) == 0 { // Intn(2)は50%の確率
		// 取引成立の場合
		orderStatus = true
		contractPrice = rot.OrderPrice
		contractQuantity = rot.OrderQuantity
	} else {
		orderStatus = false
		contractPrice = -99999999
		contractQuantity = 0
	}
	return EditContractTradeStruct(rot, orderStatus, contractPrice, contractQuantity)
}

func EditContractTradeStruct(rot *RecieveOrderTrade, status bool, price float64, quantity int) *ContractTrade {

	var contractDate int
	var contractTime string

	if status {
		// get today date and time in JST
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		t := time.Now().UTC().In(jst)
		contractDateString := t.Format("20060102")
		contractDate, _ = strconv.Atoi(contractDateString)
		contractTime = t.Format("03:04:05")
	} else {
		contractDate = 0
		contractTime = ""
	}

	ct := ContractTrade{
		TradeNo:          rot.TradeNo,
		OrderStatus:      status,
		ContractDate:     contractDate,
		ContractTime:     contractTime,
		StockCode:        rot.StockCode,
		TradeType:        rot.TradeType,
		OrderType:        rot.OrderType,
		ContractPrice:    price,
		ContractQuantity: quantity,
	}

	return &ct
}

func MatchingOrders(r, ServiceName, SegmentName string) {

	var rot RecieveOrderTrade
	var ct *ContractTrade

	// メッセージからJSONパース
	if err := json.Unmarshal([]byte(r), &rot); err != nil {
		fmt.Println(err)
	}

	// 指値か成行で処理分岐  0: 成行  1: 指値
	switch {
	case rot.OrderType == "0":
		ct = ValidateMarketOrder(&rot)
	case rot.OrderType == "1":
		ct = ValidateLimitOrder(&rot)
	default:
		fmt.Println("ERROR!")
	}

	bytes, err := json.Marshal(*ct)
	if err != nil {
		fmt.Println("JSON marshal Error: ", err)
	}

	if OrderProtocol == "sns" {
		common.PublishMessageToSNS(TopicARN, TopicName, SegmentName, string(bytes))
	} else {
		common.SendMessageToSQS(ContractTradeQueURL, ContractTradeQueName, ServiceName, string(bytes))
	}

}
