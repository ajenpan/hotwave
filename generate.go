package hotwave

//event
//go:generate protoc --go_out=. --go-grpc_out=. ./event/proto/*.proto

//service
//go:generate protoc --go_out=. --go-grpc_out=. ./service/gateway/proto/*.proto
//go:generate protoc --go_out=. --go-grpc_out=. ./service/auth/proto/*.proto
