package main

import (
	"github.piwriw.go-amqp/rabbitmq/routing"
)

func main() {
	kutengone := routing.NewRabbitMQRouting("kuteng", "kuteng_one")
	kutengone.ReceiveRouting()
	go func() {
		kutengtwo := routing.NewRabbitMQRouting("kuteng", "kuteng_two")
		kutengtwo.ReceiveRouting()
	}()
}
