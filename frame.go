package hotwave

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"hotwave/logger"
	"hotwave/metadata"
	"hotwave/registry"
	utilAddr "hotwave/util/addr"
	"hotwave/util/backoff"
)

func New(opts ...Option) *Frame {
	ret := &Frame{
		exit: make(chan chan error),
	}

	options := newOptions(opts...)
	for _, opt := range opts {
		opt(&options)
	}

	if options.NodeId == "" {
		options.NodeId = uuid.Must(uuid.NewUUID()).String()
	}

	ret.opts = options

	ret.grpcServer = grpc.NewServer()

	return ret
}

type Frame struct {
	// protocol.UnimplementedNodeBaseServer

	sync.RWMutex
	// marks the serve as started
	started bool
	// used for first registration
	registered bool
	rsvc       *registry.Service
	opts       Options
	grpcServer *grpc.Server

	exit chan chan error
}

// export for grpc handler
func (f *Frame) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	f.grpcServer.RegisterService(desc, impl)
}

func (f *Frame) Start() error {
	f.RLock()
	if f.started {
		f.RUnlock()
		return nil
	}
	f.RUnlock()

	config := f.Options()

	lis, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}
	logger.Infof("Server %s-%s started", config.Name, config.NodeId)

	recverr := make(chan error, 1)
	go func() {
		defer close(recverr)

		//update for real address when the port is random, eg. localhost:0
		addr := f.opts.Address
		f.opts.Address = lis.Addr().String()
		defer func() {
			//set back the address
			f.opts.Address = addr
		}()

		err := f.grpcServer.Serve(lis)
		if err != nil {
			logger.Errorf("grpc server error: %s", err)
			recverr <- err
		}
	}()

	select {
	case err = <-recverr:
		//wait 1s for error
	case <-time.After(time.Second * 1):
	}

	if err != nil {
		return err
	}

	// register self to the world
	if err := f.register(); err != nil {
		return err
	}
	go f.keepRegistered()

	f.Lock()
	f.started = true
	f.Unlock()
	return nil
}

func (f *Frame) Stop() error {
	f.RLock()
	if !f.started {
		f.RUnlock()
		return nil
	}
	f.RUnlock()

	recvChan := make(chan error)
	f.exit <- recvChan
	f.grpcServer.Stop()

	f.Lock()
	f.rsvc = nil
	f.started = false
	f.Unlock()
	return <-recvChan
}

func (f *Frame) Options() Options {
	f.RLock()
	opts := f.opts
	f.RUnlock()
	return opts
}

func (f *Frame) keepRegistered() {
	config := f.Options()
	t := new(time.Ticker)
	// only process if it exists
	if f.opts.RegisterInterval > time.Duration(0) {
		// new ticker
		t = time.NewTicker(f.opts.RegisterInterval)
	}

	var ch chan error

	func() {
		var err error

		for {
			select {
			// register self on interval
			case <-t.C:
				f.RLock()
				registered := f.registered
				f.RUnlock()
				rerr := f.opts.RegisterCheck(context.Background())
				if rerr != nil && registered {
					logger.Errorf("Server %s-%s register check error: %s, deregister it", config.Name, config.NodeId, err)
					// deregister self in case of error
					if err := f.deregister(); err != nil {
						logger.Errorf("Server %s-%s deregister error: %s", config.Name, config.NodeId, err)
					}
				} else if rerr != nil && !registered {
					logger.Errorf("Server %s-%s register check error: %s", config.Name, config.NodeId, err)
					continue
				}
				if err := f.register(); err != nil {
					logger.Errorf("Server %s-%s register error: %s", config.Name, config.NodeId, err)
				}
			// wait for exit
			case ch = <-f.exit:
				t.Stop()
				return
			}
		}
	}()

	f.RLock()
	registered := f.registered
	f.RUnlock()
	if registered {
		// deregister self
		err := f.deregister()
		if err != nil {
			logger.Errorf("Server %s-%s deregister error: %s", config.Name, config.NodeId, err)
			ch <- err
			return
		}
	}
	ch <- nil
}

func (f *Frame) register() error {
	f.RLock()
	rsvc := f.rsvc
	f.RUnlock()
	config := f.Options()

	regFunc := func(service *registry.Service) error {
		// create registry options
		rOpts := []registry.RegisterOption{registry.RegisterTTL(config.RegisterTTL)}

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

	var cacheService bool
	var err error
	var advt, host, port string

	advt = config.Address
	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, port, err = net.SplitHostPort(advt)
		if err != nil {
			return err
		}
	} else {
		host = advt
	}

	if ip := net.ParseIP(host); ip != nil {
		cacheService = true
	}

	addr, err := utilAddr.Extract(host)
	if err != nil {
		return err
	}

	if port != "" {
		addr = utilAddr.HostPort(addr, port)
	}

	// register service
	node := &registry.Node{
		Id:       config.Name + "." + config.NodeId,
		Address:  addr,
		Metadata: metadata.Copy(config.Metadata),
	}

	node.Metadata["server"] = "mucp"
	node.Metadata["registry"] = config.Registry.String()
	node.Metadata["protocol"] = "mucp"

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	f.RLock()
	// get registered value
	registered := f.registered
	f.RUnlock()

	if !registered {
		if logger.V(logger.InfoLevel, logger.DefaultLogger) {
			logger.Infof("[%s] Registering node:%s, address:%s", config.Registry.String(), node.Id, node.Address)
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

	f.Lock()
	defer f.Unlock()

	if cacheService {
		f.rsvc = service
	}
	f.registered = true

	return nil
}

func (f *Frame) deregister() error {
	var err error
	var advt, host, port string

	config := f.Options()

	advt = config.Address
	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, port, err = net.SplitHostPort(advt)
		if err != nil {
			return err
		}
	} else {
		host = advt
	}

	addr, err := utilAddr.Extract(host)
	if err != nil {
		return err
	}

	if port != "" {
		addr = utilAddr.HostPort(addr, port)
	}

	node := &registry.Node{
		Id:      config.Name + "." + config.NodeId,
		Address: addr,
	}

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	logger.Infof("Registry [%s] Deregistering node: %s", config.Registry.String(), node.Id)

	if err := config.Registry.Deregister(service); err != nil {
		return err
	}

	f.Lock()
	f.rsvc = nil

	if !f.registered {
		f.Unlock()
		return nil
	}

	f.registered = false

	f.Unlock()
	return nil
}
