package subfunc

import (
	"encoding/json"
	"fmt"

	// for gorm
	"./common"
	_ "github.com/go-sql-driver/mysql"
)

// TimePrice ... テーブル構造体
type TimePrice struct {
	StockCode string
	Date      int
	Time      string
	Price     float64
}

// StockPrice .. 受け取ったJSON構造体
type StockPrice struct {
	StockCode string      `json:"stockCode"`
	StockName string      `json:"stockName"`
	DatePrice []DatePrice `json:"datePrice"`
}

// DatePrice
type DatePrice struct {
	Date      int           `json:"date"`
	TimePrice []PriceAtTime `json:"timePrice"`
}

// PriceAtTime ...
type PriceAtTime struct {
	Time  string  `json:"time"`
	Price float64 `json:"price"`
}

// InsertData ... DB insertメソッド
func InsertData(TimePrice TimePrice) {

	db := common.ConnectDB("write")
	// ON DUPLICATE KEY UPDATE
	db.NewRecord(TimePrice)
	db.Set(
		"gorm:insert_option",
		fmt.Sprintf("ON DUPLICATE KEY UPDATE `price` = %g", TimePrice.Price),
	).Create(&TimePrice)

}

// UpdatePrices ...
func UpdatePrices(r string) {

	var sp StockPrice

	// stringで受け取った値をUnmarshallでDICT化
	if err := json.Unmarshal([]byte(r), &sp); err != nil {
		fmt.Println(err)
	}

	code := sp.StockCode

	// JSONに含まれる日付ごとのループ処理
	for _, v1 := range sp.DatePrice {
		date := v1.Date

		for _, v2 := range v1.TimePrice {
			time := v2.Time
			price := v2.Price
			tp := TimePrice{
				StockCode: code,
				Date:      date,
				Time:      time,
				Price:     price,
			}
			InsertData(tp)
		}

	}

}
