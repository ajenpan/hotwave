// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: service/lobby/proto/lobby.proto

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

type UserPropsInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UserPropsInfoRequest) Reset() {
	*x = UserPropsInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_lobby_proto_lobby_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserPropsInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserPropsInfoRequest) ProtoMessage() {}

func (x *UserPropsInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_lobby_proto_lobby_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserPropsInfoRequest.ProtoReflect.Descriptor instead.
func (*UserPropsInfoRequest) Descriptor() ([]byte, []int) {
	return file_service_lobby_proto_lobby_proto_rawDescGZIP(), []int{0}
}

type PropsInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Count int32  `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *PropsInfo) Reset() {
	*x = PropsInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_lobby_proto_lobby_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PropsInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PropsInfo) ProtoMessage() {}

func (x *PropsInfo) ProtoReflect() protoreflect.Message {
	mi := &file_service_lobby_proto_lobby_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PropsInfo.ProtoReflect.Descriptor instead.
func (*PropsInfo) Descriptor() ([]byte, []int) {
	return file_service_lobby_proto_lobby_proto_rawDescGZIP(), []int{1}
}

func (x *PropsInfo) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PropsInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PropsInfo) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type UserPropsInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Props []*PropsInfo `protobuf:"bytes,1,rep,name=props,proto3" json:"props,omitempty"`
}

func (x *UserPropsInfoResponse) Reset() {
	*x = UserPropsInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_lobby_proto_lobby_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserPropsInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserPropsInfoResponse) ProtoMessage() {}

func (x *UserPropsInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_lobby_proto_lobby_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserPropsInfoResponse.ProtoReflect.Descriptor instead.
func (*UserPropsInfoResponse) Descriptor() ([]byte, []int) {
	return file_service_lobby_proto_lobby_proto_rawDescGZIP(), []int{2}
}

func (x *UserPropsInfoResponse) GetProps() []*PropsInfo {
	if x != nil {
		return x.Props
	}
	return nil
}

type UserMatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameName string `protobuf:"bytes,1,opt,name=game_name,json=gameName,proto3" json:"game_name,omitempty"`
	RoomId   int32  `protobuf:"varint,2,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
}

func (x *UserMatchRequest) Reset() {
	*x = UserMatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_lobby_proto_lobby_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserMatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserMatchRequest) ProtoMessage() {}

func (x *UserMatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_lobby_proto_lobby_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserMatchRequest.ProtoReflect.Descriptor instead.
func (*UserMatchRequest) Descriptor() ([]byte, []int) {
	return file_service_lobby_proto_lobby_proto_rawDescGZIP(), []int{3}
}

func (x *UserMatchRequest) GetGameName() string {
	if x != nil {
		return x.GameName
	}
	return ""
}

func (x *UserMatchRequest) GetRoomId() int32 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

type UserMatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UserMatchResponse) Reset() {
	*x = UserMatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_lobby_proto_lobby_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserMatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserMatchResponse) ProtoMessage() {}

func (x *UserMatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_lobby_proto_lobby_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserMatchResponse.ProtoReflect.Descriptor instead.
func (*UserMatchResponse) Descriptor() ([]byte, []int) {
	return file_service_lobby_proto_lobby_proto_rawDescGZIP(), []int{4}
}

type PlayerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid    int64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	SeatId int32 `protobuf:"varint,2,opt,name=seat_id,json=seatId,proto3" json:"seat_id,omitempty"`
}

func (x *PlayerInfo) Reset() {
	*x = PlayerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_lobby_proto_lobby_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerInfo) ProtoMessage() {}

func (x *PlayerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_service_lobby_proto_lobby_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerInfo.ProtoReflect.Descriptor instead.
func (*PlayerInfo) Descriptor() ([]byte, []int) {
	return file_service_lobby_proto_lobby_proto_rawDescGZIP(), []int{5}
}

func (x *PlayerInfo) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *PlayerInfo) GetSeatId() int32 {
	if x != nil {
		return x.SeatId
	}
	return 0
}

type UserGameStartNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BattleId string        `protobuf:"bytes,1,opt,name=battle_id,json=battleId,proto3" json:"battle_id,omitempty"`
	Players  []*PlayerInfo `protobuf:"bytes,2,rep,name=players,proto3" json:"players,omitempty"`
	Errcode  int32         `protobuf:"varint,3,opt,name=errcode,proto3" json:"errcode,omitempty"`
	Errmsg   string        `protobuf:"bytes,4,opt,name=errmsg,proto3" json:"errmsg,omitempty"`
}

func (x *UserGameStartNotify) Reset() {
	*x = UserGameStartNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_lobby_proto_lobby_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserGameStartNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserGameStartNotify) ProtoMessage() {}

func (x *UserGameStartNotify) ProtoReflect() protoreflect.Message {
	mi := &file_service_lobby_proto_lobby_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserGameStartNotify.ProtoReflect.Descriptor instead.
func (*UserGameStartNotify) Descriptor() ([]byte, []int) {
	return file_service_lobby_proto_lobby_proto_rawDescGZIP(), []int{6}
}

func (x *UserGameStartNotify) GetBattleId() string {
	if x != nil {
		return x.BattleId
	}
	return ""
}

func (x *UserGameStartNotify) GetPlayers() []*PlayerInfo {
	if x != nil {
		return x.Players
	}
	return nil
}

func (x *UserGameStartNotify) GetErrcode() int32 {
	if x != nil {
		return x.Errcode
	}
	return 0
}

func (x *UserGameStartNotify) GetErrmsg() string {
	if x != nil {
		return x.Errmsg
	}
	return ""
}

type UserGameOverNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UserGameOverNotify) Reset() {
	*x = UserGameOverNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_lobby_proto_lobby_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserGameOverNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserGameOverNotify) ProtoMessage() {}

