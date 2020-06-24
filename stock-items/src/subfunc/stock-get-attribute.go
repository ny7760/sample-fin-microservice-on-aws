package subfunc

import (
	"github.com/gin-gonic/gin"
	// for gorm
	_ "github.com/go-sql-driver/mysql"

	"./common"
)

// テーブル定義

// Attribute ...
type Attribute struct { 
	StockCode    string `gorm:"column:stock_code"` // カラム名=スネークケースで自動的に合う。なので、本来は書かなくても動く
	StockName    string `gorm:"column:stock_name"`
	IndustryCode string `gorm:"column:industry_code"`
	UpdateTime   string `gorm:"column:update_time"`
}

// テーブル定義とJSON定義が同じなので共通の定義
// Industry
type Industry struct {
	IndustryCode string `json:"industryCode"`
	IndustryName string `json:"industryName"`
}

// 構造体定義

// Stock
type Stock struct {
	StockCode  string   `json:"stockCode"`
	StockName  string   `json:"stockName"`
	Industry   Industry `json:"industry"`
	UpdateTime string   `json:"updateTime"`
}

// SelectAtrribute ... DB selectメソッド
func SelectAtrribute(key string) (attribute Attribute) {

	db := common.ConnectDB("read")
	tableName := "attributes"
	if key != "" {
		db.Debug().Table(tableName).Where("stock_code = ?", key).Find(&attribute)
	} else {
		db.Debug().Table(tableName).Find(&attribute)
	}
	defer db.Close()
	return attribute
}

// SelectIndustry ... DB selectメソッド
func SelectIndustry(key string) (industry Industry) {

	db := common.ConnectDB("read")
	tableName := "industries"
	if key != "" {
		db.Debug().Table(tableName).Where("industry_code = ?", key).Find(&industry)
	} else {
		db.Debug().Table(tableName).Find(&industry)
	}
	defer db.Close()
	return industry
}

// GetAttribute ...
func GetAttribute(c *gin.Context) {

	attribute := SelectAtrribute(c.Query("stockCode"))
	industry := SelectIndustry(attribute.IndustryCode)

	stock := Stock{
		StockCode: attribute.StockCode,
		StockName: attribute.StockName,
		Industry: Industry{
			IndustryCode: attribute.IndustryCode,
			IndustryName: industry.IndustryName,
		},
		UpdateTime: attribute.UpdateTime,
	}

	c.JSON(200, stock)
}
