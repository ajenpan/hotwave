package client

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"

	"hotwave/registry"
)

func init() {
	encoding.RegisterCodec(wrapCodec{jsonCodec{}})
	encoding.RegisterCodec(wrapCodec{protoCodec{}})
	encoding.RegisterCodec(wrapCodec{bytesCodec{}})
}

type grpcClient struct {
	opts Options
	pool *pool
	once atomic.Value
}

// secure returns the dial option for whether its a secure or insecure connection
func (g *grpcClient) secure(addr string) grpc.DialOption {
	// first we check if theres'a  tls config
	// if g.opts.Context != nil {
	// 	if v := g.opts.Context.Value(tlsAuth{}); v != nil {
	// 		tls := v.(*tls.Config)
	// 		creds := credentials.NewTLS(tls)
	// 		// return tls config if it exists
	// 		return grpc.WithTransportCredentials(creds)
	// 	}
	// }

	// default config
	tlsConfig := &tls.Config{}
	defaultCreds := grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))

	// check if the address is prepended with https
	if strings.HasPrefix(addr, "https://") {
		return defaultCreds
	}

	// if no port is specified or port is 443 default to tls
	_, port, err := net.SplitHostPort(addr)
	// assuming with no port its going to be secured
	if port == "443" {
		return defaultCreds
	} else if err != nil && strings.Contains(err.Error(), "missing port in address") {
		return defaultCreds
	}

	// other fallback to insecure
	return grpc.WithInsecure()
}

func (g *grpcClient) call(ctx context.Context, node *registry.Node, req interface{}, rsp interface{}, opts CallOptions) error {

	return nil
}

func (g *grpcClient) Init(opts ...Option) error {
	size := g.opts.PoolSize
	ttl := g.opts.PoolTTL

	for _, o := range opts {
		o(&g.opts)
	}

	// update pool configuration if the options changed
	if size != g.opts.PoolSize || ttl != g.opts.PoolTTL {
		g.pool.Lock()
		g.pool.size = g.opts.PoolSize
		g.pool.ttl = int64(g.opts.PoolTTL.Seconds())
		g.pool.Unlock()
	}
	return nil
}

func (g *grpcClient) Options() Options {
	return g.opts
}

func (g *grpcClient) NewMessage(topic string, msg interface{}, opts ...MessageOption) Message {
	// return newGRPCEvent(topic, msg, g.opts.ContentType, opts...)
	return nil
}
func (g *grpcClient) NewRequest(service, method string, req interface{}, reqOpts ...RequestOption) Request {
	// return newGRPCRequest(service, method, req, g.opts.ContentType, reqOpts...)
	return nil
}

func (g *grpcClient) Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error {

	return nil
}

func (g *grpcClient) Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
	return nil, nil
}

func (g *grpcClient) Publish(ctx context.Context, p Message, opts ...PublishOption) error {
	return nil
}

func (g *grpcClient) String() string {
	return "grpc"
}

func NewClient(opts ...Option) Client {
	// return newClient(opts...)
	options := NewOptions()
	// default content type for grpc
	options.ContentType = "application/grpc+proto"

	for _, o := range opts {
		o(&options)
	}

	rc := &grpcClient{
		opts: options,
	}
	rc.once.Store(false)

	rc.pool = newPool(options.PoolSize, options.PoolTTL, 50, 20)

	c := Client(rc)

	// wrap in reverse
	for i := len(options.Wrappers); i > 0; i-- {
		c = options.Wrappers[i-1](c)
	}
	return c
}
