package subfunc

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	TopicARN = os.Getenv("TRADE_TOPIC_ARN")
)

type PublishOrderTrade struct {
	TradeNo       int     `json:"tradeNo"`
	OrderDate     int     `json:"orderDate"`
	OrderTime     string  `json:"orderTime"`
	StockCode     string  `json:"stockCode"`
	TradeType     string  `json:"tradeType"`
	OrderType     string  `json:"orderType"`
	OrderPrice    float64 `json:"orderPrice"`
	OrderQuantity int     `json:"orderQuantity"`
}

func PublishTradeToMarket(tradeNo int, ot *OrderTrade) string {

	// get today date and time in JST
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Now().UTC().In(jst)
	orderTime := t.Format("03:04:05")

	// 注文構造体を作成
	orderJson := PublishOrderTrade{
		TradeNo:       tradeNo,
		OrderDate:     ot.OrderDate,
		OrderTime:     orderTime,
		StockCode:     ot.StockCode,
		TradeType:     ot.TradeType,
		OrderType:     ot.OrderType,
		OrderPrice:    ot.OrderPrice,
		OrderQuantity: ot.OrderQuantity,
	}
	bytes, err := json.Marshal(orderJson)
	if err != nil {
		fmt.Println("JSON marshal Error: ", err)
	}
	message := string(bytes)

	// SNSのセッション作成
	mySession := session.Must(session.NewSession())
	svc := sns.New(mySession, aws.NewConfig().WithRegion(AwsRegion))

	// SNSフォーマットのJSON文字列を作成
	messageJson := map[string]string{
		"default": message,
		"sqs":     message,
	}
	bytes2, err := json.Marshal(messageJson)
	if err != nil {
		fmt.Println("JSON marshal Error: ", err)
	}
	messageForSNS := string(bytes2)

	inputPublish := &sns.PublishInput{
		Message:          aws.String(messageForSNS),
		MessageStructure: aws.String("json"),
		TopicArn:         aws.String(TopicARN),
	}

	// SNSにメッセージPublish
	MessageId, err := svc.Publish(inputPublish)
	if err != nil {
		fmt.Println("Publish Error: ", err)
	}

	return *MessageId.MessageId

}
