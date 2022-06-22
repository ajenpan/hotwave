package httpsvr

import (
	"errors"
	"log"
	"mime"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"hotwave/marshal"
)

// MIMEWildcard is the fallback MIME type used for requests which do not match
// a registered MIME type.
const MIMEWildcard = "*"

var (
	acceptHeader      = http.CanonicalHeaderKey("Accept")
	contentTypeHeader = http.CanonicalHeaderKey("Content-Type")

	defaultMarshaler = &marshal.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
)

type Marshaler = marshal.Marshaler

// MarshalerForRequest returns the inbound/outbound marshalers for this request.
// It checks the registry on the ServeMux for the MIME type set by the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*"/InboundMarshaler/OutboundMarshaler.
func MarshalerForRequest(marshalers *marshalerRegistry, r *http.Request) (inbound Marshaler, outbound Marshaler) {
	for _, acceptVal := range r.Header[acceptHeader] {
		if m, ok := marshalers.mimeMap[acceptVal]; ok {
			outbound = m
			break
		}
	}

	for _, contentTypeVal := range r.Header[contentTypeHeader] {
		contentType, _, err := mime.ParseMediaType(contentTypeVal)
		if err != nil {
			log.Printf("Failed to parse Content-Type %s: %v", contentTypeVal, err)
			continue
		}
		if m, ok := marshalers.mimeMap[contentType]; ok {
			inbound = m
			break
		}
	}

	if inbound == nil {
		inbound = marshalers.mimeMap[MIMEWildcard]
	}
	if outbound == nil {
		outbound = inbound
	}

	return inbound, outbound
}

// marshalerRegistry is a mapping from MIME types to Marshalers.
type marshalerRegistry struct {
	mimeMap map[string]marshal.Marshaler
}

// Add adds a marshaler for a case-sensitive MIME type string ("*" to match any MIME type).
func (m marshalerRegistry) Add(mime string, marshaler Marshaler) error {
	if len(mime) == 0 {
		return errors.New("empty MIME type")
	}
	m.mimeMap[mime] = marshaler
	return nil
}

func WithMarshalerOption(mime string, marshaler Marshaler) Option {
	return func(mux *Options) {
		if err := mux.Marshalers.Add(mime, marshaler); err != nil {
			panic(err)
		}
	}
}
