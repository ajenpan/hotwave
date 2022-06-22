package tcpsvr

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	log "hotwave/logger"
)

type ClientOption func(*ClientOptions)

type ClientOptions struct {
	RemoteAddress string
	OnMessage     func(s *Client, p *Packet)
	OnConnStat    func(s *Client, state SocketStat)
}

// Address to bind to - host:port
// func WithAddress(a string) ClientOption {
// 	return func(o *ClinetOptions) {
// 		o.Address = a
// 	}
// }

// func WithOnMessageFunc(f tcp.OnMessageFunc) ClientOption {
// 	return func(o *ClinetOptions) {
// 		o.OnMessage = f
// 	}
// }

// func WithOnConnStatFunc(f tcp.OnConnStatFunc) ClientOption {
// 	return func(o *ClinetOptions) {
// 		o.OnConnStat = f
// 	}
// }

var DefaultClientOptions = ClientOptions{
	RemoteAddress: "",
	OnMessage:     func(s *Client, p *Packet) {},
	OnConnStat:    func(s *Client, state SocketStat) {},
}

func NewClient(opts *ClientOptions) *Client {
	ret := &Client{
		Opt:    opts,
		Socket: NewSocket(nil),
	}
	return ret
}

type Client struct {
	*Socket
	Opt   *ClientOptions
	mutex sync.Mutex
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

	conn, err := net.DialTimeout("tcp", c.Opt.RemoteAddress, 10*time.Second)
	if err != nil {
		return err
	}
	socket := NewSocket(conn)

	//send ack
	err = socket.writePacket(NewAckPacket(nil))
	if err != nil {
		socket.Close()
		return err
	}

	p := &Packet{}
	//read ack
	err = socket.readPacket(p)
	if err != nil {
		socket.Close()
		return err
	}

	c.Socket = socket
	c.Socket.id = string(p.Raw)

	//here is connect finished
	go func() {
		defer socket.Close()
		go socket.writeWork()

		if c.Opt.OnConnStat != nil {
			c.Opt.OnConnStat(c, SocketStatConnected)
			defer c.Opt.OnConnStat(c, SocketStatDisconnected)
		}

		go func() {
			tk := time.NewTicker(5 * time.Second)
			defer tk.Stop()

			for {
				select {
				case <-tk.C:
					socket.Send(HeartbeatPakcet)
				case <-socket.chClosed:
					return
				}
			}
		}()

		var socketErr error = nil
		for {
			p.Reset()
			if socketErr = socket.readPacket(p); socketErr != nil {
				break
			}
			switch p.Typ {
			case PacketTypeAck:
			case PacketTypeHeartbeat:
				continue
			default:
				if c.Opt.OnMessage != nil {
					c.Opt.OnMessage(c, p)
				}
			}
		}
		log.Info(socketErr)
	}()
	return nil
}
