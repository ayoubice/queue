package sns

import (
	"errors"
	"fmt"
	"strings"

	awsSns "github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

const (
	protocolSNS = "sns"
	protocolSQS = "sqs"
)

// QueueConfig is used to create a SNS service
type QueueConfig struct {
	QueueName          string
	QueueProvider      string
	SubscriberEndpoint *string
	SubscriberProtocol *string
}

// Service is the implementation of DynamoService
type Service struct {
	snsAPI      snsiface.SNSAPI
	queueConfig QueueConfig
}

// NewService returns ready to use instance of *DB
func NewService(snsAPI snsiface.SNSAPI, queueConfig QueueConfig) *Service {
	return &Service{
		snsAPI:      snsAPI,
		queueConfig: queueConfig,
	}
}

// PublishInput represent the input for Publish action.
// One of PhoneNumber, TopicArn or TargetArn must be specified.
type PublishInput struct {
	Message           *string                           `type:"string" required:"true"`
	MessageAttributes map[string]*MessageAttributeValue `type:"map"`
	MessageStructure  *string                           `type:"string"`
	PhoneNumber       *string                           `type:"string"`
	Subject           *string                           `type:"string"`
	TargetArn         *string                           `type:"string"`
	TopicArn          *string                           `type:"string"`
}

// MessageAttributeValue represent the attributes for the message
type MessageAttributeValue struct {
	BinaryValue []byte  `type:"blob"`
	DataType    *string `type:"string" required:"true"`
	StringValue *string `type:"string"`
}

// SubscribeInput for Subscribe action.
type SubscribeInput struct {
	Endpoint *string `type:"string"`
	Protocol *string `type:"string" required:"true"`
	TopicArn *string `type:"string" required:"true"`
}

// Create is used to create a Topic and subscribe endpoints. A provider and a consumer should be given.
func (service *Service) Create() (string, error) {
	input := &awsSns.CreateTopicInput{Name: &service.queueConfig.QueueName}

	output, err := service.snsAPI.CreateTopic(input)
	if err != nil {
		return "", err
	}
	if output.TopicArn == nil {
		return "", errors.New("Topic Arn is null")
	}

	queueEndpoint := service.snsToSqsEndPoint(*output.TopicArn)
	subscribeInput := SubscribeInput{
		Endpoint: &queueEndpoint,
		Protocol: &service.queueConfig.QueueProvider,
		TopicArn: output.TopicArn,
	}
	_, err = service.SubscribeTopic(subscribeInput)
	if err != nil {
		return "", fmt.Errorf("unable to subscribe queue for arn: %s. Error: %s", *output.TopicArn, err)
	}
	subscribeInput = SubscribeInput{
		Endpoint: service.queueConfig.SubscriberEndpoint,
		Protocol: service.queueConfig.SubscriberProtocol,
		TopicArn: output.TopicArn,
	}
	_, err = service.SubscribeTopic(subscribeInput)
	if err != nil {
		return "", err
	}

	return *output.TopicArn, nil
}

// SubscribeTopic subscribe an endpoint with a given protocol to the topic.
func (service *Service) SubscribeTopic(subscribeInput SubscribeInput) (string, error) {
	input := service.publishSubscribeInput(subscribeInput)
	output, err := service.snsAPI.Subscribe(input)
	if err != nil {
		return "", err
	}
	if output.SubscriptionArn == nil {
		return "", errors.New("subscription Arn is null")
	}

	return *output.SubscriptionArn, nil
}

// PublishToTopic prepares a message and publish it to a Topic.
func (service *Service) PublishToTopic(publishInput PublishInput) (string, error) {
	input := service.publishInputAssembler(publishInput)
	output, err := service.snsAPI.Publish(input)
	if err != nil {
		return "", err
	}
	if output.MessageId == nil {
		return "", errors.New("message id is null")
	}

	return *output.MessageId, nil
}

func (service *Service) publishInputAssembler(publishInput PublishInput) *awsSns.PublishInput {
	return &awsSns.PublishInput{
		Message:           publishInput.Message,
		MessageAttributes: service.publishInputMessageAttributesAssembler(publishInput.MessageAttributes),
		MessageStructure:  publishInput.MessageStructure,
		PhoneNumber:       publishInput.PhoneNumber,
		Subject:           publishInput.Subject,
		TargetArn:         publishInput.TargetArn,
		TopicArn:          publishInput.TopicArn,
	}
}

func (service *Service) publishInputMessageAttributesAssembler(
	publishInputAttributes map[string]*MessageAttributeValue,
) map[string]*awsSns.MessageAttributeValue {
	attributes := make(map[string]*awsSns.MessageAttributeValue, len(publishInputAttributes))

	for index, input := range publishInputAttributes {
		attributes[index] = &awsSns.MessageAttributeValue{
			BinaryValue: input.BinaryValue,
			DataType:    input.DataType,
			StringValue: input.StringValue,
		}
	}

	return attributes
}

func (service *Service) publishSubscribeInput(subscribeInput SubscribeInput) *awsSns.SubscribeInput {
	return &awsSns.SubscribeInput{
		Endpoint: subscribeInput.Endpoint,
		Protocol: subscribeInput.Protocol,
		TopicArn: subscribeInput.TopicArn,
	}
}

func (service *Service) snsToSqsEndPoint(snsEndpoint string) string {
	return strings.Replace(snsEndpoint, protocolSNS, protocolSQS, 1)
}
