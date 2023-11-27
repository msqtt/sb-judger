// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: sandbox/output_message.proto

package pb_sb

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

type Output struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CaseId        uint32 `protobuf:"varint,1,opt,name=case_id,json=caseId,proto3" json:"case_id,omitempty"`
	CpuTimeUsage  uint32 `protobuf:"varint,2,opt,name=cpu_time_usage,json=cpuTimeUsage,proto3" json:"cpu_time_usage,omitempty"`
	RealTimeUsage uint32 `protobuf:"varint,3,opt,name=real_time_usage,json=realTimeUsage,proto3" json:"real_time_usage,omitempty"`
	MemoryUsage   uint32 `protobuf:"varint,4,opt,name=memory_usage,json=memoryUsage,proto3" json:"memory_usage,omitempty"`
	Status        Status `protobuf:"varint,5,opt,name=status,proto3,enum=v1.sandbox.Status" json:"status,omitempty"`
	OutPut        string `protobuf:"bytes,6,opt,name=out_put,json=outPut,proto3" json:"out_put,omitempty"`
}

func (x *Output) Reset() {
	*x = Output{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sandbox_output_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Output) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Output) ProtoMessage() {}

func (x *Output) ProtoReflect() protoreflect.Message {
	mi := &file_sandbox_output_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Output.ProtoReflect.Descriptor instead.
func (*Output) Descriptor() ([]byte, []int) {
	return file_sandbox_output_message_proto_rawDescGZIP(), []int{0}
}

func (x *Output) GetCaseId() uint32 {
	if x != nil {
		return x.CaseId
	}
	return 0
}

func (x *Output) GetCpuTimeUsage() uint32 {
	if x != nil {
		return x.CpuTimeUsage
	}
	return 0
}

func (x *Output) GetRealTimeUsage() uint32 {
	if x != nil {
		return x.RealTimeUsage
	}
	return 0
}

func (x *Output) GetMemoryUsage() uint32 {
	if x != nil {
		return x.MemoryUsage
	}
	return 0
}

func (x *Output) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_AC
}

func (x *Output) GetOutPut() string {
	if x != nil {
		return x.OutPut
	}
	return ""
}

type CollectOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CaseOuts []*Output `protobuf:"bytes,1,rep,name=case_outs,json=caseOuts,proto3" json:"case_outs,omitempty"`
}

func (x *CollectOutput) Reset() {
	*x = CollectOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sandbox_output_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CollectOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CollectOutput) ProtoMessage() {}

func (x *CollectOutput) ProtoReflect() protoreflect.Message {
	mi := &file_sandbox_output_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CollectOutput.ProtoReflect.Descriptor instead.
func (*CollectOutput) Descriptor() ([]byte, []int) {
	return file_sandbox_output_message_proto_rawDescGZIP(), []int{1}
}

func (x *CollectOutput) GetCaseOuts() []*Output {
	if x != nil {
		return x.CaseOuts
	}
	return nil
}

var File_sandbox_output_message_proto protoreflect.FileDescriptor

var file_sandbox_output_message_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x76, 0x31, 0x2e, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x1a, 0x19, 0x73, 0x61, 0x6e, 0x64,
	0x62, 0x6f, 0x78, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd7, 0x01, 0x0a, 0x06, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x63, 0x61, 0x73, 0x65, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x63, 0x70, 0x75,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0c, 0x63, 0x70, 0x75, 0x54, 0x69, 0x6d, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x26, 0x0a, 0x0f, 0x72, 0x65, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x72, 0x65, 0x61, 0x6c, 0x54, 0x69,
	0x6d, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x6d, 0x6f, 0x72,
	0x79, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x6d,
	0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x76, 0x31, 0x2e,
	0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x6f, 0x75, 0x74, 0x5f, 0x70, 0x75,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x50, 0x75, 0x74, 0x22,
	0x40, 0x0a, 0x0d, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x12, 0x2f, 0x0a, 0x09, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x6f, 0x75, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78,
	0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x08, 0x63, 0x61, 0x73, 0x65, 0x4f, 0x75, 0x74,
	0x73, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6d, 0x73, 0x71, 0x74, 0x74, 0x2f, 0x73, 0x62, 0x2d, 0x6a, 0x75, 0x64, 0x67, 0x65, 0x72, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f,
	0x78, 0x3b, 0x70, 0x62, 0x5f, 0x73, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sandbox_output_message_proto_rawDescOnce sync.Once
	file_sandbox_output_message_proto_rawDescData = file_sandbox_output_message_proto_rawDesc
)

func file_sandbox_output_message_proto_rawDescGZIP() []byte {
	file_sandbox_output_message_proto_rawDescOnce.Do(func() {
		file_sandbox_output_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_sandbox_output_message_proto_rawDescData)
	})
	return file_sandbox_output_message_proto_rawDescData
}

var file_sandbox_output_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_sandbox_output_message_proto_goTypes = []interface{}{
	(*Output)(nil),        // 0: v1.sandbox.Output
	(*CollectOutput)(nil), // 1: v1.sandbox.CollectOutput
	(Status)(0),           // 2: v1.sandbox.Status
}
var file_sandbox_output_message_proto_depIdxs = []int32{
	2, // 0: v1.sandbox.Output.status:type_name -> v1.sandbox.Status
	0, // 1: v1.sandbox.CollectOutput.case_outs:type_name -> v1.sandbox.Output
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_sandbox_output_message_proto_init() }
func file_sandbox_output_message_proto_init() {
	if File_sandbox_output_message_proto != nil {
		return
	}
	file_sandbox_status_enum_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_sandbox_output_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Output); i {
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
		file_sandbox_output_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CollectOutput); i {
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
			RawDescriptor: file_sandbox_output_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sandbox_output_message_proto_goTypes,
		DependencyIndexes: file_sandbox_output_message_proto_depIdxs,
		MessageInfos:      file_sandbox_output_message_proto_msgTypes,
	}.Build()
	File_sandbox_output_message_proto = out.File
	file_sandbox_output_message_proto_rawDesc = nil
	file_sandbox_output_message_proto_goTypes = nil
	file_sandbox_output_message_proto_depIdxs = nil
}
