package transport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ExtractParseGRpcMethod(packageName string, ms protoreflect.ServiceDescriptors, h interface{}) *CallTable {
	refh := reflect.TypeOf(h)

	ret := &CallTable{
		list: make(map[string]*Method),
	}

	for i := 0; i < ms.Len(); i++ {
		rpcMethods := ms.Get(i).Methods()

		for j := 0; j < rpcMethods.Len(); j++ {
			rpcMethod := rpcMethods.Get(j)
			rpcMethodName := string(rpcMethod.Name())

			method, has := refh.MethodByName(rpcMethodName)
			if !has {
				continue
			}

			epn := strings.Join([]string{packageName, rpcMethodName}, endpointSplit)

			if method.Type.NumIn() != 3 || method.Type.NumOut() != 2 {
				continue
			}

			reqType := method.Type.In(2).Elem()
			respType := method.Type.Out(0).Elem()

			ret.list[epn] = &Method{
				Imp:          method,
				Style:        StyleMicro,
				RequestType:  reqType,
				ResponseType: respType,
			}
		}
	}
	return ret
}

func ExtractAsyncMethod(packageName string, ms protoreflect.MessageDescriptors, h interface{}) *CallTable {
	const Messagesuffix string = "Request"
	const MethodPrefix string = "On"
	refh := reflect.TypeOf(h)

	ret := &CallTable{
		list: make(map[string]*Method),
	}

	for i := 0; i < ms.Len(); i++ {
		msg := ms.Get(i)
		requestName := string(msg.Name())
		if !strings.HasSuffix(requestName, Messagesuffix) {
			continue
		}

		method, has := refh.MethodByName(MethodPrefix + requestName)
		if !has {
			continue
		}

		epn := strings.Join([]string{packageName, strings.TrimSuffix(requestName, Messagesuffix)}, endpointSplit)

		// func (context.Context, proto.Message) (error)
		if method.Type.NumIn() != 3 {
			continue
		}

		ret.list[epn] = &Method{
			Imp:         method,
			Style:       StyleAsync,
			RequestType: method.Type.In(2).Elem(),
		}
	}
	return ret
}

func ExtractProtoFile(fd protoreflect.FileDescriptor, handler interface{}) *CallTable {
	ret := &CallTable{
		list: make(map[string]*Method),
	}

	href := reflect.TypeOf(handler)

	rpcTable := ExtractParseGRpcMethod(string(fd.Package()), fd.Services(), href)
	asyncTalbe := ExtractAsyncMethod(string(fd.Package()), fd.Messages(), href)

	ret.Merge(rpcTable, false)
	ret.Merge(asyncTalbe, false)

	return ret
}

func ServerGRPCMethodForHttp(handler interface{}, ct *CallTable) func(w http.ResponseWriter, r *http.Request) {
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
		path = strings.Trim(path, "/")
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
			return
		}

		var respErr error
		if !respArgs[1].IsNil() {
			respErr = respArgs[1].Interface().(error)
		}

		var respData json.RawMessage

		if !respArgs[0].IsNil() {
			if resp, ok := respArgs[0].Interface().(protobuf.Message); ok {
				if data, err := protojson.Marshal(resp); err == nil {
					respData = data
				} else {
					respWithError(nil, respErr)
				}
			}
		}

		respWithError(respData, respErr)
	}
}
