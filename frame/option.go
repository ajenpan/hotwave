package frame

import (
	"context"
	"time"

	"hotwave/registry"
)

type Options struct {
	Name     string            //required
	NodeId   string            //option
	Version  string            //required
	Address  string            //option
	Metadata map[string]string //option
	Registry registry.Registry //option

	// The register expiry time
	RegisterTTL time.Duration
	// The interval on which to register
	RegisterInterval time.Duration
	// RegisterCheck runs a check function before registering the service
	RegisterCheck func(context.Context) error
}

func newOptions(opts ...Option) Options {
	opt := DefaultOptions
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

type Option func(*Options)

var DefaultOptions = Options{
	Name:             "unknown.service",
	Version:          "latest",
	NodeId:           "",
	Registry:         registry.DefaultRegistry,
	Address:          ":0",
	RegisterCheck:    func(context.Context) error { return nil },
	RegisterInterval: time.Second * 30,
	RegisterTTL:      time.Second * 90,
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}
