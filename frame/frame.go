package frame

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	protocol "hotwave/frame/proto"
	"hotwave/frame/router"
	"hotwave/logger"
	"hotwave/registry"
	utilAddr "hotwave/util/addr"
	"hotwave/util/backoff"
)

func NewFrame(opts ...Option) *frame {
	ret := &frame{
		exit: make(chan chan error),
	}

	options := newOptions(opts...)

	if options.NodeId == "" {
		options.NodeId = uuid.Must(uuid.NewUUID()).String()
	}
	ret.opts = options

	grpcServer := grpc.NewServer()
	protocol.RegisterNodeBaseServer(grpcServer, ret)
	ret.grpcServer = grpcServer

	ret.router = router.NewRouter(router.WithRegistry(options.Registry))
	return ret
}

type frame struct {
	protocol.UnimplementedNodeBaseServer

	sync.RWMutex
	// marks the serve as started
	started bool
	// used for first registration
	registered bool
	rsvc       *registry.Service
	opts       Options
	grpcServer *grpc.Server

	exit   chan chan error
	router router.Router
}

func (f *frame) GetService(name string) *router.Service {
	return f.router.GetService(name)
}

// export for grpc handler
func (f *frame) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	f.grpcServer.RegisterService(desc, impl)
}

func (f *frame) Stop() error {
	recvChan := make(chan error)
	f.exit <- recvChan
	f.grpcServer.Stop()
	f.router.Close()
	return <-recvChan
}

func (f *frame) Start() error {
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

	//TODO: how to hold the error?
	go func() {
		err := f.grpcServer.Serve(lis)
		if err != nil {
			logger.Errorf("grpc server error: %s", err)
		}
	}()

	//update for real address when the port is random, eg. localhost:0
	addr := f.opts.Address
	f.opts.Address = lis.Addr().String()
	logger.Info(addr)

	// register self to the world
	if err := f.register(); err != nil {
		return err
	}

	go f.keepRegistered()

	if err := f.router.Start(); err != nil {
		return err
	}

	f.Lock()
	f.started = true
	f.Unlock()
	return nil
}

func (f *frame) Options() Options {
	f.RLock()
	opts := f.opts
	f.RUnlock()
	return opts
}

func (f *frame) keepRegistered() {
	config := f.Options()
	t := new(time.Ticker)
	// only process if it exists
	if f.opts.RegisterInterval > time.Duration(0) {
		// new ticker
		t = time.NewTicker(f.opts.RegisterInterval)
	}
	// return error chan
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

func (f *frame) register() error {
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

	var err error
	var advt, host, port string
	var cacheService bool

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

	// make copy of metadata
	mapCopy := func(originalMap map[string]string) map[string]string {
		targetMap := make(map[string]string, len(originalMap))
		for key, value := range originalMap {
			targetMap[key] = value
		}
		return targetMap
	}
	md := mapCopy(config.Metadata)

	if port != "" {
		addr = utilAddr.HostPort(addr, port)
	}

	// register service
	node := &registry.Node{
		Id:       config.Name + "-" + config.NodeId,
		Address:  addr,
		Metadata: md,
	}

	node.Metadata["server"] = "mucp"
	node.Metadata["registry"] = config.Registry.String()
	node.Metadata["protocol"] = "mucp"

	f.RLock()

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	// get registered value
	registered := f.registered

	f.RUnlock()

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

	f.Lock()
	defer f.Unlock()

	if cacheService {
		f.rsvc = service
	}
	f.registered = true

	return nil
}

func (f *frame) deregister() error {
	var err error
	var advt, host string

	config := f.Options()

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

	addr, err := utilAddr.Extract(host)
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
