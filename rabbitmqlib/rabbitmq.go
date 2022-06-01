package rabbitmqlib

import (
	"fmt"

	"github.com/streadway/amqp"
)

//RabbitMqProxy RabbitMqProxy
type RabbitMqProxy struct {
	con *Connection
}

//address amqp://user:pass@hostName:port/vhost
//Dial connect rabbitmq
func (r *RabbitMqProxy) Dial(address string) {
	r.con = new(Connection)
	r.con.Connect(address)
}

//ReceiveQueryMessage 获取Query
func (r *RabbitMqProxy) ReceiveQueryMessage(queueName string, messageReceived func(<-chan amqp.Delivery)) error {
	channelItem, err := r.con.GetChannel()
	if err == nil {
		defer r.con.CloseChannel(channelItem)
		var queue amqp.Queue
		args := amqp.Table{}
		args["x-message-ttl"] = 5000
		queue, err = channelItem.CH.QueueDeclare(queueName, false, true, false, false, args)
		if err == nil {
			err = channelItem.CH.Qos(
				1,     // prefetch count
				0,     // prefetch size
				false, // global
			)
			if err == nil {
				var channel <-chan amqp.Delivery
				channel, err = channelItem.CH.Consume(queue.Name, "", false, false, false, false, nil)
				if err == nil {
					messageReceived(channel)
				}
			}
		}
	}
	return err
}

//SendQueryMessage 发送Query
func (r *RabbitMqProxy) SendQueryMessage(queueName string, datas []byte) error {
	channelItem, err := r.con.GetChannel()
	if err != nil {
		return err
	}
	defer r.con.CloseChannel(channelItem)

	args := amqp.Table{}
	args["x-message-ttl"] = 5000
	_, err = channelItem.CH.QueueDeclare(queueName, false, true, false, false, args)
	if err != nil {
		return err
	}
	err = channelItem.CH.Publish("", queueName, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        datas,
		})

	return err
}

//ReceiveRouterMessage receive p/s message router
func (r *RabbitMqProxy) ReceiveRouterMessage(exchangeName string, routingKeys ...string) (<-chan amqp.Delivery, error) {
	channelItem, err := r.con.GetChannel()
	if err != nil {
		return nil, err
	}
	defer r.con.CloseChannel(channelItem)
	err = channelItem.CH.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	q, err := channelItem.CH.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, err
	}
	err = channelItem.CH.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}
	for _, routingKey := range routingKeys {
		err = channelItem.CH.QueueBind(q.Name, routingKey, exchangeName, false, nil)
		if err != nil {
			return nil, err
		}
	}
	return channelItem.CH.Consume(q.Name, "", false, false, false, false, nil)
}

//SendRouterMessage send p/s message router
func (r *RabbitMqProxy) SendRouterMessage(exchangeName string, routingKey string, datas []byte) error {
	channelItem, err := r.con.GetChannel()
	if err != nil {
		return err
	}
	defer r.con.CloseChannel(channelItem)
	err = channelItem.CH.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = channelItem.CH.Publish(exchangeName, routingKey, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        datas,
		})
	if err != nil {
		return err
	}
	return nil
}

//ReceiveFanoutMessage receive p/s message router
func (r *RabbitMqProxy) ReceiveFanoutMessage(subscribe *Subscribe) error {
	channelItem, err := r.con.GetChannel()
	if err != nil {
		return err
	}
	defer r.con.CloseChannel(channelItem)
	err = channelItem.CH.ExchangeDeclare(subscribe.ExchangeName, "fanout", false, true, false, false, nil)
	if err != nil {
		return err
	}
	isDurable := subscribe.QueueName != ""
	q, err := channelItem.CH.QueueDeclare(subscribe.QueueName, isDurable, false, true, false, nil)
	if err != nil {
		return err
	}
	err = channelItem.CH.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}
	err = channelItem.CH.QueueBind(q.Name, "", subscribe.ExchangeName, false, nil)
	if err != nil {
		return err
	}
	delivery, err := channelItem.CH.Consume(q.Name, "", true, false, false, false, nil)
	if err == nil {
		for d := range delivery {
			subscribe.MessageReceived(d.Body)
		}
	} else {
		fmt.Println(err, "Consume Failed")
	}
	return err
}

//SendFanoutMessage send p/s message router
func (r *RabbitMqProxy) SendFanoutMessage(exchangeName string, datas []byte) error {
	channelItem, err := r.con.GetChannel()
	if err != nil {
		return err
	}
	defer r.con.CloseChannel(channelItem)
	err = channelItem.CH.ExchangeDeclare(exchangeName, "fanout", false, true, false, false, nil)
	if err != nil {
		return err
	}
	err = channelItem.CH.Publish(exchangeName, "", false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        datas,
		})
	if err != nil {
		return err
	}
	return nil
}

//ReceiveTopicMessage ReceiveTopicMessage
func (r *RabbitMqProxy) ReceiveTopicMessage(subscribe *Subscribe) error {
	channelItem, err := r.con.GetChannel()
	if err == nil {
		defer r.con.CloseChannel(channelItem)
		err = channelItem.CH.ExchangeDeclare(subscribe.ExchangeName, "topic", true, false, false, false, nil)
		if err == nil {
			isDurable := subscribe.QueueName != ""
			var queue amqp.Queue
			queue, err = channelItem.CH.QueueDeclare(subscribe.QueueName, isDurable, false, false, false, nil)
			if err == nil {
				err = channelItem.CH.Qos(
					1,     // prefetch count
					0,     // prefetch size
					false, // global
				)
				if err == nil {
					err = channelItem.CH.QueueBind(queue.Name, subscribe.RoutingKey, subscribe.ExchangeName, false, nil)
					if err == nil {
						delivery, err := channelItem.CH.Consume(queue.Name, "", true, false, false, false, nil)
						if err == nil {
							for d := range delivery {
								subscribe.MessageReceived(d.Body)
							}
						} else {
							fmt.Println(err, "Consume Failed")
						}
					} else {
						fmt.Println(err, "QueueBind Failed")
					}
				}
			} else {
				fmt.Println(err, "Queue declare Failed")
			}
		} else {
			fmt.Println(err, "Exhange declare Failed")
		}
	}
	return err
}

//SendTopicMessage send p/s message router
func (r *RabbitMqProxy) SendTopicMessage(exchangeName string, routingKey string, datas []byte) error {
	channelItem, err := r.con.GetChannel()
	if err != nil {
		return err
	}
	defer r.con.CloseChannel(channelItem)
	err = channelItem.CH.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = channelItem.CH.Publish(exchangeName, routingKey, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        datas,
		})
	return err
}
