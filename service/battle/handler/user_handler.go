package handler

import (
	log "hotwave/logger"
	gwclient "hotwave/service/gateway/client"
	gwproto "hotwave/service/gateway/proto"
)

func (a *Handler) OnUserMessage(s *gwclient.UserSession, msg *gwproto.ToServerMessage) {
	log.Info("UserMessage", s.UID, msg.Name)

	itable := a.geBattleByUid(s.UID)
	if itable == nil {
		return
	}

	itable.OnPlayerMessage(s.UID, msg.Name, msg.Data)
}
