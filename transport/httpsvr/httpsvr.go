package httpsvr

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	protobuf "google.golang.org/protobuf/proto"

	"hotwave/utils/calltable"
	"hotwave/utils/marshal"
)

type HttpSvr struct {
	CT      *calltable.CallTable[string]
	Marshal marshal.Marshaler
	// Log     logger.Logger
	svr  *http.Server
	addr string
}

func (s *HttpSvr) Start() error {
	mux := http.NewServeMux()
	s.serverCallTable(mux, s.CT)
	s.svr = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}
	return s.svr.ListenAndServe()
}

func (s *HttpSvr) Stop() error {
	return s.svr.Close()
}

func (s *HttpSvr) serverCallTable(mux *http.ServeMux, ct *calltable.CallTable[string]) {
	ct.Range(func(key string, method *calltable.Method) bool {
		pattern := "/" + key
		fmt.Println("register http method:", pattern)
		cb := handleCallback(method, s.Marshal)
		mux.HandleFunc(pattern, cb)
		return true
	})
}

func handleCallback(method *calltable.Method, marshal marshal.Marshaler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		req := reflect.New(method.RequestType).Interface().(protobuf.Message)

		if err := marshal.Unmarshal(raw, req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// here call method
		respArgs := method.Call(r.Context(), req)

		if len(respArgs) != 2 {
			return
		}

		var respErr error

		if !respArgs[1].IsNil() {
			respErr = respArgs[1].Interface().(error)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(respErr.Error()))
			return
		}

		if !respArgs[0].IsNil() {
			respData, err := marshal.Marshal(respArgs[0].Interface())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(respErr.Error()))
				return
			}
			w.Write(respData)
		}
		w.WriteHeader(http.StatusOK)
	}
}
