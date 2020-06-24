package subfunc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	ServiceName            = os.Getenv("SERVICE_NAME")
	StockOrderTradeQueURL  = os.Getenv("STOCK_ORDER_TRADE_QUE")
	StockOrderTradeQueName = os.Getenv("STOCK_ORDER_TRADE_QUE_NAME")
)

type SendOrderTrade struct {
	TradeNo       int     `json:"tradeNo"`
	OrderDate     int     `json:"orderDate"`
	OrderTime     string  `json:"orderTime"`
	StockCode     string  `json:"stockCode"`
	TradeType     string  `json:"tradeType"`
	OrderType     string  `json:"orderType"`
	OrderPrice    float64 `json:"orderPrice"`
	OrderQuantity int     `json:"orderQuantity"`
}

func SendTradeToMarket(tradeNo int, ot *OrderTrade) string {

	// get today date and time in JST
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Now().UTC().In(jst)
	orderTime := t.Format("03:04:05")

	// 注文構造体を作成
	orderJson := SendOrderTrade{
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

	// SQSのセッション作成
	mySession := session.Must(session.NewSession())
	svc := sqs.New(mySession, aws.NewConfig().WithRegion(AwsRegion))
	xray.AWS(svc.Client)

	// セグメント定義
	ctx, seg := xray.BeginSegment(context.Background(), ServiceName)
	subctx, subseg := xray.BeginSubsegment(ctx, StockOrderTradeQueName)

	inputParams := &sqs.SendMessageInput{
		QueueUrl:    aws.String(StockOrderTradeQueURL),
		MessageBody: aws.String(message),
	}

	resp, err := svc.SendMessageWithContext(subctx, inputParams)
	if err != nil {
		fmt.Println("Send Message Error: ", err)
	}
	fmt.Println(resp)

	subseg.Close(nil)
	seg.Close(nil)

	return *resp.MessageId

}
