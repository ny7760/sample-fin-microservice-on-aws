package common

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func SendMessageToSQS(QueURL, QueName, SegmentName, message string) {
	mySession := session.Must(session.NewSession())
	svc := sqs.New(mySession, aws.NewConfig().WithRegion(AwsRegion))
	xray.AWS(svc.Client)

	// セグメント定義
	ctx, seg := xray.BeginSegment(context.Background(), SegmentName)
	subctx, subseg := xray.BeginSubsegment(ctx, QueName)

	inputParams := &sqs.SendMessageInput{
		QueueUrl:    aws.String(QueURL),
		MessageBody: aws.String(message),
	}

	resp, err := svc.SendMessageWithContext(subctx, inputParams)
	if err != nil {
		fmt.Println("Send Message Error: ", err)
	}
	fmt.Println(resp)

	subseg.Close(nil)
	seg.Close(nil)

}
