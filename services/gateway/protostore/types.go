package protostore

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MessageMap map[string]protoreflect.MessageDescriptor
type MethodMap map[string]protoreflect.MethodDescriptor
type FilesMap map[string]protoreflect.FileDescriptor
