package subfunc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	DefaultPrice, _ = strconv.Atoi(os.Getenv("FIRST_PRICE"))
	SampleStockCode = os.Getenv("SAMPLE_STOCK_CODE")
	SampleStockName = os.Getenv("SAMPLE_STOCK_NAME")
)

// StockPrice
type StockPrice struct {
	StockCode string      `json:"stockCode"`
	StockName string      `json:"stockName"`
	DatePrice []DatePrice `json:"datePrice"`
}

// DatePrice
type DatePrice struct {
	Date      int         `json:"date"`
	TimePrice []TimePrice `json:"timePrice"`
}

// TimePrice ...
type TimePrice struct {
	Time  string  `json:"time"`
	Price float64 `json:"price"`
}

func SupplyRandomPrices() string {
	// get today date and time in JST
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Now().UTC().In(jst)

	t1Date, _ := strconv.Atoi(t.Format("20060102"))
	t1Time := t.Format("03:04")

	// set seed
	rand.Seed(t.UnixNano())
	// random Price
	var t1Price int
	if rand.Intn(2) == 0 { // Intn(2)は50%の確率
		t1Price = DefaultPrice + rand.Intn(50)
	} else {
		t1Price = DefaultPrice + (-1 * rand.Intn(50))
	}

	// make JSON string for publish
	priceJson := StockPrice{
		StockCode: SampleStockCode,
		StockName: SampleStockName,
		DatePrice: []DatePrice{
			{
				Date: t1Date,
				TimePrice: []TimePrice{
					{
						Time:  t1Time,
						Price: float64(t1Price),
					},
				},
			},
		},
	}
	bytes, err := json.Marshal(priceJson)
	if err != nil {
		fmt.Println("JSON marshal Error: ", err)
	}
	return string(bytes)

}
