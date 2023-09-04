package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	log "hotwave/logger"
	"hotwave/service/auth/handler"
	"hotwave/service/auth/proto"
	"hotwave/service/auth/store/cache"
	"hotwave/service/auth/store/models"
	"hotwave/transport/httpsvr"
	"hotwave/utils/calltable"
	"hotwave/utils/marshal"
	"hotwave/utils/rsagen"
	utilSignal "hotwave/utils/signal"
)

var Version string = "unknown"
var GitCommit string = "unknown"
var BuildAt string = "unknown"
var BuildBy string = "unknown"
var Name string = "auth"

var ConfigPath string = ""
var ListenAddr string = ""
var PrintConf bool = false

const PrivateKeyFile = "private.pem"
const PublicKeyFile = "public.pem"

func ReadRSAKey() ([]byte, []byte, error) {
	privateRaw, err := os.ReadFile(PrivateKeyFile)
	if err != nil {
		privateKey, publicKey, err := rsagen.GenerateRsaPem(2048)
		if err != nil {
			return nil, nil, err
		}
		privateRaw = []byte(privateKey)
		os.WriteFile(PrivateKeyFile, []byte(privateKey), 0644)
		os.WriteFile(PublicKeyFile, []byte(publicKey), 0644)
	}
	publicRaw, err := os.ReadFile(PublicKeyFile)
	if err != nil {
		return nil, nil, err
	}
	return privateRaw, publicRaw, nil
}

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
			Value:       ":30020",
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

func CreateMysqlClient(dsn string) *gorm.DB {
	dbc, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableNestedTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return dbc
}

func CreateSQLiteClient(dsn string) *gorm.DB {
	dbc, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dbc.AutoMigrate(
		&models.Users{},
	)
	if err != nil {
		panic(err)
	}
	return dbc
}

func RealMain(c *cli.Context) error {
	privateRaw, publicRaw, err := ReadRSAKey()
	if err != nil {
		panic(err)
	}
	pk, err := rsagen.ParseRsaPrivateKeyFromPem(privateRaw)
	if err != nil {
		return err
	}

	h := handler.NewAuth(handler.AuthOptions{
		PK:        pk,
		PublicKey: publicRaw,
		DB:        CreateMysqlClient("root:123456@tcp(test122:13306)/auth?charset=utf8mb4&parseTime=True&loc=Local"),
		Cache:     cache.NewMemory(),
	})

	ct := calltable.ExtractParseGRpcMethod(proto.File_proto_auth_proto.Services(), h)
	svr := &httpsvr.HttpSvr{
		Mux: http.NewServeMux(),
		Marshal: &marshal.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseEnumNumbers: true,
				UseProtoNames:  true,
			},
		},
		Addr: ListenAddr,
	}
	svr.ServerCallTable(ct)

	go svr.Start()
	defer svr.Stop()

	fmt.Println("start http server at:", ListenAddr)

	signal := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", signal.String())
	return nil
}

func AuthMiddleWrap(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}
