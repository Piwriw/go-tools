package main

import (
	"github.piwriw.go-amqp/rabbitmq/publish"
)

func main() {
	go func() {
		rabbitmq := publish.NewRabbitMQPubSub("PublishMQ")
		rabbitmq.ReceiveSub()
	}()
	rabbitmq := publish.NewRabbitMQPubSub("PublishMQ")
	rabbitmq.ReceiveSub()

}
