// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: frame/proto/frame.proto

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

type EventMessageWraper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic    string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Body     []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	CreateBy string `protobuf:"bytes,3,opt,name=create_by,json=createBy,proto3" json:"create_by,omitempty"`
	CreateAt string `protobuf:"bytes,4,opt,name=create_at,json=createAt,proto3" json:"create_at,omitempty"`
}

func (x *EventMessageWraper) Reset() {
	*x = EventMessageWraper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_frame_proto_frame_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventMessageWraper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventMessageWraper) ProtoMessage() {}

func (x *EventMessageWraper) ProtoReflect() protoreflect.Message {
	mi := &file_frame_proto_frame_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventMessageWraper.ProtoReflect.Descriptor instead.
func (*EventMessageWraper) Descriptor() ([]byte, []int) {
	return file_frame_proto_frame_proto_rawDescGZIP(), []int{0}
}

func (x *EventMessageWraper) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *EventMessageWraper) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *EventMessageWraper) GetCreateBy() string {
	if x != nil {
		return x.CreateBy
	}
	return ""
}

func (x *EventMessageWraper) GetCreateAt() string {
	if x != nil {
		return x.CreateAt
	}
	return ""
}

type SteamClosed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *SteamClosed) Reset() {
	*x = SteamClosed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_frame_proto_frame_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SteamClosed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SteamClosed) ProtoMessage() {}

func (x *SteamClosed) ProtoReflect() protoreflect.Message {
	mi := &file_frame_proto_frame_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SteamClosed.ProtoReflect.Descriptor instead.
func (*SteamClosed) Descriptor() ([]byte, []int) {
	return file_frame_proto_frame_proto_rawDescGZIP(), []int{1}
}

func (x *SteamClosed) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_frame_proto_frame_proto protoreflect.FileDescriptor

var file_frame_proto_frame_proto_rawDesc = []byte{
	0x0a, 0x17, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x72,
	0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x66, 0x72, 0x61, 0x6d, 0x65,
	0x22, 0x78, 0x0a, 0x12, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x57, 0x72, 0x61, 0x70, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x12, 0x1b, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x79, 0x12, 0x1b, 0x0a,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74, 0x22, 0x1f, 0x0a, 0x0b, 0x53, 0x74,
	0x65, 0x61, 0x6d, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x32, 0x4d, 0x0a, 0x08, 0x4e,
	0x6f, 0x64, 0x65, 0x42, 0x61, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x19, 0x2e, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x72, 0x61, 0x70,
	0x65, 0x72, 0x1a, 0x12, 0x2e, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2e, 0x53, 0x74, 0x65, 0x61, 0x6d,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x22, 0x00, 0x28, 0x01, 0x42, 0x15, 0x5a, 0x13, 0x2e, 0x2f,
	0x66, 0x72, 0x61, 0x6d, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_frame_proto_frame_proto_rawDescOnce sync.Once
	file_frame_proto_frame_proto_rawDescData = file_frame_proto_frame_proto_rawDesc
)

func file_frame_proto_frame_proto_rawDescGZIP() []byte {
	file_frame_proto_frame_proto_rawDescOnce.Do(func() {
		file_frame_proto_frame_proto_rawDescData = protoimpl.X.CompressGZIP(file_frame_proto_frame_proto_rawDescData)
	})
	return file_frame_proto_frame_proto_rawDescData
}

var file_frame_proto_frame_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_frame_proto_frame_proto_goTypes = []interface{}{
	(*EventMessageWraper)(nil), // 0: frame.EventMessageWraper
	(*SteamClosed)(nil),        // 1: frame.SteamClosed
}
var file_frame_proto_frame_proto_depIdxs = []int32{
	0, // 0: frame.NodeBase.EventMessage:input_type -> frame.EventMessageWraper
	1, // 1: frame.NodeBase.EventMessage:output_type -> frame.SteamClosed
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_frame_proto_frame_proto_init() }
func file_frame_proto_frame_proto_init() {
	if File_frame_proto_frame_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_frame_proto_frame_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventMessageWraper); i {
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
		file_frame_proto_frame_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SteamClosed); i {
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
			RawDescriptor: file_frame_proto_frame_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_frame_proto_frame_proto_goTypes,
		DependencyIndexes: file_frame_proto_frame_proto_depIdxs,
		MessageInfos:      file_frame_proto_frame_proto_msgTypes,
	}.Build()
	File_frame_proto_frame_proto = out.File
	file_frame_proto_frame_proto_rawDesc = nil
	file_frame_proto_frame_proto_goTypes = nil
	file_frame_proto_frame_proto_depIdxs = nil
}
