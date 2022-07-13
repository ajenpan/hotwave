package event

import (
	"io"
	"sync"

	"github.com/google/uuid"

	"hotwave/event/proto"
)

type GrpcEventPublisher struct {
	proto.UnimplementedEventServer

	senderlock sync.RWMutex
	senders    map[string]Publisher
}

type grpcEventSender struct {
	topics sync.Map
	que    chan *Event
}

func (sender *grpcEventSender) Publish(ev *Event) {
	sender.que <- ev
}

func (s *GrpcEventPublisher) AddSubPublisher(id string, pub Publisher) string {
	if id == "" {
		id = uuid.NewString()
	}
	s.senderlock.Lock()
	defer s.senderlock.Unlock()
	s.senders[id] = pub
	return id
}

func (s *GrpcEventPublisher) RemoveSubPublisher(id string) {
	s.senderlock.Lock()
	defer s.senderlock.Unlock()

	delete(s.senders, id)
}

func (s *GrpcEventPublisher) Subscribe(req *proto.SubscribeRequest, svr proto.Event_SubscribeServer) error {
	sender := &grpcEventSender{
		que: make(chan *Event, 10),
	}

	for _, topic := range req.Topics {
		sender.topics.Store(topic, true)
	}

	id := uuid.NewString()
	s.AddSubPublisher(id, sender)
	defer s.RemoveSubPublisher(id)

	for {
		select {
		case msg, ok := <-sender.que:
			if !ok {
				return io.EOF
			}
			if _, ok := sender.topics.Load(msg.Topic); ok {
				if err := svr.Send(msg); err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func (s *GrpcEventPublisher) Publish(msg *Event) {
	s.senderlock.RLock()
	defer s.senderlock.RUnlock()
	for _, sender := range s.senders {
		sender.Publish(msg)
	}
}
