package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	AwsRegion = os.Getenv("AWS_REGION")
)

func PublishMessageToSNS(TopicARN, TopicName, SegmentName, message string) {
	mySession := session.Must(session.NewSession())
	svc := sns.New(mySession, aws.NewConfig().WithRegion(AwsRegion))
	xray.AWS(svc.Client)

	// SNSフォーマットのJSON文字列を作成
	messageJson := map[string]string{
		"default": message,
		"sqs":     message,
	}
	bytes, err := json.Marshal(messageJson)
	if err != nil {
		fmt.Println("JSON marshal Error: ", err)
	}
	messageForSNS := string(bytes)

	inputPublish := &sns.PublishInput{
		Message:          aws.String(messageForSNS),
		MessageStructure: aws.String("json"),
		TopicArn:         aws.String(TopicARN),
	}

	ctx, seg := xray.BeginSegment(context.Background(), SegmentName)
	ctx, subseg := xray.BeginSubsegment(ctx, TopicName)

	// SNSにメッセージPublish
	MessageId, err := svc.PublishWithContext(ctx, inputPublish)
	if err != nil {
		fmt.Println("Publish Error: ", err)
	}
	subseg.Close(nil)
	seg.Close(nil)

	fmt.Println(MessageId)

}
