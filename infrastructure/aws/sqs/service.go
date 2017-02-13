package sqs

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// Service is the implementation of DynamoService
type SQSService struct {
	sqsiface.SQSAPI
	QueueName string
}

// NewService returns ready to use instance of *DB
func NewService(sqsAPI sqsiface.SQSAPI, queueName string) *SQSService {
	return &SQSService{
		SQSAPI:    sqsAPI,
		QueueName: queueName,
	}
}

// Create creates a queue in the sqs.
func (service *SQSService) Create(attributes map[string]*string) (string, error) {
	input := &sqs.CreateQueueInput{
		Attributes: attributes,
		QueueName:  &service.QueueName,
	}

	output, err := service.CreateQueue(input)
	if err != nil {
		return "", err
	}
	if output.QueueUrl == nil {
		return "", errors.New("QueueUrl is null")
	}

	return *output.QueueUrl, nil
}

// CreateWithoutAttributes creates a queue in the sqs.
func (service *SQSService) CreateWithoutAttributes() (string, error) {
	input := &sqs.CreateQueueInput{
		QueueName: &service.QueueName,
	}
	output, err := service.CreateQueue(input)
	if err != nil {
		return "", fmt.Errorf("Unable to create queu: %s", err)
	}
	if output.QueueUrl == nil {
		return "", fmt.Errorf("QueueUrl is null")
	}

	return *output.QueueUrl, nil
}

// Send puts a message into the queue
func (service *SQSService) Send(message string) (string, error) {
	queueURL, err := service.queueURL()
	if err != nil {
		return "", err
	}

	SendMessageInput := &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &queueURL,
	}

	output, err := service.SendMessage(SendMessageInput)
	if err != nil {
		return "", err
	}

	return *output.MessageId, nil
}

// Consume will consume a queue by queueURL returns the message
func (service *SQSService) Consume() (string, error) {
	queueURL, err := service.queueURL()
	if err != nil {
		return "", err
	}

	inputReceive := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(1),
		QueueUrl:            aws.String(queueURL),
	}
	output, err := service.SQSAPI.ReceiveMessage(inputReceive)
	if err != nil {
		return "", fmt.Errorf("unable to receive message from queue: %s ", err)
	}

	if output == nil || len(output.Messages) == 0 {
		log.Println("no messages in queue")
		return "", nil
	}

	firstMessage := *output.Messages[0]

	inputDelete := &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: firstMessage.ReceiptHandle,
	}
	_, err = service.SQSAPI.DeleteMessage(inputDelete)
	if err != nil {
		return "", fmt.Errorf("unable to delete message from queue %s ", err)
	}

	return *firstMessage.Body, nil
}

func (service *SQSService) queueURL() (string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: &service.QueueName,
	}

	output, err := service.GetQueueUrl(input)
	if err != nil {
		return "", err
	}

	if output.QueueUrl == nil {
		return "", errors.New("QueueUrl is null")
	}

	return *output.QueueUrl, nil
}

