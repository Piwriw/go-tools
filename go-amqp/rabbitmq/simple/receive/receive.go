package main

import "github.piwriw.go-amqp/rabbitmq/simple"

func main() {
	mqSimple := simple.NewRabbitMQSimple("simpleRb")
	mqSimple.ConsumeSimple()
}
