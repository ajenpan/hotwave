// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: service/gateway/proto/gateway.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ToUserMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data       []byte      `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	ToUid      int64       `protobuf:"varint,3,opt,name=to_uid,json=toUid,proto3" json:"to_uid,omitempty"`
	ToSocketid string      `protobuf:"bytes,4,opt,name=to_socketid,json=toSocketid,proto3" json:"to_socketid,omitempty"`
	Mime       MIMEType    `protobuf:"varint,5,opt,name=mime,proto3,enum=gateway.MIMEType" json:"mime,omitempty"`
	Type       MessageType `protobuf:"varint,6,opt,name=type,proto3,enum=gateway.MessageType" json:"type,omitempty"`
}

func (x *ToUserMessage) Reset() {
	*x = ToUserMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_gateway_proto_gateway_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToUserMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToUserMessage) ProtoMessage() {}

func (x *ToUserMessage) ProtoReflect() protoreflect.Message {
	mi := &file_service_gateway_proto_gateway_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToUserMessage.ProtoReflect.Descriptor instead.
func (*ToUserMessage) Descriptor() ([]byte, []int) {
	return file_service_gateway_proto_gateway_proto_rawDescGZIP(), []int{0}
}

func (x *ToUserMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ToUserMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ToUserMessage) GetToUid() int64 {
	if x != nil {
		return x.ToUid
	}
	return 0
}

func (x *ToUserMessage) GetToSocketid() string {
	if x != nil {
		return x.ToSocketid
	}
	return ""
}

func (x *ToUserMessage) GetMime() MIMEType {
	if x != nil {
		return x.Mime
	}
	return MIMEType_Protobuf
}

func (x *ToUserMessage) GetType() MessageType {
	if x != nil {
		return x.Type
	}
	return MessageType_Aync
}

type ToServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data         []byte      `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	FromUid      int64       `protobuf:"varint,3,opt,name=from_uid,json=fromUid,proto3" json:"from_uid,omitempty"`
	FromSocketid string      `protobuf:"bytes,4,opt,name=from_socketid,json=fromSocketid,proto3" json:"from_socketid,omitempty"`
	Mime         MIMEType    `protobuf:"varint,5,opt,name=mime,proto3,enum=gateway.MIMEType" json:"mime,omitempty"`
	Type         MessageType `protobuf:"varint,6,opt,name=type,proto3,enum=gateway.MessageType" json:"type,omitempty"`
}

func (x *ToServerMessage) Reset() {
	*x = ToServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_gateway_proto_gateway_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToServerMessage) ProtoMessage() {}

func (x *ToServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_service_gateway_proto_gateway_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToServerMessage.ProtoReflect.Descriptor instead.
func (*ToServerMessage) Descriptor() ([]byte, []int) {
	return file_service_gateway_proto_gateway_proto_rawDescGZIP(), []int{1}
}

func (x *ToServerMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ToServerMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ToServerMessage) GetFromUid() int64 {
	if x != nil {
		return x.FromUid
	}
	return 0
}

func (x *ToServerMessage) GetFromSocketid() string {
	if x != nil {
		return x.FromSocketid
	}
	return ""
}

func (x *ToServerMessage) GetMime() MIMEType {
	if x != nil {
		return x.Mime
	}
	return MIMEType_Protobuf
}

func (x *ToServerMessage) GetType() MessageType {
	if x != nil {
		return x.Type
	}
	return MessageType_Aync
}

type AddGateAllowListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Names []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
}

func (x *AddGateAllowListRequest) Reset() {
	*x = AddGateAllowListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_gateway_proto_gateway_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddGateAllowListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddGateAllowListRequest) ProtoMessage() {}

func (x *AddGateAllowListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_gateway_proto_gateway_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddGateAllowListRequest.ProtoReflect.Descriptor instead.
func (*AddGateAllowListRequest) Descriptor() ([]byte, []int) {
	return file_service_gateway_proto_gateway_proto_rawDescGZIP(), []int{2}
}

func (x *AddGateAllowListRequest) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

type AddGateAllowListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddGateAllowListResponse) Reset() {
	*x = AddGateAllowListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_gateway_proto_gateway_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddGateAllowListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddGateAllowListResponse) ProtoMessage() {}

func (x *AddGateAllowListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_gateway_proto_gateway_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddGateAllowListResponse.ProtoReflect.Descriptor instead.
func (*AddGateAllowListResponse) Descriptor() ([]byte, []int) {
	return file_service_gateway_proto_gateway_proto_rawDescGZIP(), []int{3}
}

var File_service_gateway_proto_gateway_proto protoreflect.FileDescriptor

var file_service_gateway_proto_gateway_proto_rawDesc = []byte{
	0x0a, 0x23, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61,
	0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x1a, 0x22,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xc0, 0x01, 0x0a, 0x0d, 0x54, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x15, 0x0a, 0x06,
	0x74, 0x6f, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f,
	0x55, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x5f, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x6f, 0x53, 0x6f, 0x63, 0x6b,
	0x65, 0x74, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x11, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x4d, 0x49, 0x4d,
	0x45, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xca, 0x01, 0x0a, 0x0f, 0x54, 0x6f, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x66, 0x72, 0x6f, 0x6d, 0x55, 0x69, 0x64, 0x12, 0x23, 0x0a, 0x0d,
	0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x72, 0x6f, 0x6d, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x69,
	0x64, 0x12, 0x25, 0x0a, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x11, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x4d, 0x49, 0x4d, 0x45, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x22, 0x2f, 0x0a, 0x17, 0x41, 0x64, 0x64, 0x47, 0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c,
	0x6f, 0x77, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x22, 0x1a, 0x0a, 0x18, 0x41, 0x64, 0x64, 0x47, 0x61, 0x74, 0x65, 0x41, 0x6c,
	0x6c, 0x6f, 0x77, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32,
	0xa5, 0x01, 0x0a, 0x07, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x12, 0x3f, 0x0a, 0x05, 0x50,
	0x72, 0x6f, 0x78, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x54,
	0x6f, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x18, 0x2e, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x54, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x59, 0x0a, 0x10,
	0x41, 0x64, 0x64, 0x47, 0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x20, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x41, 0x64, 0x64, 0x47, 0x61,
	0x74, 0x65, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x21, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x41, 0x64, 0x64,
	0x47, 0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x17, 0x5a, 0x15, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_gateway_proto_gateway_proto_rawDescOnce sync.Once
	file_service_gateway_proto_gateway_proto_rawDescData = file_service_gateway_proto_gateway_proto_rawDesc
)

func file_service_gateway_proto_gateway_proto_rawDescGZIP() []byte {
	file_service_gateway_proto_gateway_proto_rawDescOnce.Do(func() {
		file_service_gateway_proto_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_gateway_proto_gateway_proto_rawDescData)
	})
	return file_service_gateway_proto_gateway_proto_rawDescData
}

var file_service_gateway_proto_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_service_gateway_proto_gateway_proto_goTypes = []interface{}{
	(*ToUserMessage)(nil),            // 0: gateway.ToUserMessage
	(*ToServerMessage)(nil),          // 1: gateway.ToServerMessage
	(*AddGateAllowListRequest)(nil),  // 2: gateway.AddGateAllowListRequest
	(*AddGateAllowListResponse)(nil), // 3: gateway.AddGateAllowListResponse
	(MIMEType)(0),                    // 4: gateway.MIMEType
	(MessageType)(0),                 // 5: gateway.MessageType
}
var file_service_gateway_proto_gateway_proto_depIdxs = []int32{
	4, // 0: gateway.ToUserMessage.mime:type_name -> gateway.MIMEType
	5, // 1: gateway.ToUserMessage.type:type_name -> gateway.MessageType
	4, // 2: gateway.ToServerMessage.mime:type_name -> gateway.MIMEType
	5, // 3: gateway.ToServerMessage.type:type_name -> gateway.MessageType
	0, // 4: gateway.Gateway.Proxy:input_type -> gateway.ToUserMessage
	2, // 5: gateway.Gateway.AddGateAllowList:input_type -> gateway.AddGateAllowListRequest
	1, // 6: gateway.Gateway.Proxy:output_type -> gateway.ToServerMessage
	3, // 7: gateway.Gateway.AddGateAllowList:output_type -> gateway.AddGateAllowListResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_service_gateway_proto_gateway_proto_init() }
func file_service_gateway_proto_gateway_proto_init() {
	if File_service_gateway_proto_gateway_proto != nil {
		return
	}
	file_service_gateway_proto_client_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_service_gateway_proto_gateway_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToUserMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_gateway_proto_gateway_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToServerMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_gateway_proto_gateway_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddGateAllowListRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_gateway_proto_gateway_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddGateAllowListResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_gateway_proto_gateway_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_gateway_proto_gateway_proto_goTypes,
		DependencyIndexes: file_service_gateway_proto_gateway_proto_depIdxs,
		MessageInfos:      file_service_gateway_proto_gateway_proto_msgTypes,
	}.Build()
	File_service_gateway_proto_gateway_proto = out.File
	file_service_gateway_proto_gateway_proto_rawDesc = nil
	file_service_gateway_proto_gateway_proto_goTypes = nil
	file_service_gateway_proto_gateway_proto_depIdxs = nil
}
