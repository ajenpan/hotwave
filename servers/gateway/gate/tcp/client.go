package tcp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type ClientOption func(*ClientOptions)

type ClientOptions struct {
	RemoteAddress string
	OnMessage     func(s *Client, p *Packet)
	OnConnStat    func(s *Client, state SocketStat)
}

// Address to bind to - host:port
// func Address(a string) ClientOption {
// 	return func(o *ClinetOptions) {
// 		o.Address = a
// 	}
// }

// func MessageCall(f tcp.OnMessageFunc) ClientOption {
// 	return func(o *ClinetOptions) {
// 		o.OnMessage = f
// 	}
// }

// func ConnStatCall(f tcp.OnConnStatFunc) ClientOption {
// 	return func(o *ClinetOptions) {
// 		o.OnConnStat = f
// 	}
// }

var DefaultClientOptions = ClientOptions{
	RemoteAddress: "",
	OnMessage:     func(s *Client, p *Packet) {},
	OnConnStat: func(s *Client, state SocketStat) {
	},
}

func NewClient(opts *ClientOptions) *Client {
	ret := &Client{opt: opts}
	return ret
}

type Client struct {
	*Socket
	opt   *ClientOptions
	mutex sync.Mutex
	// chDie chan bool // wait for close
}

//retrun socket work status
func (a *Client) Status() int32 {
	return atomic.LoadInt32(&a.state)
}

func (c *Client) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Socket != nil {
		c.Socket.Close()
	}

	conn, err := net.DialTimeout("tcp", c.opt.RemoteAddress, 10*time.Second)
	if err != nil {
		return err
	}
	socket := NewSocket(conn)
	c.Socket = socket

	// go func() {
	//send ack
	err = socket.writePacket(NewAckPacket())
	if err != nil {
		return err
	}
	p := &Packet{}

	//read ack
	err = socket.readPacket(p)
	if err != nil {
		return err
	}

	go func() {
		defer socket.Close()
		go socket.writeWork()

		if c.opt.OnConnStat != nil {
			c.opt.OnConnStat(c, SocketStatConnected)
			defer c.opt.OnConnStat(c, SocketStatDisconnected)
		}

		var socketErr error = nil
		for {
			p.Reset()
			if socketErr = socket.readPacket(p); socketErr != nil {
				break
			}

			if c.opt.OnMessage != nil {
				c.opt.OnMessage(c, p)
			}
		}
	}()
	return nil
}
