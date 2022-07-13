package common

import (
	"context"

	"google.golang.org/protobuf/proto"

	"hotwave/transport"
	"hotwave/utils/calltable"
)

func CallHelper(method *calltable.Method, session transport.Session, iraw interface{}) error {
	//decode request
	var req proto.Message
	switch iraw := iraw.(type) {
	case proto.Message:
		req = iraw
	case []byte:
		req = method.NewRequest().(proto.Message)
		if err := proto.Unmarshal(iraw, req); err != nil {
			return err
		}
	}

	switch method.Style {
	case calltable.StyleRequest:
		res := method.Call(session, req)
		if !res[1].IsNil() {
			return res[1].Interface().(error)
		}
		if !res[0].IsNil() {
			return session.Send(res[0].Interface().(proto.Message))
		}
	case calltable.StyleAsync:
		res := method.Call(session, req)
		if len(res) == 1 && !res[0].IsNil() {
			return res[0].Interface().(error)
		}
	case calltable.StyleGRpc:
		ctx := CtxWithSocket(context.Background(), session)
		res := method.Call(ctx, req)
		if !res[1].IsNil() {
			return res[1].Interface().(error)
		}
		if !res[0].IsNil() {
			return session.Send(res[0].Interface().(proto.Message))
		}
	}
	return nil
}
