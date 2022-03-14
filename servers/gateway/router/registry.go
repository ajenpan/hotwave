package router

import (
	"strings"
	"sync"
	"time"

	log "hotwave/logger"
	"hotwave/registry"
	"hotwave/registry/cache"
)

type Info struct {
	endpoint string
	svrname  string
}

type endpointName = string
type serviceName = string
type nodeName = string

type ServiceMap = map[serviceName]*Service
type NodeMap = map[nodeName]*registry.Node
type EndpointMap = map[endpointName]*registry.Endpoint

type Service struct {
	rw        sync.RWMutex
	Name      string
	Version   string
	Nodes     NodeMap
	Endpoints EndpointMap
}

func (s *Service) GetNode(id string) *registry.Node {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if node, has := s.Nodes[id]; has {
		return node
	}
	return nil
}

func (s *Service) NodeSize() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.Nodes)
}

func serviceConvert(s *registry.Service) *Service {
	nodeConvert := func(nodes []*registry.Node) NodeMap {
		ret := NodeMap{}
		for _, node := range nodes {
			ret[node.Id] = node
		}
		return ret
	}
	eptConvert := func(epts []*registry.Endpoint) EndpointMap {
		ret := EndpointMap{}
		for _, ept := range epts {
			ret[ept.Name] = ept
		}
		return ret
	}

	return &Service{
		Name:      s.Name,
		Version:   s.Version,
		Nodes:     nodeConvert(s.Nodes),
		Endpoints: eptConvert(s.Endpoints),
	}
}

// router is the default router
type registryRouter struct {
	exit chan bool
	opts Options

	// registry cache
	rc cache.Cache

	//lockSvrAddress sync.RWMutex
	//svrAddress     map[string]*NodeAddress

	lockEps2svr sync.RWMutex
	eps2svr     map[endpointName]*Info //endpoint:info

	//eps2service map[endpointName]ServiceMap
	// req2svr map[string]ServiceMap

	lockServices sync.RWMutex
	services     ServiceMap
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

		// for each gService, get gService and store endpoints
		for _, s := range services {
			service, err := r.rc.GetService(s.Name)
			if err != nil {
				log.Errorf("unable to get gService: %v", err)
				continue
			}
			r.store(service)
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
	// skip these things
	if res == nil || res.Service == nil {
		return
	}
	if res.Action == "delete" {
		r.remove(res.Service.Name)
		return
	}
	// get entry from cache
	service, err := r.rc.GetService(res.Service.Name)
	if err != nil {
		log.Errorf("unable to get gService: %v", err)
		return
	}
	// update our local endpoints
	r.store(service)
}

func (r *registryRouter) remove(serviceName string) {
	r.lockEps2svr.Lock()
	defer r.lockEps2svr.Unlock()

	for key, value := range r.eps2svr {
		if value.svrname != serviceName {
			continue
		}
		delete(r.eps2svr, key)
	}
}

// store local endpoint cache
func (r *registryRouter) store(services []*registry.Service) {

	// services
	eps2svr := map[string]*Info{}
	smap := ServiceMap{}
	//request2services := map[string]ServiceMap{}

	// create a new endpoint mapping
	for _, service := range services {
		// set names we need later
		smap[service.Name] = serviceConvert(service)

		//map per endpoint
		for _, ept := range service.Endpoints {
			eptName := strings.ToLower(ept.Name)
			if _, has := eps2svr[eptName]; has {
				log.Warnf("exist endpoint name :%s", eptName)
				continue
			}
			eps2svr[eptName] = &Info{
				endpoint: ept.Name,
				svrname:  service.Name,
			}

			//if s, has := request2services[ept.Request.Type]; !has {
			//	s = ServiceMap{}
			//	request2services[ept.Request.Type] = smap
			//	s[service.Name] = smap[service.Name]
			//} else {
			//	s[service.Name] = smap[service.Name]
			//}
		}
	}

	//reflashRequest:=func(){
	//	r.lockEps2svr.Lock()
	//	defer r.lockEps2svr.Unlock()
	//	//删除已经存在的
	//	for key, value := range r.req2svr {
	//
	//		if _, has := smap[value.svrname]; !has {
	//			continue
	//		}
	//
	//		delete(r.eps2svr, key)
	//	}
	//	//重新刷新
	//	for key, value := range eps2svr {
	//		r.eps2svr[key] = value
	//	}
	//}
	//reflashRequest()

	reflashEps2svr := func() {
		r.lockEps2svr.Lock()
		defer r.lockEps2svr.Unlock()
		//删除已经存在的
		for key, value := range r.eps2svr {
			if _, has := smap[value.svrname]; !has {
				continue
			}
			delete(r.eps2svr, key)
		}
		//重新刷新
		for key, value := range eps2svr {
			r.eps2svr[key] = value
		}
	}
	reflashEps2svr()

	reflashServices := func() {
		r.lockServices.Lock()
		defer r.lockServices.Unlock()
		for k, v := range smap {
			r.services[k] = v
		}
	}
	reflashServices()
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
		r.rc.Stop()
	}
	return nil
}

func (r *registryRouter) GetService(name string) *Service {
	r.lockEps2svr.RLock()
	defer r.lockEps2svr.RUnlock()

	if s, has := r.services[name]; has {
		return s
	}
	return nil
}

func (r *registryRouter) GetServiceInfo(path string) *Info {
	if r.isClosed() {
		return nil
	}
	path = strings.ToLower(path)
	r.lockEps2svr.RLock()
	defer r.lockEps2svr.RUnlock()

	ret, ok := r.eps2svr[path]
	if !ok {
		return nil
	}
	return ret
}

func (r *registryRouter) Start() {
	go r.watch()
	go r.refresh()
}

// NewRouter returns the default router
func NewRouter(opts ...Option) *registryRouter {
	options := NewOptions(opts...)
	r := &registryRouter{
		exit:     make(chan bool),
		opts:     options,
		rc:       cache.New(options.Registry),
		eps2svr:  make(map[string]*Info),
		services: ServiceMap{},
	}
	return r
}
