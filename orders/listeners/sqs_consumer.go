package listeners

import (
	"awesomeProject5/OMS/orders/services"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"time"
)

func StartConsume(queueURL string, ctx context.Context) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Failed to load AWS config:", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	errorCount := 0 // Track consecutive errors to avoid infinite failures

	for {
		receiveMessages := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20, // Increase wait time to reduce empty polling
		}

		resp, err := sqsClient.ReceiveMessage(ctx, receiveMessages)
		if err != nil {
			log.Println("Error receiving message:", err)
			errorCount++
			if errorCount > 5 { // Stop after 5 consecutive failures
				log.Fatal("Too many errors, stopping consumer.")
			}
			time.Sleep(5 * time.Second) // Backoff before retrying
			continue
		}
		errorCount = 0 // Reset on success

		for _, message := range resp.Messages {
			if message.Body == nil {
				log.Println("Received an empty message, skipping...")
				continue
			}

			// Process the message
			services.ParseCSV(*message.Body, ctx) // FIX: Correct services import

			// Delete the message after processing
			_, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Println("Failed to delete message:", err)
			} else {
				fmt.Println("Message deleted successfully")
			}
		}
	}
}
