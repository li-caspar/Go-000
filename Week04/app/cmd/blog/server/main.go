package main

import (
	pb "app/api/blog/v1"
	"app/internal/service"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	fmt.Println("server start")
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	server := grpc.NewServer()
	s := &service.PostService{}
	pb.RegisterPostServer(server, s)
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("failed to grpc serve:%v", err)
	}
}
