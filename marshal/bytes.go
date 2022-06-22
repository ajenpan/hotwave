package marshal

import (
	"errors"
)

type BytesMarshaler struct{}

func (n BytesMarshaler) Marshal(v interface{}) ([]byte, error) {
	switch ve := v.(type) {
	case *[]byte:
		return *ve, nil
	case []byte:
		return ve, nil
	}
	return nil, errors.New("invalid message")
}

func (n BytesMarshaler) Unmarshal(d []byte, v interface{}) error {
	switch ve := v.(type) {
	case *[]byte:
		*ve = d
	}
	return errors.New("invalid message")
}

func (n BytesMarshaler) String() string {
	return "bytes"
}
