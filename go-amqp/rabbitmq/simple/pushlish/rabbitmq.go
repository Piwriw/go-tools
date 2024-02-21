package main

import "github.piwriw.go-amqp/rabbitmq/simple"

func main() {
	mqSimple := simple.NewRabbitMQSimple("simpleRb")
	mqSimple.PublishSimple("send simple message to receive")
	mqSimple.PublishSimple("Second send simple message to receive")

}
