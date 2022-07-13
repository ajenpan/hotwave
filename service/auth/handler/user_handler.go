package handler

import (
	"context"
	"strings"

	"google.golang.org/protobuf/proto"

	log "hotwave/logger"
	gwclient "hotwave/service/gateway/client"
	gwproto "hotwave/service/gateway/proto"
	"hotwave/utils/calltable"
)

func (a *Auth) OnUserMessage(s *gwclient.UserSession, msg *gwproto.ToServerMessage) {
	// here, msg.Name is something like "auth.LoginRequest"
	// method key likes "Auth/Login"

	//TODO: batter way to match the method key
	methodName := msg.Name
	methodName = strings.TrimPrefix(methodName, "auth.")
	methodName = strings.TrimSuffix(methodName, "Request")
	methodName = "Auth/" + methodName

	method := a.CT.Get(methodName)

	if method == nil {
		log.Warn("method not found: ", methodName)
		return
	}

	req := method.NewRequest().(proto.Message)
	if err := proto.Unmarshal(msg.Data, req); err != nil {
		log.Error(err)
		return
	}

	if method.Style == calltable.StyleGRpc {
		res := method.Call(context.Background(), req)
		if !res[1].IsNil() {
			log.Error(res[1].Interface().(error))
		}

		if !res[0].IsNil() {
			s.Send(res[0].Interface().(proto.Message))
		}
	}
}
