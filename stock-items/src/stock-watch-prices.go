package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"

	"./subfunc"
)

var (
	QueURL         = os.Getenv("STOCK_PRICE_QUE")
	QueName        = os.Getenv("STOCK_PRICE_QUE_NAME")
	AwsRegion      = os.Getenv("AWS_REGION")
	ServiceName    = os.Getenv("SERVICE_NAME")
	MaxMessages, _ = strconv.ParseInt(os.Getenv("STOCK_PRICE_MAX_MESSAGES"), 10, 64)
	PollingTime, _ = strconv.ParseInt(os.Getenv("STOCK_PRICE_POLLING_TIME"), 10, 64)
)

var svc *sqs.SQS

func GetMessage() error {
	params := &sqs.ReceiveMessageInput{
		AttributeNames:      aws.StringSlice([]string{"AWSTraceHeader"}),
		QueueUrl:            aws.String(QueURL),
		MaxNumberOfMessages: aws.Int64(MaxMessages), // 一度に取得するメッセージの最大数
		WaitTimeSeconds:     aws.Int64(PollingTime), // ロングポーリングの時間
	}
	ctx, seg := xray.BeginSegment(context.Background(), ServiceName)
	subctx, subseg := xray.BeginSubsegment(ctx, QueName)

	resp, err := svc.ReceiveMessageWithContext(subctx, params)
	if err != nil {
		return err
	}
	num := len(resp.Messages)
	fmt.Printf("Number of messages: %d\n", num)
	if num == 0 {
		fmt.Println("queus is enmpty")
		return nil
	}
	for _, msg := range resp.Messages {
		// メッセージが取得できたら更新
		fmt.Println(*msg.Body)
		subfunc.UpdatePrices(*msg.Body)
		if err := DeleteMessage(subctx, msg); err != nil {
			fmt.Println(err)
		}
	}
	subseg.Close(nil)
	seg.Close(nil)

	return nil
}

func DeleteMessage(ctx context.Context, msg *sqs.Message) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(QueURL),
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	_, err := svc.DeleteMessageWithContext(ctx, params)

	if err != nil {
		return err
	}
	return nil

}

func main() {
	priceSession := session.Must(session.NewSession())
	svc = sqs.New(priceSession, aws.NewConfig().WithRegion(AwsRegion))
	xray.AWS(svc.Client)

	// Polling
	for {
		if err := GetMessage(); err != nil {
			log.Fatal(err)
		}
	}
}
