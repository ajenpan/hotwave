package protostore

import (
	"fmt"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func NewMomoryStore() *MomoryStore {
	return &MomoryStore{
		methodStore:  make(map[string]MethodMap),
		messageStore: make(map[string]MessageMap),
	}
}

type MomoryStore struct {
	methodLock  sync.RWMutex
	methodStore map[string]MethodMap

	messageLock  sync.RWMutex
	messageStore map[string]MessageMap
}

func (store *MomoryStore) StoreProtoFiles(server string, files *descriptorpb.FileDescriptorSet) error {
	f, err := protodesc.NewFiles(files)
	if err != nil {
		return err
	}

	methodMap := make(MethodMap)
	msgMap := make(MessageMap)

	f.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		ss := fd.Services()
		for i := 0; i < ss.Len(); i++ {
			s := ss.Get(i)
			methods := s.Methods()
			for j := 0; j < methods.Len(); j++ {
				m := methods.Get(j)
				fmt.Println(m.Name())
				eps := strings.Join([]string{string(fd.Package()), string(s.FullName()), string(m.Name())}, "/")
				methodMap[eps] = m
			}
		}

		msgs := fd.Messages()
		for i := 0; i < msgs.Len(); i++ {
			m := msgs.Get(i)
			msgMap[string(m.FullName())] = m
		}
		return true
	})

	store.methodLock.Lock()
	store.methodStore[server] = methodMap
	store.methodLock.Unlock()

	store.messageLock.Lock()
	store.messageStore[server] = msgMap
	store.messageLock.Unlock()
	return nil
}

func (store *MomoryStore) NewTypeByMethod(server string, method string) (proto.Message, proto.Message, error) {
	store.methodLock.RLock()
	methodMap, ok := store.methodStore[server]
	store.methodLock.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("server %s not found", server)
	}

	methodDescriptor, ok := methodMap[method]
	if !ok {
		return nil, nil, fmt.Errorf("method %s not found", method)
	}

	req := dynamicpb.NewMessage(methodDescriptor.Input())
	resp := dynamicpb.NewMessage(methodDescriptor.Output())
	return req, resp, nil
}
