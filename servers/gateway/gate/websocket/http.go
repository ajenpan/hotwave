package websocket

import (
	"io/ioutil"
	"net"
	"net/http"

	"hotwave/servers/gateway/gate/codec"
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

	// httpsvr : &http.Server{Addr: addr, Handler: handler},

	ln, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	go func() {
		err := s.httpsvr.Serve(ln)
		if err != nil {

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
	// c, err := upgrader.Upgrade(w, r, nil)

	defer r.Body.Close()
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	if s.opts.Adapter != nil {

		session := NewSession(rw, r)

		s.opts.Adapter.OnGateMessage(session, &codec.AsyncMessage{
			MsgName: r.URL.Path,
			Body:    raw,
		})
	}

}
