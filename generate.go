package hotwave

//event
//go:generate protoc  -I=./tools/protoc/include -I=./event/proto --go_out=. --go-grpc_out=. ./event/proto/*.proto

//service
//go:generate protoc -I=./tools/protoc/include -I=./service/gateway/proto --go_out=. --go-grpc_out=. ./service/gateway/proto/*.proto
//go:generate protoc -I=./tools/protoc/include -I=./service/auth/proto --go_out=. --go-grpc_out=. ./service/auth/proto/*.proto
//go:generate protoc -I=./tools/protoc/include -I=./service/battle/proto --go_out=. --go-grpc_out=. ./service/battle/proto/*.proto

//game
//go:generate protoc -I=./tools/protoc/include -I=./game/niuniu --go_out=. --go-grpc_out=. ./game/niuniu/*.proto
