// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.3
// source: internal/idl/transport.proto

package transport

import (
	context "context"
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

type Meta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *Meta) Reset() {
	*x = Meta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_idl_transport_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Meta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Meta) ProtoMessage() {}

func (x *Meta) ProtoReflect() protoreflect.Message {
	mi := &file_internal_idl_transport_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Meta.ProtoReflect.Descriptor instead.
func (*Meta) Descriptor() ([]byte, []int) {
	return file_internal_idl_transport_proto_rawDescGZIP(), []int{0}
}

func (x *Meta) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type Transport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr    string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Session int64  `protobuf:"varint,2,opt,name=session,proto3" json:"session,omitempty"`
	Meta    *Meta  `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
	Cmd     string `protobuf:"bytes,4,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Msg     []byte `protobuf:"bytes,5,opt,name=msg,proto3" json:"msg,omitempty"`
	Traces  []byte `protobuf:"bytes,7,opt,name=traces,proto3" json:"traces,omitempty"`
}

func (x *Transport) Reset() {
	*x = Transport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_idl_transport_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transport) ProtoMessage() {}

func (x *Transport) ProtoReflect() protoreflect.Message {
	mi := &file_internal_idl_transport_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transport.ProtoReflect.Descriptor instead.
func (*Transport) Descriptor() ([]byte, []int) {
	return file_internal_idl_transport_proto_rawDescGZIP(), []int{1}
}

func (x *Transport) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Transport) GetSession() int64 {
	if x != nil {
		return x.Session
	}
	return 0
}

func (x *Transport) GetMeta() *Meta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *Transport) GetCmd() string {
	if x != nil {
		return x.Cmd
	}
	return ""
}

func (x *Transport) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *Transport) GetTraces() []byte {
	if x != nil {
		return x.Traces
	}
	return nil
}

var File_internal_idl_transport_proto protoreflect.FileDescriptor

var file_internal_idl_transport_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69, 0x64, 0x6c, 0x2f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x1a, 0x0a, 0x04, 0x4d, 0x65, 0x74,
	0x61, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x75, 0x75, 0x69, 0x64, 0x22, 0x9a, 0x01, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70,
	0x6f, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x23, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x72,
	0x61, 0x63, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x74, 0x72, 0x61, 0x63,
	0x65, 0x73, 0x32, 0x46, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x14,
	0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x1a, 0x14, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x6c, 0x73, 0x77, 0x2f, 0x69, 0x6b,
	0x75, 0x6e, 0x65, 0x74, 0x2f, 0x6b, 0x69, 0x74, 0x65, 0x78, 0x5f, 0x67, 0x65, 0x6e, 0x2f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_idl_transport_proto_rawDescOnce sync.Once
	file_internal_idl_transport_proto_rawDescData = file_internal_idl_transport_proto_rawDesc
)

func file_internal_idl_transport_proto_rawDescGZIP() []byte {
	file_internal_idl_transport_proto_rawDescOnce.Do(func() {
		file_internal_idl_transport_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_idl_transport_proto_rawDescData)
	})
	return file_internal_idl_transport_proto_rawDescData
}

var file_internal_idl_transport_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_idl_transport_proto_goTypes = []interface{}{
	(*Meta)(nil),      // 0: transport.Meta
	(*Transport)(nil), // 1: transport.Transport
}
var file_internal_idl_transport_proto_depIdxs = []int32{
	0, // 0: transport.Transport.meta:type_name -> transport.Meta
	1, // 1: transport.TransportService.Call:input_type -> transport.Transport
	1, // 2: transport.TransportService.Call:output_type -> transport.Transport
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_idl_transport_proto_init() }
func file_internal_idl_transport_proto_init() {
	if File_internal_idl_transport_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_idl_transport_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Meta); i {
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
		file_internal_idl_transport_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transport); i {
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
			RawDescriptor: file_internal_idl_transport_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_idl_transport_proto_goTypes,
		DependencyIndexes: file_internal_idl_transport_proto_depIdxs,
		MessageInfos:      file_internal_idl_transport_proto_msgTypes,
	}.Build()
	File_internal_idl_transport_proto = out.File
	file_internal_idl_transport_proto_rawDesc = nil
	file_internal_idl_transport_proto_goTypes = nil
	file_internal_idl_transport_proto_depIdxs = nil
}

var _ context.Context

// Code generated by Kitex v0.10.3. DO NOT EDIT.

type TransportService interface {
	Call(ctx context.Context, req *Transport) (res *Transport, err error)
}
