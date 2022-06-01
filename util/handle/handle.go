package handle

import (
	"encoding/json"
	"fmt"
	"hotwave/logger"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	protobuf "google.golang.org/protobuf/proto"
)

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
			//TODO:
			logger.Warn("method %s return %d args", path, len(respArgs))
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
