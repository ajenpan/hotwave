package websocket

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"

	codecjson "hotwave/codec/json"
	"hotwave/logger"
	"hotwave/servers/gateway/proto"
)

func NewServer(opts Options) *Server {
	ret := &Server{
		opts: opts,
		die:  make(chan bool),
	}
	httpsvr := &http.Server{Addr: opts.Address, Handler: ret}
	ret.httpsvr = httpsvr
	return ret
}

type Server struct {
	opts Options
	die  chan bool

	httpsvr *http.Server
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	go func() {
		err := s.httpsvr.Serve(ln)
		if err != nil {
			logger.Error(err)
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	select {
	case <-s.die:
	default:
		close(s.die)
	}
	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	if s.opts.Adapter == nil {
		rw.WriteHeader(http.StatusNotImplemented)
		rw.Write([]byte("not implemented"))
		return
	}

	req := &proto.ClientMessageWraper{
		Method: r.URL.Path,
		Body:   raw,
	}

	resp, err := s.opts.Adapter.OnGateMethod(r.Context(), req)

	type Response struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}

	respWrap := &Response{
		Message: "ok",
	}

	if err != nil {
		respWrap.Code = -1
		respWrap.Message = err.Error()
		return
	}

	if resp != nil {
		raw, err := codecjson.Marshaler{}.Marshal(resp)
		if err != nil {
			respWrap.Code = -1
			respWrap.Message = err.Error()
		} else {
			respWrap.Data = raw
		}
	}

	wrapraw, err := json.Marshal(respWrap)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(wrapraw)
}
