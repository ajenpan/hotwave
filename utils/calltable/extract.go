package calltable

import (
	"reflect"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func ExtractParseGRpcMethod(ms protoreflect.ServiceDescriptors, h interface{}) *CallTable {
	refh := reflect.TypeOf(h)

	ret := &CallTable{
		list: make(map[string]*Method),
	}

	for i := 0; i < ms.Len(); i++ {
		service := ms.Get(i)
		methods := service.Methods()
		svrName := string(service.Name())

		for j := 0; j < methods.Len(); j++ {
			rpcMethod := methods.Get(j)
			rpcMethodName := string(rpcMethod.Name())

			method, has := refh.MethodByName(rpcMethodName)
			if !has {
				continue
			}

			epn := strings.Join([]string{svrName, rpcMethodName}, endpointSplit)

			if method.Type.NumIn() != 3 || method.Type.NumOut() != 2 {
				continue
			}

			reqType := method.Type.In(2).Elem()
			respType := method.Type.Out(0).Elem()

			m := &Method{
				Imp:          method,
				Style:        StyleMicro,
				RequestType:  reqType,
				ResponseType: respType,
			}
			m.init()
			ret.list[epn] = m
		}
	}
	return ret
}

func ExtractAsyncMethod(packageName string, ms protoreflect.MessageDescriptors, h interface{}) *CallTable {
	const Messagesuffix string = "Request"
	const MethodPrefix string = "On"
	refh := reflect.TypeOf(h)

	ret := &CallTable{
		list: make(map[string]*Method),
	}

	for i := 0; i < ms.Len(); i++ {
		msg := ms.Get(i)
		requestName := string(msg.Name())
		if !strings.HasSuffix(requestName, Messagesuffix) {
			continue
		}

		method, has := refh.MethodByName(MethodPrefix + requestName)
		if !has {
			continue
		}

		epn := strings.Join([]string{packageName, strings.TrimSuffix(requestName, Messagesuffix)}, endpointSplit)

		if method.Type.NumIn() != 3 {
			continue
		}

		m := &Method{
			H:           h,
			Imp:         method,
			Style:       StyleAsync,
			RequestType: method.Type.In(2).Elem(),
		}

		m.init()
		ret.list[epn] = m
	}
	return ret
}

func ExtractProtoFile(fd protoreflect.FileDescriptor, handler interface{}) *CallTable {
	ret := &CallTable{
		list: make(map[string]*Method),
	}

	rpcTable := ExtractParseGRpcMethod(fd.Services(), handler)
	asyncTalbe := ExtractAsyncMethod(string(fd.Package()), fd.Messages(), handler)

	ret.Merge(rpcTable, false)
	ret.Merge(asyncTalbe, false)

	return ret
}