func (x *UserGameOverNotify) ProtoReflect() protoreflect.Message {
	mi := &file_service_lobby_proto_lobby_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserGameOverNotify.ProtoReflect.Descriptor instead.
func (*UserGameOverNotify) Descriptor() ([]byte, []int) {
	return file_service_lobby_proto_lobby_proto_rawDescGZIP(), []int{7}
}

var File_service_lobby_proto_lobby_proto protoreflect.FileDescriptor

var file_service_lobby_proto_lobby_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x07, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x22, 0x16, 0x0a, 0x14, 0x55, 0x73,
	0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x45, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x41, 0x0a, 0x15, 0x55, 0x73, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x28, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x50, 0x72, 0x6f, 0x70,
	0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x22, 0x48, 0x0a, 0x10,
	0x55, 0x73, 0x65, 0x72, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x61, 0x6d, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a,
	0x07, 0x72, 0x6f, 0x6f, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x22, 0x13, 0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x37, 0x0a, 0x0a, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x73,
	0x65, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x65,
	0x61, 0x74, 0x49, 0x64, 0x22, 0x93, 0x01, 0x0a, 0x13, 0x55, 0x73, 0x65, 0x72, 0x47, 0x61, 0x6d,
	0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x1b, 0x0a, 0x09,
	0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x2d, 0x0a, 0x07, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x67, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x07, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x72, 0x72, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6d, 0x73, 0x67, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6d, 0x73, 0x67, 0x22, 0x14, 0x0a, 0x12, 0x55, 0x73,
	0x65, 0x72, 0x47, 0x61, 0x6d, 0x65, 0x4f, 0x76, 0x65, 0x72, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x42, 0x15, 0x5a, 0x13, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x6c, 0x6f, 0x62, 0x62,
	0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_lobby_proto_lobby_proto_rawDescOnce sync.Once
	file_service_lobby_proto_lobby_proto_rawDescData = file_service_lobby_proto_lobby_proto_rawDesc
)

func file_service_lobby_proto_lobby_proto_rawDescGZIP() []byte {
	file_service_lobby_proto_lobby_proto_rawDescOnce.Do(func() {
		file_service_lobby_proto_lobby_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_lobby_proto_lobby_proto_rawDescData)
	})
	return file_service_lobby_proto_lobby_proto_rawDescData
}

var file_service_lobby_proto_lobby_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_service_lobby_proto_lobby_proto_goTypes = []interface{}{
	(*UserPropsInfoRequest)(nil),  // 0: gateway.UserPropsInfoRequest
	(*PropsInfo)(nil),             // 1: gateway.PropsInfo
	(*UserPropsInfoResponse)(nil), // 2: gateway.UserPropsInfoResponse
	(*UserMatchRequest)(nil),      // 3: gateway.UserMatchRequest
	(*UserMatchResponse)(nil),     // 4: gateway.UserMatchResponse
	(*PlayerInfo)(nil),            // 5: gateway.PlayerInfo
	(*UserGameStartNotify)(nil),   // 6: gateway.UserGameStartNotify
	(*UserGameOverNotify)(nil),    // 7: gateway.UserGameOverNotify
}
var file_service_lobby_proto_lobby_proto_depIdxs = []int32{
	1, // 0: gateway.UserPropsInfoResponse.props:type_name -> gateway.PropsInfo
	5, // 1: gateway.UserGameStartNotify.players:type_name -> gateway.PlayerInfo
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_service_lobby_proto_lobby_proto_init() }
func file_service_lobby_proto_lobby_proto_init() {
	if File_service_lobby_proto_lobby_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_lobby_proto_lobby_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserPropsInfoRequest); i {
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
		file_service_lobby_proto_lobby_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PropsInfo); i {
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
		file_service_lobby_proto_lobby_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserPropsInfoResponse); i {
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
		file_service_lobby_proto_lobby_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserMatchRequest); i {
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
		file_service_lobby_proto_lobby_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserMatchResponse); i {
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
		file_service_lobby_proto_lobby_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerInfo); i {
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
		file_service_lobby_proto_lobby_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserGameStartNotify); i {
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
		file_service_lobby_proto_lobby_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserGameOverNotify); i {
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
			RawDescriptor: file_service_lobby_proto_lobby_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_service_lobby_proto_lobby_proto_goTypes,
		DependencyIndexes: file_service_lobby_proto_lobby_proto_depIdxs,
		MessageInfos:      file_service_lobby_proto_lobby_proto_msgTypes,
	}.Build()
	File_service_lobby_proto_lobby_proto = out.File
	file_service_lobby_proto_lobby_proto_rawDesc = nil
	file_service_lobby_proto_lobby_proto_goTypes = nil
	file_service_lobby_proto_lobby_proto_depIdxs = nil
}