package router

import (
	"hotwave/registry"
)

type Router struct{}

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
