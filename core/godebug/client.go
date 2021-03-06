package godebug

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/jmigpin/editor/core/godebug/debug"
)

//var logger = log.New(os.Stdout, "godebug: ", log.Lshortfile)
var logger = log.New(ioutil.Discard, "godebug: ", 0)

type Client struct {
	Conn     net.Conn
	Messages chan interface{}
	done     chan struct{}
}

func NewClient(ctx context.Context) (*Client, error) {
	client := &Client{
		Messages: make(chan interface{}, 512),
		done:     make(chan struct{}),
	}
	if err := client.connect(ctx); err != nil {
		return nil, err
	}

	// receive msgs from server and send to channel
	go func() {
		client.receiveLoop()
	}()

	return client, nil
}

func (client *Client) Wait() {
	<-client.done
}

func (client *Client) Close() error {
	if client.Conn != nil {
		return client.Conn.Close()
	}
	return nil
}

func (client *Client) connect(ctx context.Context) error {
	// connect to server with retries during a period
	end := time.Now().Add(5 * time.Second)
	for {
		// connect
		var dialer net.Dialer
		conn0, err := dialer.DialContext(ctx, "tcp", ":8070")
		if err != nil {
			// retry while the end time is not reached
			if time.Now().Before(end) {
				timer := time.NewTimer(250 * time.Millisecond)
				select {
				case <-timer.C:
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			return err
		}

		// connected
		client.Conn = conn0
		return nil
	}
}

func (client *Client) receiveLoop() {
	for {
		msg, err := debug.DecodeMessage(client.Conn)
		if err != nil {
			logger.Print(err)

			// unable to read (server was probably closed)
			if operr, ok := err.(*net.OpError); ok {
				if operr.Op == "read" {
					break
				}
			}
			// connection ended gracefully by the client
			if err == io.EOF {
				break
			}

			continue
		}

		//logger.Printf("recv msg")
		client.Messages <- msg
	}

	// no more msgs
	close(client.Messages)

	close(client.done)
}
