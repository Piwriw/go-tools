package main

import (
	"context"
	"github.com/piwriw/protobuf/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type UserServer struct {
}

func (u UserServer) mustEmbedUnimplementedHelloServer() {
	//TODO implement me
	panic("implement me")
}

func (u UserServer) SayHello(ctx context.Context, request *service.HelloRequest) (*service.HelloReply, error) {

	//TODO implement me
	panic("implement me")
}

func main() {
	grpcServer := grpc.NewServer()
	service.RegisterHelloServer(grpcServer, &UserServer{})
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(listen)

}
