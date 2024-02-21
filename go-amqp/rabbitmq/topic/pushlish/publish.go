package main

import (
	"fmt"
	"github.piwriw.go-amqp/rabbitmq/topic"
	"strconv"
	"time"
)

func main() {
	kutengOne := topic.NewRabbitMQTopic("JoohwanTopic", "joohwan.topic.one")
	kutengTwo := topic.NewRabbitMQTopic("JoohwanTopic", "joohwan.topic.two")
	for i := 0; i <= 100; i++ {
		kutengOne.PublishTopic("Hello joohwan topic one!" + strconv.Itoa(i))
		kutengTwo.PublishTopic("Hello joohwan topic Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
