package server

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestServer(t *testing.T) {
	s, err := NewServer("0.0.0.0", 13562,
		WithClientConnectionClosed(func(c *Client, err error) {
			t.Log("close client:", c.conn.RemoteAddr())
		}),
		WithNewClientCallback(func(c *Client) {
			t.Log("new client connect", c.conn.RemoteAddr())
		}),
		WithNewMessage(func(c *Client, message string) {
			t.Logf("recv msg from client:%s msg:%s", c.conn.RemoteAddr().String(), message)

			c.Write([]byte("hello from server"))
		}),
	)
	if err != nil {
		t.Fatal("new server err:" + err.Error())
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-signalChan:
			t.Log("recv exit signal")
			s.Stop()
			for _, client := range s.clients {
				if client == nil {
					continue
				}
				client.cancel()
			}
			t.Log("recv exit signal deal end")
			return
		}
	}()

	err = s.Run()
	if err != nil {
		t.Fatal("new server run err:" + err.Error())
	}
	t.Log("wait")
	s.wg.Wait()
}
