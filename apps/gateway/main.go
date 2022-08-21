package main

import (
	"crypto/rsa"
	"fmt"
	"net"
	"os"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"hotwave/event"
	evProto "hotwave/event/proto"
	log "hotwave/logger"
	gwAuth "hotwave/service/gateway/auth"
	gwHandler "hotwave/service/gateway/handler"
	gwProto "hotwave/service/gateway/proto"
	tcpGate "hotwave/transport/tcp"
	"hotwave/utils/rsagen"
	utilSignal "hotwave/utils/signal"
)

var Version string = "unknown"
var GitCommit string = "unknown"
var BuildAt string = "unknown"
var BuildBy string = "unknown"
var Name string = "allinone"

var ConfigPath string = ""
var ListenAddr string = ""
var PrintConf bool = false

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("project:", Name)
		fmt.Println("version:", Version)
		fmt.Println("git commit:", GitCommit)
		fmt.Println("build at:", BuildAt)
		fmt.Println("build by:", BuildBy)
	}

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       "config.yaml",
			Destination: &ConfigPath,
		}, &cli.StringFlag{
			Name:        "listen",
			Aliases:     []string{"l"},
			Value:       ":10110",
			Destination: &ListenAddr,
		}, &cli.BoolFlag{
			Name:        "print-config",
			Destination: &PrintConf,
			Hidden:      true,
		},
	}

	app.Action = RealMain

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}

func ReadRSAKey() (*rsa.PrivateKey, error) {
	const privateFile = "private.pem"
	const publicFile = "public.pem"

	raw, err := os.ReadFile(privateFile)
	if err != nil {
		privateKey, publicKey, err := rsagen.GenerateRsaPem(2048)
		if err != nil {
			return nil, err
		}
		raw = []byte(privateKey)
		os.WriteFile(privateFile, []byte(privateKey), 0644)
		os.WriteFile(publicFile, []byte(publicKey), 0644)
	}
	return rsagen.ParseRsaPrivateKeyFromPem(raw)
}

var PK *rsa.PrivateKey

func RealMain(c *cli.Context) error {
	var err error
	PK, err = ReadRSAKey()
	if err != nil {
		panic(err)
	}

	grpcs := grpc.NewServer()
	publisher := &event.GrpcEventPublisher{}
	evProto.RegisterEventServer(grpcs, publisher)

	gw, err := gwHandler.NewGateway(gwHandler.GatewayOption{
		AuthClient: &gwAuth.LocalAuth{
			PK: &PK.PublicKey,
		},
		Publisher: publisher,
	})

	if err != nil {
		panic(err)
	}

	gw.AddAllowlistByMsg(&gwProto.LoginGateRequest{})
	gwProto.RegisterGatewayServer(grpcs, gw)

	lis, err := net.Listen("tcp", ":20000")
	if err != nil {
		panic(err)
	}

	go func() {
		err = grpcs.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	defer grpcs.Stop()

	gate := tcpGate.NewServer(tcpGate.ServerOptions{
		Address:   ":10000",
		OnMessage: gw.OnGateMessage,
		OnConn:    gw.OnGateConnStat,
	})

	if err := gate.Start(); err != nil {
		panic(err)
	}
	defer gate.Stop()

	signal := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", signal.String())
	return nil
}
