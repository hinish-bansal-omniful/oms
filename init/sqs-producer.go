package init

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/omniful/go_commons/sqs"
)

var awsConfig *aws.Config
var exportedPublisher *sqs.Publisher

func GetNewSQSPublisher() *sqs.Publisher {
	return exportedPublisher
}

func setNewSQSPublisher(publisher *sqs.Publisher) {
	exportedPublisher = publisher
}

func getSQSUrl(ctx context.Context) *string {
	sqsURL, err := sqs.GetUrl(ctx, sqs.GetSQSConfig(ctx, false, "", "eu-north-1", "831926594619", "https://sqs.eu-north-1.amazonaws.com/831926594619/oms-wms"), "oms-wms")
	if err != nil {
		fmt.Println("error in connecting sqs")
	}
	return sqsURL
}

func initialiseSQSProducer(ctx context.Context) {

	sqsURL := getSQSUrl(ctx)
	fmt.Println("Queue URL: ", *sqsURL)

	newQueue, err := sqs.NewStandardQueue(ctx, "oms-wms", sqs.GetSQSConfig(ctx, false, "", "eu-north-1", "831926594619", "https://sqs.eu-north-1.amazonaws.com/831926594619/oms-wms"))
	if err != nil {
		fmt.Println("Error in creating queue")
	}

	NewSQSPublisher := sqs.NewPublisher(newQueue)
	setNewSQSPublisher(NewSQSPublisher)
}
