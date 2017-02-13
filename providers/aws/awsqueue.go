package providers

import (
	"fmt"

	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/raffo0707/queue/providers"
)

// AWSQueue represents an implementation of QueueInterface over sqs and sns services
type AWSQueue struct {
	Queue providers.QueueInterface
	sqsiface.SQSAPI
	QueueName string
}

// NewAWSQueue returns ready to use instance of *DB
func NewAWSQueue(sqsAPI sqsiface.SQSAPI, queueName string) *AWSQueue {
	return &AWSQueue{
		SQSAPI:    sqsAPI,
		QueueName: queueName,
	}
}

// Create creates a queue in the sqs.
func (service *AWSQueue) Create(attributes map[string]*string) (string, error) {
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
func (service *AWSQueue) CreateWithoutAttributes() (string, error) {
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
func (service *AWSQueue) Send(message string) (string, error) {
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

// Receive will consume a queue by queueURL returns the message
func (service *AWSQueue) Receive(MaxNumberOfMessages int) error {
	queueURL, err := service.queueURL()
	if err != nil {
		return err
	}

	inputReceive := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(MaxNumberOfMessages),
		QueueUrl:            aws.String(queueURL),
	}
	output, err := service.SQSAPI.ReceiveMessage(inputReceive)
	if err != nil {
		return fmt.Errorf("unable to receive message from queue: %s ", err)
	}

	log.Println("Received ", len(output.Messages))

	if output == nil || len(output.Messages) == 0 {
		log.Println("no messages in queue")
		return nil
	}

	err = service.parseMessages(*output.Messages)
	if err != nil {
		return err
	}

	return nil
}

func (service *AWSQueue) Delete(message *sqs.Message) error {
	queueURL, err := service.queueURL()
	if err != nil {
		return err
	}

	inputDelete := &sqs.DeleteMessageInput{QueueUrl: &queueURL, ReceiptHandle: message.ReceiptHandle}
	_, err = service.SQSAPI.DeleteMessage(inputDelete)
	if err != nil {
		return fmt.Errorf("unable to delete message from queue %s ", err)
	}
}

func (service *AWSQueue) queueURL() (string, error) {
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

func (service *AWSQueue) parseMessages(messages []*sqs.Message) {
	var err error
	for _, message := range messages{
		err = service.Queue.Consume(message.Body)
		if err != nil {
			return err
		}

		err = service.Delete(message)
		if err != nil {
			return err
		}
	}
}

