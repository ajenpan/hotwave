package handle

import (
	"reflect"
	"strings"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type MethodType int //MethodType
const (
	CallTableType_RPC   = iota
	CallTableType_Async = iota
)

const endpointSplit = "/"

type Method struct {
	Method reflect.Method
	Typ    MethodType

	RequestType  reflect.Type
	ResponseType reflect.Type
}

func (mm *Method) Call(args ...reflect.Value) []reflect.Value {
	return mm.Method.Func.Call(args)
}

type CallTable struct {
	sync.RWMutex
	list map[string]*Method
}

func (m *CallTable) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.list)
}

func (m *CallTable) Has(name string) bool {
	m.RLock()
	defer m.RUnlock()
	_, has := m.list[name]
	return has
}

func (m *CallTable) Get(name string) *Method {
	m.RLock()
	defer m.RUnlock()

	ret, has := m.list[name]
	if has {
		return ret
	}
	return nil
}

func (m *CallTable) Range(f func(key string, value *Method) bool) {
	m.Lock()
	defer m.Unlock()
	for k, v := range m.list {
		if !f(k, v) {
			return
		}
	}
}

func (m *CallTable) Merge(other *CallTable, overWrite bool) int {
	ret := 0
	other.RWMutex.RLock()
	defer other.RWMutex.RUnlock()

	m.Lock()
	defer m.Unlock()

	for k, v := range other.list {
		_, has := m.list[k]
		if has && !overWrite {
			continue
		}
		m.list[k] = v
		ret++
	}
	return ret
}

func ExtractParseGRpcMethod(packageName string, ms protoreflect.ServiceDescriptors, refh reflect.Type) *CallTable {
	ret := &CallTable{
		list: make(map[string]*Method),
	}

	for i := 0; i < ms.Len(); i++ {
		rpcMethods := ms.Get(i).Methods()

		for j := 0; j < rpcMethods.Len(); j++ {
			rpcMethod := rpcMethods.Get(j)
			rpcMethodName := string(rpcMethod.Name())

			method, has := refh.MethodByName(rpcMethodName)
			if !has {
				continue
			}

			epn := strings.Join([]string{packageName, rpcMethodName}, endpointSplit)

			if method.Type.NumIn() != 3 || method.Type.NumOut() != 2 {
				continue
			}

			reqType := method.Type.In(2).Elem()
			respType := method.Type.Out(0).Elem()

			ret.list[epn] = &Method{
				Method:       method,
				Typ:          CallTableType_RPC,
				RequestType:  reqType,
				ResponseType: respType,
			}
		}
	}
	return ret
}

func ExtractAsyncMethod(packageName string, ms protoreflect.MessageDescriptors, refh reflect.Type) *CallTable {
	const Messagesuffix string = "Request"
	const MethodPrefix string = "On"

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

		// object - session - requestmessage
		if method.Type.NumIn() != 3 {
			continue
		}

		ret.list[epn] = &Method{
			Method:      method,
			Typ:         CallTableType_Async,
			RequestType: method.Type.In(2).Elem(),
		}
	}
	return ret
}

func ExtractProtoFile(fd protoreflect.FileDescriptor, handler interface{}) *CallTable {
	ret := &CallTable{
		list: make(map[string]*Method),
	}

	href := reflect.TypeOf(handler)

	rpcTable := ExtractParseGRpcMethod(string(fd.Package()), fd.Services(), href)
	asyncTalbe := ExtractAsyncMethod(string(fd.Package()), fd.Messages(), href)

	ret.Merge(rpcTable, false)
	ret.Merge(asyncTalbe, false)

	return ret
}
