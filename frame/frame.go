package frame

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	protobuf "google.golang.org/protobuf/proto"

	protocol "hotwave/frame/proto"
	"hotwave/logger"
	"hotwave/registry"
	"hotwave/util/addr"
	"hotwave/util/backoff"
)

func NewCore(opts ...Option) (*Core, error) {
	ret := &Core{
		exit: make(chan chan error),
		// wg:          wait(.Context),

	}

	options := newOptions(opts...)

	if options.NodeId == "" {
		options.NodeId = uuid.Must(uuid.NewUUID()).String()
	}
	ret.opts = options

	grpcServer := grpc.NewServer()
	protocol.RegisterNodeBaseServer(grpcServer, ret)
	ret.grpcServer = grpcServer

	return ret, nil
}

type Core struct {
	protocol.UnimplementedNodeBaseServer

	sync.RWMutex
	// marks the serve as started
	started bool
	// used for first registration
	registered bool
	wg         *sync.WaitGroup
	rsvc       *registry.Service
	opts       Options
	grpcServer *grpc.Server
	exit       chan chan error
}

func (c *Core) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	c.grpcServer.RegisterService(desc, impl)
}

func (c *Core) UserMessage(svr protocol.NodeBase_UserMessageServer) error {
	// um ,err:= svr.Recv()

	return nil
}

func (c *Core) EventMessage(svr protocol.NodeBase_EventMessageServer) error {
	return nil
}

func (c *Core) Stop() error {
	c.grpcServer.Stop()
	return nil
}

func (c *Core) Start() error {
	c.RLock()
	if c.started {
		c.RUnlock()
		return nil
	}
	c.RUnlock()

	config := c.Options()

	lis, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}
	go c.grpcServer.Serve(lis)

	//update for real address when the port is random, eg. localhost:0
	addr := c.opts.Address
	c.opts.Address = lis.Addr().String()
	logger.Info(addr)

	if err = c.Register(); err != nil {
		return err
	}

	// signals := make(chan os.Signal, 1)
	// signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// s := <-signals
	// logger.Infof("recv signal: %v", s.String())

	c.Lock()
	c.started = true
	c.Unlock()
	return nil
}

func (c *Core) SendMessageToUser(user User, msg protobuf.Message) error {

	return nil
}

func (s *Core) Options() Options {
	s.RLock()
	opts := s.opts
	s.RUnlock()
	return opts
}

func (s *Core) KeepRegistered() {
	go func() {
		config := s.Options()
		t := new(time.Ticker)
		// only process if it exists
		if s.opts.RegisterInterval > time.Duration(0) {
			// new ticker
			t = time.NewTicker(s.opts.RegisterInterval)
		}

		var err error
	Loop:
		for {
			select {
			// register self on interval
			case <-t.C:
				s.RLock()
				registered := s.registered
				s.RUnlock()
				rerr := s.opts.RegisterCheck(context.Background())
				if rerr != nil && registered {
					logger.Errorf("Server %s-%s register check error: %s, deregister it", config.Name, config.NodeId, err)
					// deregister self in case of error
					if err := s.Deregister(); err != nil {
						logger.Errorf("Server %s-%s deregister error: %s", config.Name, config.NodeId, err)
					}
				} else if rerr != nil && !registered {
					logger.Errorf("Server %s-%s register check error: %s", config.Name, config.NodeId, err)
					continue
				}
				if err := s.Register(); err != nil {
					logger.Errorf("Server %s-%s register error: %s", config.Name, config.NodeId, err)
				}
			// wait for exit
			case ch := <-s.exit:
				t.Stop()
				ch <- err
				break Loop
			}
		}

		s.RLock()
		registered := s.registered
		s.RUnlock()
		if registered {
			// deregister self
			if err := s.Deregister(); err != nil {
				if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
					logger.Errorf("Server %s-%s deregister error: %s", config.Name, config.NodeId, err)
				}
			}
		}

		s.Lock()
		swg := s.wg
		s.Unlock()

		// wait for requests to finish
		if swg != nil {
			swg.Wait()
		}
	}()
}

func (s *Core) Register() error {
	s.RLock()
	rsvc := s.rsvc
	s.RUnlock()
	config := s.Options()

	regFunc := func(service *registry.Service) error {
		// create registry options
		rOpts := []registry.RegisterOption{registry.RegisterTTL(10 * time.Second)}

		var regErr error

		for i := 0; i < 3; i++ {
			// attempt to register
			if err := config.Registry.Register(service, rOpts...); err != nil {
				// set the error
				regErr = err
				// backoff then retry
				time.Sleep(backoff.Do(i + 1))
				continue
			}
			// success so nil error
			regErr = nil
			break
		}

		return regErr
	}

	// have we registered before?
	if rsvc != nil {
		if err := regFunc(rsvc); err != nil {
			return err
		}
		return nil
	}

	var err error
	var advt, host string
	var cacheService bool

	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address

	advt = config.Address

	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, _, err = net.SplitHostPort(advt)
		if err != nil {
			return err
		}
	} else {
		host = advt
	}

	if ip := net.ParseIP(host); ip != nil {
		cacheService = true
	}

	addr, err := addr.Extract(host)
	if err != nil {
		return err
	}

	// make copy of metadata
	mapCopy := func(originalMap map[string]string) map[string]string {
		targetMap := make(map[string]string, len(originalMap))
		for key, value := range originalMap {
			targetMap[key] = value
		}
		return targetMap
	}
	md := mapCopy(config.Metadata)

	// register service
	node := &registry.Node{
		Id:       config.Name + "-" + config.NodeId,
		Address:  addr,
		Metadata: md,
	}

	node.Metadata["registry"] = config.Registry.String()
	node.Metadata["protocol"] = "mucp"

	s.RLock()

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	// get registered value
	registered := s.registered

	s.RUnlock()

	if !registered {
		if logger.V(logger.InfoLevel, logger.DefaultLogger) {
			logger.Infof("Registry [%s] Registering node: %s", config.Registry.String(), node.Id)
		}
	}

	// register the service
	if err := regFunc(service); err != nil {
		return err
	}

	// already registered? don't need to register subscribers
	if registered {
		return nil
	}

	s.Lock()
	defer s.Unlock()

	if cacheService {
		s.rsvc = service
	}
	s.registered = true

	return nil
}

func (s *Core) Deregister() error {
	var err error
	var advt, host string

	config := s.Options()

	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address

	advt = config.Address

	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, _, err = net.SplitHostPort(advt)
		if err != nil {
			return err
		}
	} else {
		host = advt
	}

	addr, err := addr.Extract(host)
	if err != nil {
		return err
	}

	node := &registry.Node{
		Id:      config.Name + "-" + config.NodeId,
		Address: addr,
	}

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	if logger.V(logger.InfoLevel, logger.DefaultLogger) {
		logger.Infof("Registry [%s] Deregistering node: %s", config.Registry.String(), node.Id)
	}

	if err := config.Registry.Deregister(service); err != nil {
		return err
	}

	s.Lock()
	s.rsvc = nil

	if !s.registered {
		s.Unlock()
		return nil
	}

	s.registered = false

	s.Unlock()
	return nil
}
