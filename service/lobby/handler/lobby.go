package handler

import (
	"context"
	"sync"

	log "hotwave/logger"
	nodetp "hotwave/node/transport"
	"hotwave/node/usergater"
	battleProto "hotwave/service/battle/proto"
	"hotwave/service/lobby/proto"
	protocal "hotwave/service/lobby/proto"
	"hotwave/service/lobby/user"
)

type Lobby struct {
	rwlock sync.RWMutex
	users  []*user.UserInfo

	battleclient battleProto.BattleClient
}

func (l *Lobby) OnUserPropsInfoRequest(socket usergater.UserSocket, in *protocal.UserPropsInfoRequest) {

}

func (l *Lobby) OnUserMatchRequest(socket usergater.UserSocket, in *protocal.UserMatchRequest) {
	l.addUserToQueue(&user.UserInfo{
		UserId: socket.UID(),
		Socket: socket,
	})
}

func (l *Lobby) addUserToQueue(u *user.UserInfo) {
	l.rwlock.Lock()
	defer l.rwlock.Unlock()

	l.users = append(l.users, u)

	const expert = 4

	if len(l.users) >= expert {

		users := l.users[:expert]

		battlePlayerInfo := make([]*battleProto.PlayerInfo, len(users))
		lobbyPlayerInfo := make([]*proto.PlayerInfo, len(users))

		for i, v := range users {
			battlePlayerInfo[i] = &battleProto.PlayerInfo{
				Uid:     v.UserId,
				SeatId:  int32(i) + 1,
				Score:   1000,
				IsRobot: false,
			}

			lobbyPlayerInfo[i] = &proto.PlayerInfo{
				Uid:    v.UserId,
				SeatId: int32(i) + 1,
			}
		}

		//create battle
		resp, err := l.battleclient.CreateBattle(context.Background(), &battleProto.CreateBattleRequest{
			GameName:    "niuniu",
			PlayerInfos: battlePlayerInfo,
		})

		notice := &proto.UserGameStartNotify{
			BattleId: resp.BattleId,
			Players:  lobbyPlayerInfo,
		}

		if err != nil {
			log.Error(err)
			notice.Errcode = -1
			notice.Errmsg = "create battle failed"
		} else {
			notice.Errcode = 0
			notice.Errmsg = "ok"
		}

		for _, user := range users {
			user.Socket.Send(&nodetp.Message{
				Body: notice,
			})
		}
	}
}
