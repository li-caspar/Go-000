package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn    net.Conn
	Server  *Server
	ctx     context.Context
	cancel  context.CancelFunc
	wBuffer chan []byte
	rBuffer chan []byte
}

type Server struct {
	addr     string
	port     int
	protocol string
	timeout  time.Duration
	maxConns int
	tls      *tls.Config
	listen   net.Listener

	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string) //这里可以扩展根据协议命令做回调
	clients                  []*Client
	wg                       sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

type Option func(*Server)

func WithProtocol(p string) Option {
	return func(s *Server) {
		s.protocol = p
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithMaxConns(maxcons int) Option {
	return func(s *Server) {
		s.maxConns = maxcons
	}
}

func WithTls(tls *tls.Config) Option {
	return func(s *Server) {
		s.tls = tls
	}
}

func WithNewClientCallback(f func(c *Client)) Option {
	return func(s *Server) {
		s.onNewClientCallback = f
	}
}

func WithClientConnectionClosed(f func(c *Client, err error)) Option {
	return func(s *Server) {
		s.onClientConnectionClosed = f
	}
}

func WithNewMessage(f func(c *Client, message string)) Option {
	return func(s *Server) {
		s.onNewMessage = f
	}
}

func NewServer(addr string, port int, options ...Option) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	svr := &Server{
		addr:     addr,
		port:     port,
		protocol: "tcp",
		timeout:  30 * time.Second, //设置默认值
		maxConns: 1000,
		tls:      nil,
		ctx:      ctx,
		cancel:   cancel,
		clients:  make([]*Client, 1024),
	}
	for _, option := range options {
		option(svr)
	}
	return svr, nil
}

func (s *Server) Stop() {
	s.cancel()
	s.listen.Close()
}

func (s *Server) Run() (err error) {
	s.wg.Add(1)
	addr := fmt.Sprintf("%s:%d", s.addr, s.port)
	log.Print("Server Run Listen int:", addr)
	s.listen, err = net.Listen(s.protocol, addr)
	if err != nil {
		return err
	}
	log.Print("Server After Listen")
	run := true
	for run {
		select {
		case <-s.ctx.Done():
			run = false
			break
		default:
			conn, err := s.listen.Accept()
			if err != nil {
				log.Printf("Accept error:%v\n", err)
				continue
			}
			log.Print("Accept ", conn.RemoteAddr())
			ctx, cancel := context.WithCancel(context.Background())
			client := &Client{
				conn:    conn,
				Server:  s,
				ctx:     ctx,
				cancel:  cancel,
				wBuffer: make(chan []byte, 1024),
				rBuffer: make(chan []byte, 1024),
			}
			s.clients = append(s.clients, client)
			s.wg.Add(2)
			//开始goroutine监听连接
			go client.handleConnRead()
			go client.handleConnWrite()
		}
	}
	s.wg.Done()
	log.Println("Server run exit")
	return nil
}

func (c *Client) handleConnRead() {
	defer c.conn.Close()
	defer c.Server.wg.Done()
	//读写缓冲区
	rd := bufio.NewReader(c.conn)
	run := true
	for run {
		select {
		case <-c.ctx.Done():
			run = false
			break
		default:
			line, _, err := rd.ReadLine()
			if err != nil {
				log.Printf("read error:%v\n", err)
				c.conn.Close()
				c.Server.onClientConnectionClosed(c, err)
				return
			}
			c.Server.onNewMessage(c, string(line))
		}

	}
}

func (c *Client) handleConnWrite() {
	defer c.conn.Close()
	defer c.Server.wg.Done()
	//读写缓冲区
	wr := bufio.NewWriter(c.conn)
	run := true
	for run {
		select {
		case <-c.ctx.Done():
			run = false
			break
		case msg := <-c.wBuffer:
			//回复消息
			n, err := wr.Write(msg)
			if n < len(msg) && err != nil {
				log.Printf("Write error:%v\n", err)
				continue
			}
			err = wr.Flush() //一次性syscall
			if err != nil {
				log.Printf("Flush error:%v\n", err)
				continue
			}
		}
	}
}

func (c *Client) Write(msg []byte) {
	c.wBuffer <- msg
}
