package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	address := fmt.Sprintf("localhost:%d", rand.Int31n(1000)+20000)
	s := NewServer(&ServerOptions{
		Address: address,
		// OnMessage: func(s *Socket, p *Packet) {
		// 	fmt.Println("recv packet")
		// },
		// OnConnStat: func(s *Socket, stat SocketStat) {
		// 	fmt.Println("recv packet", stat)
		// },
		// The interval on which to register
		HeatbeatInterval: time.Second * 10,
	})

	if err := s.Start(); err != nil {
		t.Fail()
		return
	}

	defer s.Stop()

	_, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		t.Fail()
		return
	}

	time.Sleep(time.Second * 2)
	if s.SocketCount() == 0 {
		t.Fail()
	}

	time.Sleep(time.Second * 20)
	if s.SocketCount() != 0 {
		t.Fail()
	}
}
