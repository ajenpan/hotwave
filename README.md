## hotwave

### generate protobuf files

```sh
#event
 protoc --go_out=. event/proto/*.proto

#service
 protoc --go_out=. service/route/proto/*.proto
 protoc --go_out=. service/auth/proto/*.proto
 protoc --go_out=. service/battle/proto/*.proto
 protoc --go_out=. service/lobby/proto/*.proto
#games
 protoc --go_out=. service/games/niuniu/*.proto

```
