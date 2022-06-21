package tcp

import (
	"net"
	"sync"
	"time"

	log "hotwave/logger"
	"hotwave/services/gateway/gate"
)

type ServerOption func(*ServerOptions)

type ServerOptions struct {
	Address string

	HeatbeatInterval time.Duration
	Adapter          gate.GateAdapter
}

func NewServer(opts *ServerOptions) *Server {
	ret := &Server{
		opts:    opts,
		sockets: make(map[string]*Socket),
		die:     make(chan bool),
	}
	return ret
}

type Server struct {
	opts     *ServerOptions
	mu       sync.RWMutex
	sockets  map[string]*Socket
	die      chan bool
	wgLn     sync.WaitGroup
	wgConns  sync.WaitGroup
	listener net.Listener
}

func (s *Server) Stop() error {
	s.wgLn.Wait()
	select {
	case <-s.die:
	default:
		close(s.die)
	}
	s.wgConns.Wait()
	return nil
}

func (s *Server) Start() error {
	s.wgLn.Add(1)
	defer s.wgLn.Done()

	s.die = make(chan bool)
	listener, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}
	s.listener = listener
	go func() {
		defer listener.Close()
		var tempDelay time.Duration = 0
		for {
			select {
			case <-s.die:
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					if ne, ok := err.(net.Error); ok && ne.Temporary() {
						if tempDelay == 0 {
							tempDelay = 5 * time.Millisecond
						} else {
							tempDelay *= 2
						}
						if max := 1 * time.Second; tempDelay > max {
							tempDelay = max
						}
						time.Sleep(tempDelay)
						continue
					}
					return
				}
				tempDelay = 0

				socket := NewSocket(conn)
				go s.accept(socket)
			}
		}
	}()
	return nil
}

func (n *Server) accept(socket *Socket) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	n.wgConns.Add(1)
	defer n.wgConns.Done()
	defer socket.Close()

	//read ack
	p := &Packet{}
	if err := socket.readPacket(p); err != nil {
		log.Debugf("read ack error: %s", err)
		return
	}

	if p.Typ != PacketTypeAck {
		log.Debugf("ack type not right: %d", p.Typ)
		return
	}

	if err := socket.writePacket(p); err != nil {
		log.Debugf("write ack error: %s", err)
		return
	}

	// after ack, the connection is established
	go socket.writeWork()
	n.storeSocket(socket)
	defer n.removeSocket(socket)

	if n.opts.Adapter != nil {
		n.opts.Adapter.OnGateConnStat(socket, SocketStatConnected)
		defer n.opts.Adapter.OnGateConnStat(socket, SocketStatDisconnected)
	}

	var socketErr error = nil

	for {
		p.Reset()
		socketErr = socket.readPacket(p)
		if socketErr != nil {
			if operror, ok := socketErr.(*net.OpError); ok {
				log.Debug(operror.Error())
			}
			break
		}
		brk := false
		switch p.Typ {
		case PacketTypePacket:
			if n.opts.Adapter != nil {
				msg, err := ConverMessage(p)
				if err != nil {
					log.Error(err)
					continue
				}
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.Error(err)
						}
					}()
					n.opts.Adapter.OnGateAsync(socket, msg)
				}()
			}
		case PacketTypeHeartbeat:
			fallthrough
		case PacketTypeEcho:
			if err := socket.sendPacket(p); err != nil {
				log.Error(err)
				brk = true
			}
		default:
			brk = true
		}

		if brk {
			break
		}
	}
}

func (s *Server) GetSocket(id string) *Socket {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret, ok := s.sockets[id]
	if ok {
		return ret
	}
	return nil
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}

func (s *Server) SocketCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sockets)
}

func (s *Server) storeSocket(conn *Socket) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sockets[conn.ID()] = conn
}

func (s *Server) removeSocket(conn *Socket) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sockets, conn.ID())
}
