package main

import (
	"fmt"
	"github.piwriw.go-amqp/rabbitmq/routing"
	"strconv"
	"time"
)

func main() {
	kutengone := routing.NewRabbitMQRouting("kuteng", "kuteng_one")
	kutengtwo := routing.NewRabbitMQRouting("kuteng", "kuteng_two")
	for i := 0; i <= 100; i++ {
		kutengone.PublishRouting("Hello kuteng one!" + strconv.Itoa(i))
		kutengtwo.PublishRouting("Hello kuteng Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
