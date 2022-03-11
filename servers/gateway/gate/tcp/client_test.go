package tcp

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestClientReconnect(t *testing.T) {
	address := fmt.Sprintf("localhost:%d", rand.Int31n(1000)+20000)
	s := NewServer(&ServerOptions{
		Address: address,
	})
	err := s.Start()
	if err != nil {
		t.Fatal(err)
	}

	c := NewClient(&ClientOptions{
		Address: address,
	})

	if err := c.Connect(); err != nil {
		t.Fatal(err)
	}

}
