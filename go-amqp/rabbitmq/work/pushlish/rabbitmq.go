package main

import (
	"fmt"
	"github.piwriw.go-amqp/rabbitmq/simple"
	"strconv"
	"time"
)

func main() {
	mqSimple := simple.NewRabbitMQSimple("simpleRb")

	for i := 0; i <= 100; i++ {
		mqSimple.PublishSimple("Hello kuteng!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
