package providers

// QueueInterface define the methods the implementations will need to implement
type QueueInterface interface {
	Consume()
	Create()
	Delete()
	Process()
	Receive()
	Send()
}
