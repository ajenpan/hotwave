package hotwave

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	protobuf "google.golang.org/protobuf/proto"

	"hotwave/logger"
	"hotwave/proto"
)

type Options struct {
	Name    string  //required
	NodeId  string  //option, if empty, will be generated
	Version string  //required
	Address string  //option
	Adpater Adpater //option

	// The register expiry time
	RegisterTTL time.Duration
	// The interval on which to register
	RegisterInterval time.Duration
	// RegisterCheck runs a check function before registering the service
	RegisterCheck func(context.Context) error
}

func NewOptions(opts ...Option) Options {
	opt := DefaultOptions
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

type Option func(*Options)

type Adpater interface {
	OnUserMessage(User, *proto.UserMessageWraper)
	OnNodeEvent(string, *proto.EventMessageWraper)
}

var DefaultOptions = Options{
	Name:             "unknown.service",
	Version:          "latest",
	NodeId:           "",
	Address:          ":0",
	Adpater:          &NoopAdpater{},
	RegisterCheck:    func(context.Context) error { return nil },
	RegisterInterval: time.Second * 30,
	RegisterTTL:      time.Second * 90,
}

func New(opts Options) *Core {
	ret := &Core{
		exit: make(chan chan error),
		// wg:          wait(.Context),
	}

	if opts.NodeId == "" {
		opts.NodeId = uuid.Must(uuid.NewUUID()).String()
	}
	ret.opts = opts

	grpcServer := grpc.NewServer()
	ret.grpcServer = grpcServer

	return ret
}

type Core struct {
	// protocol.UnimplementedNodeBaseServer
	rwLock sync.RWMutex
	// marks the serve as started
	started bool
	// used for first registration
	wg *sync.WaitGroup
	// rsvc       *registry.Service
	opts       Options
	grpcServer *grpc.Server
	exit       chan chan error
}

func (c *Core) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	c.grpcServer.RegisterService(desc, impl)
}

func (c *Core) Stop() error {
	c.grpcServer.Stop()
	return nil
}

func (c *Core) Start() error {

	config := c.Options()

	lis, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}
	go c.grpcServer.Serve(lis)

	//update for real address when the port is random, eg. localhost:0
	addr := c.opts.Address
	c.opts.Address = lis.Addr().String()
	logger.Info(addr)

	c.started = true
	return nil
}

func (c *Core) SendMessageToUser(user User, msg protobuf.Message) error {

	return nil
}

func (s *Core) Options() Options {
	opts := s.opts
	return opts
}
