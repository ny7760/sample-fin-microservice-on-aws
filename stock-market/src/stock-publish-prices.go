package main

import (
	"os"
	"strconv"
	"time"

	"fin-micro/stock-market/subfunc"
	"fin-micro/stock-market/subfunc/common"
)

var (
	AwsRegion                 = os.Getenv("AWS_REGION")
	ServiceName               = os.Getenv("SERVICE_NAME")
	SegmentName               = os.Getenv("STOCK_PUBLISH_PRICES")
	TopicARN                  = os.Getenv("PRICE_TOPIC_ARN")
	TopicName                 = os.Getenv("PRICE_TOPIC_NAME")
	PriceQueURL               = os.Getenv("PRICE_QUE_URL")
	PriceQueName              = os.Getenv("PRICE_QUE_NAME")
	PriceProtocol             = os.Getenv("STOCK_PRICE_MESSAGE_MODE")
	GetPriceMinute, _         = strconv.Atoi(os.Getenv("GET_PRICE_MINUTE"))
	WaitForNextLoopSeconds, _ = strconv.Atoi(os.Getenv("GET_NEXT_PRICE_WAIT_SECOND"))
)

func main() {

	for {
		// get today date and time in JST
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		t := time.Now().UTC().In(jst)

		// get minute for judge
		minute, _ := strconv.Atoi(t.Format("04"))

		// get and publish price every 5 minute (default)
		if minute%GetPriceMinute == 0 {
			message := subfunc.SupplyRandomPrices()
			if PriceProtocol == "sns" {
				common.PublishMessageToSNS(TopicARN, TopicName, SegmentName, message)
			} else {
				common.SendMessageToSQS(PriceQueURL, PriceQueName, ServiceName, message)
			}
		}
		// default 60 seconds wait
		time.Sleep(time.Second * time.Duration(WaitForNextLoopSeconds))
	}

}
