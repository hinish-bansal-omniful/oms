package init

import (
	"context"
	"encoding/json"
	"fmt"

	// "oms-service/intersvc"
	"oms-service/domain"

	"oms-service/parse_csv"
	"oms-service/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/omniful/go_commons/log"
)

func getConfig() *aws.Config {
	if awsConfig == nil {
		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-north-1"))
		if err != nil {
			panic("Unable to connect to AWS")
		}
		awsConfig = &cfg
		return awsConfig
	}
	return awsConfig
}

func initialiseSQSConsumer(ctx context.Context) {

	sqsClient := sqs.NewFromConfig(*getConfig())

	sqURL := getSQSUrl(ctx)
	fmt.Println("Queue URL: ", *sqURL)

	// This will constantly listen to the SQS queue and print the messages
	// for {
	messagesResult, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl: sqURL,
	})
	if err != nil {
		fmt.Println("Unable to receive mesaages from SQS: ", err)
	}

	var sqsMessage domain.SQSMessage
	for _, message := range messagesResult.Messages {
		fmt.Println("Message: ", *message.Body)
		if err := json.Unmarshal([]byte(*message.Body), &sqsMessage); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}
	}

	// Parse CSV
	orders, err := parse_csv.ParseCSV(sqsMessage.FilePath)
	if err != nil {
		log.Printf("Error parsing CSV: %v", err)
	}

	// Save to MongoDB
	if err := repository.InsertOrders(ctx, orders, DB); err != nil {
		log.Printf("Error saving orders to database: %v", err)
	}

	// Interservice Call to WMS to check

}
