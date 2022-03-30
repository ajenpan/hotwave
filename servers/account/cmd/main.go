package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
	protobuf "google.golang.org/protobuf/proto"

	frame "hotwave"
	"hotwave/logger"
	"hotwave/servers/account/config"
	"hotwave/servers/account/handler"
	"hotwave/servers/account/proto"
	utilHandle "hotwave/util/handle"
	utilSignal "hotwave/util/signal"
	// "google.golang.org/protobuf/encoding/protojson"
	// "google.golang.org/protobuf/reflect/protoreflect"
	// "hotwave/servers/account/handler"
	// proto "hotwave/account/proto"
	// handlerHelper "hotwave/util/handler"
)

var Name = "account"
var Version string = "unknow"
var GitCommit string = "unknow"
var BuildAt string = "unknow"
var BuildBy string = "unknow"

func main() {
	err := Run()
	if err != nil {
		fmt.Println(err)
	}
}

func Run() error {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("project:", Name)
		fmt.Println("version:", Version)
		fmt.Println("git commit:", GitCommit)
		fmt.Println("build at:", BuildAt)
		fmt.Println("build by:", BuildBy)
	}

	app := cli.NewApp()
	app.Version = Version
	app.Name = Name

	app.Action = func(c *cli.Context) error {
		core := frame.New(
			frame.Name(Name),
			frame.Version(Version),
			frame.Address(":10010"),
		)

		h, err := handler.New(config.DefaultConf)
		if err != nil {
			panic(err)
		}
		proto.RegisterAccountServer(core, h)

		if err := core.Start(); err != nil {
			return err
		}
		defer core.Stop()

		//start http server
		httpServer(h)

		s := utilSignal.WaitShutdown()
		logger.Infof("recv signal: %v", s.String())
		return nil
	}

	err := app.Run(os.Args)
	return err
}

// var gDBConn = newDBConn("root:123456@tcp(127.0.0.1:3306)/account?charset=utf8mb4")
// var ct *handlerHelper.CallTable
// var h *handler.Handler

// func Run(ctx *cli.Context) error {
// 	var err error
// 	if h, err = handler.NewHandler(); err != nil {
// 		return err
// 	}

// 	ct = handlerHelper.ParseRpcMethod(proto.File_proto_usercenter_proto.Services(), h)

// 	ListenAddr := ":8070"
// 	fmt.Println("http listen at ", ListenAddr)
// 	http.HandleFunc("/UserCenter/", httpRoute)

// 	return http.ListenAndServe(ListenAddr, nil)
// }

//todo:
type Codec interface {
	Unmarshal(b []byte, m protobuf.Message) error
	Marshal(m protobuf.Message) ([]byte, error)
}

var calltable = utilHandle.ExtractProtoFile(proto.File_servers_account_proto_account_proto, &handler.Handler{})

func httpServer(handler interface{}) {

	go func() {
		http.HandleFunc("/account/", GRPCMethodToHttp(handler, calltable))

		err := http.ListenAndServe(":8082", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func GRPCMethodToHttp(handler interface{}, ct *utilHandle.CallTable) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		respWithError := func(data json.RawMessage, err error) {
			type HttpRespType struct {
				Data    json.RawMessage `json:"data"`
				Code    int             `json:"code"`
				Message string          `json:"message"`
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

		path := r.URL.Path
		path = strings.TrimPrefix(path, "/")
		if len(path) <= 1 {
			respWithError(nil, fmt.Errorf("method can not be: %s", path))
			return
		}
		method := ct.Get(path)
		if method == nil {
			respWithError(nil, fmt.Errorf("method not found: %s", path))
			return
		}

		raw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respWithError(nil, fmt.Errorf("read body error: %s", err.Error()))
			return
		}

		req := reflect.New(method.RequestType).Interface().(protobuf.Message)

		if err := protojson.Unmarshal(raw, req); err != nil {
			respWithError(nil, fmt.Errorf("unmarshal request error: %s", err.Error()))
			return
		}

		// here call method
		respArgs := method.Call(reflect.ValueOf(handler), reflect.ValueOf(r.Context()), reflect.ValueOf(req))

		if len(respArgs) != 2 {
			//TODO:
			return
		}

		respErr := respArgs[1].Interface().(error)
		var respData json.RawMessage

		if resp, ok := respArgs[0].Interface().(protobuf.Message); ok {
			if resp != nil {
				if data, err := protojson.Marshal(resp); err == nil {
					respData = data
				} else {
					logger.Error(err)
				}
			}
		}
		respWithError(respData, respErr)
	}
}

func httpHandler(w http.ResponseWriter, reqRaw *http.Request) {
	methodName := strings.TrimPrefix(reqRaw.URL.Path, "/")
	method := calltable.Get(methodName)
	if method == nil {
		w.Write([]byte("method not found"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// req := reflect.New(method.Req).Interface().(protoreflect.ProtoMessage)
	// resp := reflect.New(method.Resp).Interface().(protoreflect.ProtoMessage)

	// defer reqRaw.Body.Close()
	// data, err := ioutil.ReadAll(reqRaw.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// fmt.Println(string(data))
	// //parse json to protobuf
	// unmarshal := &protojson.UnmarshalOptions{}
	// if err = unmarshal.Unmarshal(data, req); err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// //call handler function
	// callResult := method.Method.Func.Call([]reflect.Value{reflect.ValueOf(h), reflect.ValueOf(context.Background()), reflect.ValueOf(req), reflect.ValueOf(resp)})

	// if !callResult[0].IsNil() {
	// 	err = callResult[0].Interface().(error)
	// }

	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(err.Error()))
	// } else {
	// 	respRaw, err := protojson.MarshalOptions{
	// 		EmitUnpopulated: true,
	// 	}.Marshal(resp)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Println(string(respRaw))
	// 	w.Write(respRaw)
	// }
}
