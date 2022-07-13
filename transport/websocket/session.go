package websocket

import (
	"fmt"
	"net/http"
	"sync"

	protobuf "google.golang.org/protobuf/proto"

	"hotwave/transport"
)

func NewSession(rw http.ResponseWriter, r *http.Request) *httpSession {
	ret := &httpSession{
		rw: rw,
		r:  r,
	}
	return ret
}

type httpSession struct {
	transport.SessionMeta

	id          string
	rw          http.ResponseWriter
	r           *http.Request
	respMsgLock sync.Mutex
	respMsg     protobuf.Message
}

func (s *httpSession) ID() string {
	return s.id
}

func (s *httpSession) Send(resp protobuf.Message) error {
	s.respMsgLock.Lock()
	defer s.respMsgLock.Unlock()
	if s.respMsg == nil {
		// s.respMsg = resp
		return nil
	}
	return fmt.Errorf("session already has respMsg")
}

func (s *httpSession) Close() {
	//do nothing here
}
