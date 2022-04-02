package proto

import (
	"google.golang.org/protobuf/proto"

	"hotwave/codec"
)

type Marshaler struct{}

func (Marshaler) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, codec.ErrInvalidMessage
	}
	return proto.Marshal(pb)
}

func (Marshaler) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return codec.ErrInvalidMessage
	}
	return proto.Unmarshal(data, pb)
}

func (Marshaler) String() string {
	return "proto"
}
