package common

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var (
	Endpoint  = os.Getenv("DYNAMODB_ENDPOINT")
	AwsRegion = os.Getenv("AWS_REGION")
)

// Dynamoテーブル定義
type MaxNumber struct {
	ColumnName string `dynamo:"column_name"`
	TradeNo    int    `dynamo:"max_value"`
}

// GetNewTradeNo .. Dynamodbからアトミックカウンタを利用して取引番号最大の次の値を取得
func GetNewTradeNo() int {

	disableSSL := false
	if len(Endpoint) != 0 {
		disableSSL = true
	}

	db := dynamo.New(session.New(), &aws.Config{
		Region:     aws.String(AwsRegion),
		Endpoint:   aws.String(Endpoint),
		DisableSSL: aws.Bool(disableSSL),
	})

	tableName := "StockTrade.max_number"
	table := db.Table(tableName)

	var result MaxNumber
	// Atomic Counter
	err := table.Update("column_name", "trade_no").Add("max_value", 1).Value(&result)
	if err != nil {
		fmt.Printf("Failed to get item[%v]\n", err)
	}
	return result.TradeNo
}
