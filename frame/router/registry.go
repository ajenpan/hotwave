package router

import (
	"sync"
	"time"

	log "hotwave/logger"
	"hotwave/registry"
)

type Service struct {
	sync.RWMutex
	Name     string
	Version  string
	Metadata map[string]string
	Nodes    map[string]*Node
}

type Node struct {
	sync.RWMutex
	Id          string
	Address     string
	ServiceName string
	Metadata    map[string]string
	// RegistryAt  time.Time
}

func (s *Service) GetNode(id string) *Node {
	s.RLock()
	defer s.RUnlock()
	if node, has := s.Nodes[id]; has {
		return node
	}
	return nil
}

func (s *Service) NodeSize() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.Nodes)
}

// router is the default router
type registryRouter struct {
	exit chan bool
	opts Options
	rw   sync.RWMutex

	nodes    sync.Map // map of nodes
	services sync.Map // serverName -> Service
}

func (r *registryRouter) isClosed() bool {
	select {
	case <-r.exit:
		return true
	default:
		return false
	}
}

// refresh list of api services
func (r *registryRouter) refresh() {
	var attempts int

	t := time.NewTicker(10 * time.Minute)
	defer t.Stop()

	for {
		services, err := r.opts.Registry.ListServices()
		if err != nil {
			attempts++
			log.Errorf("unable to list services: %v", err)
			time.Sleep(time.Duration(attempts) * time.Second)
			continue
		}

		attempts = 0

		for _, s := range services {
			r.store(s)
		}

		// refresh list in 10 minutes... cruft
		// use registry watching
		select {
		case <-t.C:
		case <-r.exit:
			return
		}
	}
}

// process watch event
func (r *registryRouter) process(res *registry.Result) {
	if res == nil || res.Service == nil {
		return
	}

	switch res.Action {
	case registry.Delete.String():
		{
			r.remove(res.Service)
		}
	case registry.Update.String():
		fallthrough
	case registry.Create.String():
		{
			r.store(res.Service)
		}
	}
}

func (r *registryRouter) remove(service *registry.Service) {
	var s *Service
	if v, has := r.services.Load(service.Name); has {
		s = v.(*Service)
	}

	for _, node := range service.Nodes {
		r.nodes.Delete(node.Id)

		if s != nil {
			s.Lock()
			delete(s.Nodes, node.Id)
			s.Unlock()
		}
	}
}

// store local endpoint cache
func (r *registryRouter) store(service *registry.Service) {
	r.rw.Lock()
	defer r.rw.Unlock()

	var s *Service
	if v, has := r.services.Load(service.Name); has {
		s = v.(*Service)
	}

	if s == nil {
		s := &Service{
			Name:     service.Name,
			Version:  service.Version,
			Metadata: service.Metadata,
		}
		r.services.Store(service.Name, s)
	}

	for _, node := range service.Nodes {
		temp := &Node{
			Id:          node.Id,
			Address:     node.Address,
			ServiceName: service.Name,
			Metadata:    node.Metadata,
		}
		r.nodes.Store(node.Id, temp)

		s.Lock()
		s.Nodes[node.Id] = temp
		s.Unlock()
	}
}

// watch for endpoint changes
func (r *registryRouter) watch() {
	var attempts int

	for {
		if r.isClosed() {
			return
		}

		// watch for changes
		w, err := r.opts.Registry.Watch()
		if err != nil {
			attempts++
			log.Errorf("error watching endpoints: %v", err)
			time.Sleep(time.Duration(attempts) * time.Second)
			continue
		}

		ch := make(chan bool)

		go func() {
			select {
			case <-ch:
				w.Stop()
			case <-r.exit:
				w.Stop()
			}
		}()

		// reset if we get here
		attempts = 0

		for {
			// process next event
			res, err := w.Next()
			if err != nil {
				log.Errorf("error getting next endoint: %v", err)
				close(ch)
				break
			}
			r.process(res)
		}
	}
}

func (r *registryRouter) Options() Options {
	return r.opts
}

func (r *registryRouter) Close() error {
	select {
	case <-r.exit:
		return nil
	default:
		close(r.exit)
	}
	return nil
}
func (r *registryRouter) GetNode(id string) *Node {
	if v, has := r.nodes.Load(id); has {
		return v.(*Node)
	}
	return nil
}

func (r *registryRouter) GetService(name string) *Service {
	if v, has := r.services.Load(name); has {
		return v.(*Service)
	}
	return nil
}

func (r *registryRouter) Start() error {
	go r.watch()
	go r.refresh()
	return nil
}

// NewRouter returns the default router
func NewRouter(opts ...Option) *registryRouter {
	options := NewOptions(opts...)
	r := &registryRouter{
		exit: make(chan bool),
		opts: options,
	}
	return r
}
