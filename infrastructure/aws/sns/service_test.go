package sns_test

import (
	"testing"

	"github.com/raffo0707/queue/config/aws"
	"github.com/raffo0707/queue/infrastructure/aws/sns"
	"github.com/raffo0707/queue/infrastructure/aws/sns/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	config := initConfig()
	service := sns.New(config)
	assert.NotNil(t, service)
	assert.Equal(t, service.APIVersion, "2010-03-31")
	assert.Equal(t, *service.Config.Endpoint, "fake-endpoint")
	assert.Equal(t, *service.Config.Region, "fake-region")
}

func TestCreation(t *testing.T) {
	snsService := initService()
	assert.NotNil(t, snsService)
}

// func TestCreate(t *testing.T) {
// 	config := initServiceConfig()
// 	topicArn := "fake-topic-Arn"
// 	snsServiceMock := initService()
// 	createTopicInput := &awsSns.CreateTopicInput{Name: &config.QueueName}
// 	createTopicOutput := &awsSns.CreateTopicOutput{TopicArn: &topicArn}

// 	snsServiceMock..On("CreateTopic", createTopicInput).Return(createTopicOutput, nil).Once()
// 	snsServiceMock.Create()
// }

func initConfig() aws.SNS {
	config := aws.Common{
		AccessKeyID:     "fake-access-key-id",
		SecretAccessKey: "fake-secret-access-key",
		SessionToken:    "fake-session-token",
		Endpoint:        "fake-endpoint",
		Region:          "fake-region",
	}

	return aws.SNS{
		Common:             config,
		TopicName:          "fake-topic-name",
		Provider:           "fake-provider",
		SubscriberEndpoint: "fake-subscriber-endpoint",
		SubscriberProtocol: "fake-subscriber-protocol",
	}
}

func initServiceConfig() sns.QueueConfig {
	SubscriberEndpoint := "fake-subscriber-endpoint"
	SubscriberProtocol := "fake-subscriber-protocol"

	return sns.QueueConfig{
		QueueName:          "fake-queue",
		QueueProvider:      "fake-provider",
		SubscriberEndpoint: &SubscriberEndpoint,
		SubscriberProtocol: &SubscriberProtocol,
	}
}
func initService() *sns.Service {
	type mockSNSClient struct {
		mocks.SNSAPI
	}

	service := mocks.SNSAPI{}
	APIMock := &mockSNSClient{service}

	snsConfig := initServiceConfig()

	return sns.NewService(APIMock, snsConfig)
}
