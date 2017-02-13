package aws

// Common allows to configure base services configurations
type Common struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Endpoint        string
	Region          string
}

//SNS defines the properties for sns service
type SNS struct {
	Common
	TopicName          string
	Provider           string
	SubscriberEndpoint string
	SubscriberProtocol string
}

//SQS defines the properties for sqs service
type SQS struct {
	Common
	QueueName string
}
