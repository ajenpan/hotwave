package gateway

import (
	"fmt"
	"testing"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func TestGateway(t *testing.T) {
	mx := gwruntime.NewServeMux()
	fmt.Println(mx.GetForwardResponseOptions())
	// proto.RegisterGateAdapterServer(mx, &noophandler{})
	// mx.Handle("/", &noophandler{})
	//runtime.AnnotateIncomingContext()
	// ctx, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/GetResponseBody", runtime.WithHTTPPathPattern("/responsebody/{data}"))

}
