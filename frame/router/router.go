package router

import (
	"hotwave/registry"
)

var DefaultRoute = NewRouter()

type Options struct {
	Handler  string
	Registry registry.Registry
	// Resolver resolver.Resolver
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Handler:  "meta",
		Registry: registry.DefaultRegistry,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
func WithRegistry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// func WithResolver(r resolver.Resolver) Option {
// 	return func(o *Options) {
// 		o.Resolver = r
// 	}
// }

type Router interface {
	Start() error
	Close() error
	GetService(name string) *Service
	GetNode(id string) *Node
}
