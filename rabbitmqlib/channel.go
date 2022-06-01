package rabbitmqlib

import "github.com/streadway/amqp"

//Channel Channel
type Channel struct {
	IsUsed   bool
	CH       *amqp.Channel
	Identify string
}
