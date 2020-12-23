package main

import (
	pb "app/api/blog/v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	fmt.Println("client start")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("faild to grpc dial:%v", err)
	}
	defer conn.Close()
	client := pb.NewPostClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := &pb.GetPostRequest{
		Id: 2,
	}
	res, err := client.GetPost(ctx, r)
	if err != nil {
		log.Fatalf("faild to Get Post:%v", err)
	}
	log.Println(res)

}
