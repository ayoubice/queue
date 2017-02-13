package providers

// func TestCreate(t *testing.T) {
// 	t.Parallel()

// 	queueName, queueURL := "queue-name", "http://queue.url"
// 	attributes := make(map[string]*string)

// 	input := &sqs.CreateQueueInput{
// 		Attributes: attributes,
// 		QueueName:  &queueName,
// 	}
// 	output := &sqs.CreateQueueOutput{
// 		QueueUrl: &queueURL,
// 	}

// 	sqsAPI := &queueMock.QueueInterface{}
// 	sqsAPI.On(
// 		"CreateQueue",
// 		input,
// 	).Return(output, nil).Once()

// 	sqsService := NewAWSQueue(sqsAPI, queueName)

// 	_, err := sqsService.Create(attributes)
// 	assert.Nil(t, err)
// }
