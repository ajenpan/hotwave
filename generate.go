package hotwave

//go:generate ./tools/protoc/bin/protoc -I ./tools/protoc/include/ -I ./services/auth/proto/    --go_out=.  ./services/auth/proto/*.proto
//go:generate ./tools/protoc/bin/protoc -I ./tools/protoc/include/ -I ./services/gateway/proto/ --go_out=.  ./services/gateway/proto/*.proto
