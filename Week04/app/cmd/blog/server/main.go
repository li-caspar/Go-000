package main

import (
	pb "app/api/blog/v1"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("server start")
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return serverSignal(ctx)
	})
	g.Go(func() error {
		return serverGRPC(ctx, lis)
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("server error exit, error:%s\n", err)
	}
}

//监听信号量来决定是否退出
func serverSignal(ctx context.Context) error {
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-sc:
		return fmt.Errorf("signal exit, %s", sig.String())
	case <-ctx.Done():
		fmt.Println("context done, signal return")
		return nil
	}
}

func serverGRPC(ctx context.Context, lis net.Listener) error {
	server := grpc.NewServer()
	postService := InitializePostService()
	pb.RegisterPostServer(server, postService)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		select {
		case <-ctx.Done():
			fmt.Println("context done, grpc server shutdown")
			server.GracefulStop()
		}
	}()
	if err := server.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to grpc server")
	}
	return nil
}
