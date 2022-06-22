package httpsvr

import (
	"context"
	"io/ioutil"
	"net/http"

	"hotwave/transport"
)

type Option func(*Options)

type Options struct {
	Address    string
	Middle     []Middleware
	Marshalers marshalerRegistry
	CT         *transport.CallTable

	// IncomingHeaderMatcher HeaderMatcherFunc
	// OutgoingHeaderMatcher HeaderMatcherFunc
}

func NewOption(opts ...Option) *Options {
	ret := &Options{}
	for _, v := range opts {
		v(ret)
	}
	return ret
}

func NewServer(opts *Options) *Server {
	s := &Server{
		Options: opts,
		mux:     http.NewServeMux(),
	}

	if s.CT != nil {
		s.CT.Range(func(key string, value *transport.Method) bool {
			callwarp := func(w http.ResponseWriter, r *http.Request) {
				s.onCall(w, r, value)
			}
			for i := len(s.Middle) - 1; i >= 0; i-- {
				callwarp = s.Middle[i](callwarp)
			}
			s.mux.HandleFunc(key, callwarp)
			return true
		})
	}

	httpsvr := &http.Server{
		Addr:    opts.Address,
		Handler: s.mux,
	}
	s.httpsvr = httpsvr
	return s
}

type Server struct {
	*Options

	mux     *http.ServeMux
	httpsvr *http.Server
}

func (s *Server) Start() error {
	err := s.httpsvr.ListenAndServe()
	s.Address = s.httpsvr.Addr
	return err
}

func (s *Server) Stop() error {
	return s.httpsvr.Shutdown(context.Background())
}

func (s *Server) onCall(w http.ResponseWriter, r *http.Request, h *transport.Method) {
	ctx := r.Context()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	in, out := MarshalerForRequest(&s.Marshalers, r)

	req := h.NewRequest()
	resp := h.NewResponse()

	err = in.Unmarshal(b, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	callResult := h.Call(nil, ctx, req, resp)

	if !callResult[0].IsNil() {
		err = callResult[0].Interface().(error)
	}

	var respRaw []byte

	if err == nil {
		respRaw, err = out.Marshal(resp)
	} else {
		respRaw, err = out.Marshal(err)
	}

	if err == nil {
		w.Header().Set("Content-Type", out.ContentType(nil))
		w.Write(respRaw)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
