package marshal

import (
	"google.golang.org/protobuf/proto"
)

type ProtoMarshaler struct {
	proto.MarshalOptions
	proto.UnmarshalOptions
}

func (m *ProtoMarshaler) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrInvalidProtobuf
	}
	return m.MarshalOptions.Marshal(pb)
}

func (m *ProtoMarshaler) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrInvalidProtobuf
	}
	return m.UnmarshalOptions.Unmarshal(data, pb)
}

func (*ProtoMarshaler) String() string {
	return "proto"
}

func (*ProtoMarshaler) ContentType(_ interface{}) string {
	return "application/protobuf"
}
