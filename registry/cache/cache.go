package cache

import (
	"math/rand"
	"sync"
	"time"

	log "hotwave/logger"
	"hotwave/metadata"
	"hotwave/registry"
)

type Cache interface {
	Stop() error
	GetNode(id string) *Node
	GetService(string, ...registry.GetOption) ([]*registry.Service, error)
}

type Service struct {
	sync.RWMutex
	Name     string
	Version  string
	Metadata metadata.Metadata
	Nodes    map[string]*Node
}

type Node struct {
	sync.RWMutex
	Id          string
	Address     string
	ServiceName string
	Metadata    metadata.Metadata
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
	registry.Registry

	exit chan bool
	opts Options

	nodes sync.Map // map of nodes

	serviceslock sync.RWMutex
	services     map[string]map[string]*Service //name-version-service
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
		services, err := r.ListServices()
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
	r.serviceslock.RLock()
	defer r.serviceslock.RUnlock()

	ss, has := r.services[service.Name]

	if !has || len(ss) == 0 {
		return
	}

	for _, s := range ss {
		if s.Version != service.Version {
			continue
		}

		for _, node := range service.Nodes {
			r.nodes.Delete(node.Id)

			s.Lock()
			delete(s.Nodes, node.Id)
			s.Unlock()
		}
	}
}

// store local endpoint cache
func (r *registryRouter) store(service *registry.Service) {
	// r.rw.Lock()
	// defer r.rw.Unlock()

	// r.serviceslock.RLock()
	// defer r.serviceslock.RUnlock()

	// ss, has := r.services[service.Name]
	// if !has {
	// 	ss:=make(map[string]*Service)
	// }

	// s, has := ss[service.Version]
	// if !has {
	// 	s = &Service{
	// 		Name:     service.Name,
	// 		Version:  service.Version,
	// 		Metadata: service.Metadata,
	// 		Nodes:    make(map[string]*Node),
	// 	}
	// 	ss[service.Version] = s
	// 	r.services[service.Name] = ss
	// }

	// if s == nil {
	// 	s := &Service{
	// 		Name:     service.Name,
	// 		Version:  service.Version,
	// 		Metadata: metadata.Copy(service.Metadata),
	// 	}
	// 	r.services.Store(service.Name, s)
	// }

	// for _, node := range service.Nodes {
	// 	temp := &Node{
	// 		Id:          node.Id,
	// 		Address:     node.Address,
	// 		ServiceName: service.Name,
	// 		Metadata:    metadata.Copy(node.Metadata),
	// 	}
	// 	r.nodes.Store(node.Id, temp)

	// 	s.Lock()
	// 	s.Nodes[node.Id] = temp
	// 	s.Unlock()
	// }
}

// watch for endpoint changes
func (r *registryRouter) watch() {
	var attempts int

	for {
		if r.isClosed() {
			return
		}

		// watch for changes
		w, err := r.Watch()
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

func (r *registryRouter) Stop() error {
	select {
	case <-r.exit:
		return nil
	default:
		close(r.exit)
	}
	return nil
}

func (r *registryRouter) GetNode(id string) *Node {
	v, has := r.nodes.Load(id)
	if !has {
		return nil
	}
	return v.(*Node)
}

func ConvertNode(node *Node) *registry.Node {
	node.RLock()
	defer node.RUnlock()

	ret := &registry.Node{
		Id:       node.Id,
		Address:  node.Address,
		Metadata: metadata.Copy(node.Metadata),
	}
	return ret
}

func ConvertService(s *Service) *registry.Service {
	s.RLock()
	defer s.RUnlock()

	ret := &registry.Service{
		Name:     s.Name,
		Version:  s.Version,
		Metadata: s.Metadata,
	}
	for _, node := range s.Nodes {
		ret.Nodes = append(ret.Nodes, ConvertNode(node))
	}
	return ret
}

func (r *registryRouter) GetService(service string, opts ...registry.GetOption) ([]*registry.Service, error) {
	//TODO:
	// v, has := r.services.Load(service)
	// if !has {
	// 	return nil, registry.ErrNotFound
	// }
	// s := v.(*Service)
	// return ConvertService(s), nil
	return nil, nil
}

func (r *registryRouter) Start() error {
	go r.watch()
	go r.refresh()
	return nil
}

func New(r registry.Registry, opts ...Option) *registryRouter {
	rand.Seed(time.Now().UnixNano())
	options := Options{
		TTL: DefaultTTL,
	}

	for _, o := range opts {
		o(&options)
	}

	ret := &registryRouter{
		Registry: r,
		exit:     make(chan bool),
		opts:     options,
	}
	ret.Start()
	return ret
}
