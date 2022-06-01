package rabbitmqlib

//Subscribe Subscribe
type Subscribe struct {
	ExchangeName    string
	QueueName       string
	RoutingKey      string
	MessageReceived func([]byte)
}
