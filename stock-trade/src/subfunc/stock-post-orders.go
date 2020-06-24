package subfunc

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"

	// for gorm
	_ "github.com/go-sql-driver/mysql"

	"./common"
)

var (
	OrderProtocol = os.Getenv("STOCK_ORDER_MESSAGE_MODE")
	AwsRegion     = os.Getenv("AWS_REGION")
)

// Trade ...
type Trade struct { // 構造体名の単数形 -> テーブル名(複数形)で自動的に合う
	TradeNo        int     `gorm:"column:trade_no"`
	OrderDate      int     `gorm:"column:order_date"`
	OrderTime      string  `gorm:"column:order_time"`
	StockCode      string  `gorm:"column:stock_code"`
	TradeType      string  `gorm:"column:trade_type"`
	OrderType      string  `gorm:"column:order_type"`
	OrderPrice     float64 `gorm:"column:order_price"`
	OrderQuantity  int     `gorm:"column:order_quantity"`
	Fee            float64 `gorm:"column:fee"`
	Tax            float64 `gorm:"column:tax"`
	EstimatedValue float64 `gorm:"column:estimated_value"`
	TradeStatus    string  `gorm:"column:trade_status"`
	UpdateTime     string  `gorm:"column:update_time"`
}

// 構造体定義

// OrderTrade
type OrderTrade struct {
	OrderDate     int     `json:"orderDate"`
	OrderTime     string  `json:"orderTime"`
	StockCode     string  `json:"stockCode"`
	TradeType     string  `json:"tradeType"`
	OrderType     string  `json:"orderType"`
	OrderPrice    float64 `json:"orderPrice"`
	OrderQuantity int     `json:"orderQuantity"`
}

// OrderAccept
type OrderAccept struct {
	MessageID        string `json:"messageId"`
	AcceptanceResult Trade  `json:"acceptanceResult"`
}

// InsertTrade ... DB insertメソッド
func InsertTrade(ot *OrderTrade) (trade Trade) {

	// dynamoから取引番号取得
	tradeNo := common.GetNewTradeNo()

	// 手数料、税金、見積もり金額を取得
	fee, tax, estimatedValue := common.CalculateFeeEst(ot.OrderPrice, ot.OrderQuantity)

	// get today date and time in JST
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Now().UTC().In(jst)
	rfcTime := t.Format(time.RFC3339)

	// メソッド名の行でTradeで型宣言したので、Tradeで初期化する必要あり
	trade = Trade{
		TradeNo:        tradeNo,
		OrderDate:      ot.OrderDate,
		OrderTime:      ot.OrderTime,
		StockCode:      ot.StockCode,
		TradeType:      ot.TradeType,
		OrderType:      ot.OrderType,
		OrderPrice:     ot.OrderPrice,
		OrderQuantity:  ot.OrderQuantity,
		Fee:            fee,
		Tax:            tax,
		EstimatedValue: estimatedValue,
		TradeStatus:    "0", // 未約定
		UpdateTime:     rfcTime,
	}

	// MySQLにinsert
	db := common.ConnectDB("write")
	db.Create(&trade)
	defer db.Close()

	return trade

}

// GetAttribute ...
func OrderTrades(c *gin.Context) {

	// POSTのBody取得
	var requestBody OrderTrade
	c.BindJSON(&requestBody)

	// POSTされた情報取得してinsert
	var r Trade
	r = InsertTrade(&requestBody)

	// insert OKならsns or sqsに送信
	var messageId string
	// 同じ階層の別ソース
	if OrderProtocol == "sns" {
		messageId = PublishTradeToMarket(r.TradeNo, &requestBody)
	} else {
		messageId = SendTradeToMarket(r.TradeNo, &requestBody)
	}

	result := OrderAccept{
		MessageID:        messageId,
		AcceptanceResult: r,
	}

	c.JSON(200, result)
}
