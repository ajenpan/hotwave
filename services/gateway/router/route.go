package router

import (
	"sync"
	"time"

	// "google.golang.org/protobuf/reflect/protodesc"
	// protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	// "google.golang.org/grpc/servicereflect"

	log "hotwave/logger"
	"hotwave/metadata"
	"hotwave/registry"
)

func New(r registry.Registry) *Router {
	ret := &Router{
		Registry: r,
		exit:     make(chan bool),
	}
	ret.Start()
	return ret
}

type Service struct {
	sync.RWMutex
	Name     string
	Version  string
	Metadata metadata.Metadata
	Nodes    map[string]*Node

	files descriptorpb.FileDescriptorSet
	// protoreflect.MethodDescriptors
}

type Node struct {
	sync.RWMutex
	Id          string
	Address     string
	ServiceName string
	Metadata    metadata.Metadata
}

type Router struct {
	registry.Registry

	exit chan bool

	nodeslock sync.RWMutex
	nodes     map[string]*Node // map of nodes

	serviceslock sync.RWMutex
	services     map[string]map[string]*Service //name-version-service
}

func (r *Router) Start() error {
	go r.watch()
	return nil
}

func (r *Router) GetNode(id string) *Node {
	// v, has := r.nodes.Load(id)
	// if !has {
	// 	return nil
	// }
	// return v.(*Node)
	return nil
}

func (r *Router) isClosed() bool {
	select {
	case <-r.exit:
		return true
	default:
		return false
	}
}

// watch for endpoint changes
func (r *Router) watch() {
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

// process watch event
func (r *Router) process(res *registry.Result) {
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

func (r *Router) remove(service *registry.Service) {
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

		// for _, node := range service.Nodes {
		// r.nodes.Delete(node.Id)
		// s.Lock()
		// delete(s.Nodes, node.Id)
		// s.Unlock()
		// }
	}
}

// store local endpoint cache
func (r *Router) store(rs *registry.Service) {
	// r.rw.Lock()
	// defer r.rw.Unlock()

	serverice := r.getService(rs.Name, rs.Version)
	if serverice != nil {
		serverice = warpRegService(rs)
	}
}

func warpRegService(service *registry.Service) *Service {
	return &Service{
		Name:     service.Name,
		Version:  service.Version,
		Metadata: metadata.Copy(service.Metadata),
	}
}

func (r *Router) getService(name, version string) *Service {
	r.serviceslock.RLock()
	defer r.serviceslock.RUnlock()
	s, ok := r.services[name]
	if !ok {
		return nil
	}
	return s[version]
}
