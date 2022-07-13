package main

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	protobuf "google.golang.org/protobuf/proto"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	log "hotwave/logger"
	"hotwave/marshal"
	"hotwave/service/auth/handler"
	"hotwave/service/auth/proto"
	"hotwave/service/auth/store/cache"
	gwclient "hotwave/service/gateway/client"
	gwproto "hotwave/service/gateway/proto"
	"hotwave/transport"
	"hotwave/utils/calltable"
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

func createMysqlClient(dsn string) *gorm.DB {
	dbc, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableNestedTransaction: true, //关闭嵌套事务
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return dbc
}

func createSQLiteClient(dsn string) *gorm.DB {
	dbc, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return dbc
}

func RealMain(c *cli.Context) error {
	pk, err := ReadRSAKey()
	if err != nil {
		panic(err)
	}
	authHandler := handler.NewAuth(handler.AuthOptions{
		PK:    pk,
		DB:    createMysqlClient("root:123456@tcp(localhost:3306)/auth?charset=utf8mb4&parseTime=True&loc=Local"),
		Cache: cache.NewMemory(),
	})

	ct := calltable.ExtractParseGRpcMethod(proto.File_service_auth_proto_auth_proto.Services(), authHandler)
	authHandler.CT = ct

	ServerCallTable(http.DefaultServeMux, authHandler, ct)
	fmt.Println("start http server at: ", ListenAddr)

	go func() {
		http.ListenAndServe(ListenAddr, nil)
	}()

	grpcConn, err := grpc.Dial("localhost:20000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	gwc := &gwclient.GRPCClient{
		GrpcConn: grpcConn,
		NodeID:   "auth-1",
		NodeName: "auth",
		OnConnStatusFunc: func(c *gwclient.GRPCClient, ss transport.SessionStat) {
			if ss == transport.Connected {
				c.GetewayClient.AddGateAllowList(context.Background(), &gwproto.AddGateAllowListRequest{
					Topics: []string{
						string(protobuf.MessageName(&proto.LoginRequest{})),
						string(protobuf.MessageName(&proto.CaptchaRequest{})),
						string(protobuf.MessageName(&proto.LogoutRequest{})),
					},
				})
			} else {
				go c.Reconnect()
			}
		},
		OnUserMessageFunc: authHandler.OnUserMessage,
	}

	gwc.Reconnect()
	defer gwc.Close()

	authHandler.Client = gwc
	signal := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", signal.String())
	return nil
}

func ServerCallTable(mux *http.ServeMux, handler interface{}, ct *calltable.CallTable) {
	respWithError := func(w http.ResponseWriter, data interface{}, err error) {
		type HttpRespType struct {
			Data    interface{} `json:"data"`
			Code    int         `json:"code"`
			Message string      `json:"message"`
		}
		respWrap := &HttpRespType{
			Data:    data,
			Message: "ok",
		}
		if err != nil {
			respWrap.Code = -1
			respWrap.Message = err.Error()
		}
		raw, _ := json.Marshal(respWrap)
		w.Write(raw)
	}

	ct.Range(func(key string, method *calltable.Method) bool {
		pattern := "/" + key
		fmt.Println("handle: ", pattern)
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			raw, err := ioutil.ReadAll(r.Body)

			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			if err != nil {
				respWithError(w, nil, fmt.Errorf("read body error: %s", err.Error()))
				return
			}

			req := reflect.New(method.RequestType).Interface().(protobuf.Message)

			// todo : get marshaler
			jsonpb := &marshal.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
					UseProtoNames:   true,
				},
			}

			if err := jsonpb.Unmarshal(raw, req); err != nil {
				respWithError(w, nil, fmt.Errorf("unmarshal request error: %s", err.Error()))
				return
			}

			// here call method
			respArgs := method.Call(handler, r.Context(), req)

			if len(respArgs) != 2 {
				return
			}

			var respErr error
			if !respArgs[1].IsNil() {
				respErr = respArgs[1].Interface().(error)
			}

			var respData json.RawMessage
			if !respArgs[0].IsNil() {
				if resp, ok := respArgs[0].Interface().(protobuf.Message); ok {
					data, err := jsonpb.Marshal(resp)
					if err == nil {
						respData = data
					} else {
						respWithError(w, nil, respErr)
					}
				}
			}

			respWithError(w, respData, respErr)
		})
		return true
	})
}
