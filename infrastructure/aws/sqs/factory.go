package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	awsconfig "github.com/raffo0707/queue/config/aws"
)

// New returns configured and ready to use instance of SQS
func New(config awsconfig.SQS) *sqs.SQS {
	return sqs.New(
		session.New(
			&aws.Config{
				Endpoint: aws.String(config.Endpoint),
				Region:   aws.String(config.Region),
				Credentials: credentials.NewStaticCredentials(
					config.AccessKeyID,
					config.SecretAccessKey,
					config.SessionToken,
				),
				// force AWS to use http.DefaultClient
				EC2MetadataDisableTimeoutOverride: aws.Bool(true),
			},
		),
	)
}

