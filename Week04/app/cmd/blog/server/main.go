package main

import (
	pb "app/api/blog/v1"
	"app/internal/config"
	"context"
	"flag"
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

var cfgFileName string

func init() {
	flag.StringVar(&cfgFileName, "c", "./configs/server.yaml", "config file path")
}

func main() {
	fmt.Println("server start")
	flag.Parse()
	cfg, err := config.NewConfig(cfgFileName)
	if err != nil {
		log.Fatalf("read config error:%v", err)
	}
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return serverSignal(ctx)
	})
	g.Go(func() error {
		return serverGRPC(ctx, cfg)
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

func serverGRPC(ctx context.Context, cfg *config.Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	server := grpc.NewServer()
	postService, err := InitializePostService(cfg)
	if err != nil {
		log.Fatalf("initialize post service error:%s", err)
	}
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
