package transport

import (
	"reflect"
	"sync"
)

type MethodStyle int

// Todo:
const (
	StyleMicro = iota // func (context.Context, proto.Message, proto.Message) ( error)
	StyleGRpc  = iota // func (context.Context, proto.Message) (proto.Message, error)
	StyleAsync = iota // func (context.Context, proto.Message) (error)
)

const endpointSplit = "/"

type Method struct {
	Imp   reflect.Method
	Style MethodStyle

	RequestType  reflect.Type
	ResponseType reflect.Type
}

func (m *Method) Call(args ...interface{}) []reflect.Value {
	values := make([]reflect.Value, 0, len(args))
	for _, v := range args {
		values = append(values, reflect.ValueOf(v))
	}
	return m.Imp.Func.Call(values)
}

func (m *Method) NewRequest() interface{} {
	return reflect.New(m.RequestType).Interface()
}

func (m *Method) NewResponse() interface{} {
	return reflect.New(m.ResponseType).Interface()
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
