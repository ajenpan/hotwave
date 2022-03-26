package websocket

import (
	"net/http"
	"net/textproto"
	"strings"
	"sync"

	"github.com/gorilla/websocket"

	"hotwave/servers/gateway/gate"
	"hotwave/servers/gateway/gate/codec"
)

var upgrader = websocket.Upgrader{} // use default options

func NewSession(rw http.ResponseWriter, r *http.Request) *httpSession {
	id := gate.NewSessionID("http")
	ret := &httpSession{
		rw:   rw,
		r:    r,
		id:   id,
		meta: header2meta(r.Header),
	}
	return ret
}

func header2meta(header http.Header) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, value := range header {
		key = textproto.CanonicalMIMEHeaderKey(key)
		ret[key] = strings.Join(value, " ")
	}
	return ret
}

type httpSession struct {
	sync.RWMutex

	id       string
	rw       http.ResponseWriter
	r        *http.Request
	meta     map[string]interface{}
	metaLock sync.RWMutex
}

func (s *httpSession) ID() string {
	return s.id
}

func (s *httpSession) UID() uint64 {
	return 0
}

func (s *httpSession) SetUID(uid uint64) {
	s.SetMeta("uid", uid)
}

func (s *httpSession) SetMeta(k string, v interface{}) {
	s.metaLock.Lock()
	defer s.metaLock.Unlock()
	s.meta[k] = v
}

func (s *httpSession) GetMeta(k string) (interface{}, bool) {
	s.metaLock.RLock()
	defer s.metaLock.RUnlock()
	v, ok := s.meta[k]
	return v, ok
}

func (s *httpSession) Send(*codec.AsyncMessage) error {
	return nil
}

func (s *httpSession) Close() {

}
