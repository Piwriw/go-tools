package main

import (
	"github.piwriw.go-amqp/rabbitmq/simple"
)

func main() {
	go receiveMQ("MQ1")
	receiveMQ("MQ2")
}
func receiveMQ(mqName string) {
	mqSimple := simple.NewRabbitMQSimple("simpleRb")
	mqSimple.ConsumeSimple()
}
