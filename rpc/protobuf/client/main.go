package main

import (
	"context"
	"fmt"
	"github.com/piwriw/protobuf/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func main() {
	name := "World"
	addr := "127.0.0.1:50051"
	// 连接gRPC服务器
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("connect error to: ", addr)
	}
	defer conn.Close()

	// 实例化一个client对象，传入参数conn
	c := service.NewHelloClient(conn)

	// 初始化上下文，设置请求超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//延迟关闭请求会话
	defer cancel()

	// 调用SayHello方法，以请求服务，然后得到响应消息
	r, err := c.SayHello(ctx, &service.HelloRequest{Name: name})
	if err != nil {
		fmt.Println("can not greet to: ", addr)
	} else {
		fmt.Println("response from server: ", r.GetMessage())
	}

}
