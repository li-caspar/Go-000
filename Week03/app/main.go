package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	//启用debug服务
	g.Go(func() error {
		return serverDebug(ctx)
	})
	//启动http服务
	g.Go(func() error {
		return server(ctx)
	})

	/*g.Go(func() error {
		return serverTest()
	})*/

	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//根据信号量来退出
	g.Go(func() error {
		return serverSignal(ctx, sc)
	})

	if err := g.Wait(); err != nil {
		fmt.Println("server error exit, error:" + err.Error())
	}
}

//信号量服务
func serverSignal(ctx context.Context, sc chan os.Signal) error {
	select {
	case sig := <-sc:
		return fmt.Errorf("signal exit, %s", sig.String())
	case <-ctx.Done():
		fmt.Println("context done, signa return")
		return nil
	}
}

//http服务
func server(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		select {
		case <-ctx.Done():
			fmt.Println("context done, server shutdown")
			_ = srv.Shutdown(ctx)
		}
	}()
	if err := srv.ListenAndServe(); err != nil {
		return errors.Wrap(err, "http server fail")
	}
	return nil
}

//debug服务
func serverDebug(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	srv := http.Server{
		Addr:    ":8090",
		Handler: mux,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		select {
		case <-ctx.Done():
			fmt.Println("context done, serverTest shutdown")
			_ = srv.Shutdown(ctx)
		}
	}()
	if err := srv.ListenAndServe(); err != nil {
		return errors.Wrap(err, "http serverTest fail")
	}
	return nil
}

//test服务 用于测试是否真的终止server
func serverTest() error {
	time.Sleep(10 * time.Second)
	fmt.Println("go test exit")
	return errors.New("test")
}
