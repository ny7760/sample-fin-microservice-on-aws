package subfunc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	// "github.com/gin-gonic/gin"

	// for gorm
	_ "github.com/go-sql-driver/mysql"

	"./common"
)

// Trade ...
type CashBalance struct { 
	CashValue  int    `gorm:"column:cash_value"`
	UpdateTime string `gorm:"column:update_time"`
}

type StockBalance struct { 
	StockCode     string `gorm:"column:stock_code"`
	StockValue    int    `gorm:"column:stock_value"`
	StockQuantity int    `gorm:"column:stock_quantity"`
	UpdateTime    string `gorm:"column:update_time"`
}

// 構造体定義

// ContractTrade
type ContractTrade struct {
	TradeNo          int     `json:"tradeNo"`
	ContractDate     int     `json:"contractDate"`
	ContractTime     string  `json:"contractTime"`
	StockCode        string  `json:"stockCode"`
	TradeType        string  `json:"tradeType"`
	OrderType        string  `json:"orderType"`
	ContractPrice    float64 `json:"contractPrice"`
	ContractQuantity int     `json:"contractQuantity"`
	Fee              float64 `json:"fee"`
	Tax              float64 `json:"tax"`
	ContractValue    int     `json:"contractValue"`
	SettlementAmount int     `json:"settlementAmount"`
	UpdateTime       string  `json:"updateTime"`
}

type UpdateStock struct {
	StockCode     string `json:"stockCode"`
	StockValue    int    `json:"stockValue"`
	StockQuantity int    `json:"stockQuantity"`
}

type ReturnValue struct {
	CashValue     int    `json:"cashValue"`
	StockCode     string `json:"stockCode"`
	StockValue    int    `json:"stockValue"`
	StockQuantity int    `json:"stockQuantity"`
}

func UpdateCash(ct *ContractTrade, rfcTime string) int {

	var cashBalance CashBalance

	db := common.ConnectDB("write")
	defer db.Close()

	tableName := "cash_balances"
	db.Debug().Table(tableName).Find(&cashBalance)
	fmt.Printf("Before Cash Balance: %v\n", cashBalance.CashValue)

	// if trade type is not Buy, it needs to be plused
	afterCashValue := cashBalance.CashValue - ct.SettlementAmount

	fmt.Printf("After Cash Balance: %v\n", afterCashValue)
	db.Debug().Table(tableName).Update(map[string]interface{}{"cash_value": afterCashValue, "update_time": rfcTime})

	return afterCashValue
}

func UpdateStockBalance(ct *ContractTrade, rfcTime string) *UpdateStock {

	var stockBalance StockBalance

	db := common.ConnectDB("write")
	defer db.Close()

	tableName := "stock_balances"
	db.Debug().Table(tableName).Find(&stockBalance)

	var afterStockValue, afterStockQuantity int

	// 初購入ならinsert
	if stockBalance.StockCode == "" {
		fmt.Println("Before Stock Value & Quantity: 0")
		afterStockValue = ct.ContractValue
		afterStockQuantity = ct.ContractQuantity

		stockBalance.StockCode = ct.StockCode
		stockBalance.StockValue = afterStockValue
		stockBalance.StockQuantity = afterStockQuantity
		stockBalance.UpdateTime = rfcTime

		db.Create(&stockBalance)

	} else {
		fmt.Printf("Before Stock Value: %v\n", stockBalance.StockValue)
		fmt.Printf("Before Stock Quantity: %v\n", stockBalance.StockQuantity)
		afterStockValue = stockBalance.StockValue + ct.ContractValue
		afterStockQuantity = stockBalance.StockQuantity + ct.ContractQuantity

		tableName := "stock_balances"
		db.Debug().Table(tableName).Where("stock_code = ?", ct.StockCode).Update(map[string]interface{}{"stock_value": afterStockValue, "stock_quantity": afterStockQuantity, "update_time": rfcTime})

	}
	fmt.Printf("Stock Code: %v\n", ct.StockCode)
	fmt.Printf("After Stock Value: %v\n", afterStockValue)
	fmt.Printf("After Stock Quantity: %v\n", afterStockQuantity)

	updateStock := UpdateStock{
		StockCode:     ct.StockCode,
		StockValue:    afterStockValue,
		StockQuantity: afterStockQuantity,
	}

	return &updateStock

}

func UpdateAssets(ct *ContractTrade, rfcTime string) ReturnValue {
	// Cash残高の更新
	afterCashValue := UpdateCash(ct, rfcTime)
	// 資産残高の更新
	var us *UpdateStock
	us = UpdateStockBalance(ct, rfcTime)

	rv := ReturnValue{
		CashValue:     afterCashValue,
		StockCode:     us.StockCode,
		StockValue:    us.StockValue,
		StockQuantity: us.StockQuantity,
	}

	return rv

}

// UpdateBalance ... if use gin
// func UpdateBalance(c *gin.Context) {
func UpdateBalance(w http.ResponseWriter, r *http.Request) {

	// get today date and time in JST
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Now().UTC().In(jst)
	rfcTime := t.Format(time.RFC3339)

	// POSTのBody取得 if use gin
	// c.BindJSON(&requestBody)
	var requestBody ContractTrade
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var returnValue ReturnValue
	returnValue = UpdateAssets(&requestBody, rfcTime)

	// return value if use gin
	// c.JSON(200, returnValue)
	returnJson, err := json.Marshal(returnValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write(returnJson)
}
