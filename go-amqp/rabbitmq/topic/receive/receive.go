package main

import (
	"github.piwriw.go-amqp/rabbitmq/topic"
)

func main() {
	kutengOne := topic.NewRabbitMQTopic("JoohwanTopic", "#")
	kutengOne.RecieveTopic()
}
