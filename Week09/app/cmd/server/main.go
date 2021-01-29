package main

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var shutdown = make(chan int8)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return serverSignal(ctx)
	})
	g.Go(func() error {
		return acceptTCP(ctx)
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("server error exit, error:%s\n", err)
	}
}

func serverSignal(ctx context.Context) error {
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-sc:
		shutdown <- 1
		return fmt.Errorf("signal exit, %s", sig.String())
	case <-ctx.Done():
		fmt.Println("content done, signal return")
		return nil
	}
}

func acceptTCP(ctx context.Context) error {
	var (
		addr *net.TCPAddr
		lis  *net.TCPListener
		conn *net.TCPConn
		err  error
	)
	//暂用固定的端口
	if addr, err = net.ResolveTCPAddr("tcp", ":8080"); err != nil {
		return fmt.Errorf("net.ResolveTCPAddr error:%s", err)
	}
	if lis, err = net.ListenTCP("tcp", addr); err != nil {
		return fmt.Errorf("net.ListenTCP error:%s", err)
	}
	//传递关闭信号,关闭监听连接
	go func(lis *net.TCPListener) {
		<-shutdown
		_ = lis.Close()
	}(lis)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done, TCP server shutdown")
			return nil
		default:
			if conn, err = lis.AcceptTCP(); err != nil {
				return err
			}
			go func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
				serverTCP(ctx, conn)
			}()
		}
	}
}

//处理连接
func serverTCP(ctx context.Context, conn *net.TCPConn) {
	ch := make(chan []byte, 10)
	ctxTCP, cancel := context.WithCancel(context.Background())
	defer func() {
		close(ch)
		cancel()
		_ = conn.Close()
	}()
	go dispatchTCP(conn, ctxTCP, ch) //负责写
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.ReadBytes('\n') //读取到换行
			if err != nil {
				fmt.Printf("reader error:%s\n", err)
				break
			}
			ch <- msg
		}
	}
}

//接收消息并写消息到连接
func dispatchTCP(conn *net.TCPConn, ctx context.Context, ch chan []byte) {
	wr := bufio.NewWriter(conn)
	for {
		select {
		case msg := <-ch:
			fmt.Printf("serverTCP msg:%s", msg)
			_, _ = wr.Write([]byte(fmt.Sprintf("rev:%s", msg))) //写数据到连接
			_ = wr.Flush()
		case <-ctx.Done():
			return
		}
	}
}
