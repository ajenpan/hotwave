package hotwave

//event
//go:generate ./tools/protoc/bin/protoc --go_out=. --go-grpc_out=. event/proto/*.proto

//service
//go:generate ./tools/protoc/bin/protoc --go_out=. --go-grpc_out=. service/gateway/proto/*.proto
//go:generate ./tools/protoc/bin/protoc --go_out=. --go-grpc_out=. service/auth/proto/*.proto
//go:generate ./tools/protoc/bin/protoc --go_out=. --go-grpc_out=. service/battle/proto/*.proto

//game
//go:generate ./tools/protoc/bin/protoc --go_out=. --go-grpc_out=. game/niuniu/*.proto
