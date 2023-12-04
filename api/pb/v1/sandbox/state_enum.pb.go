// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: sandbox/state_enum.proto

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

type State int32

const (
	State_CE  State = 0
	State_AC  State = 1
	State_WA  State = 2
	State_RE  State = 3
	State_TLE State = 4
	State_MLE State = 5
	State_UE  State = 6
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "CE",
		1: "AC",
		2: "WA",
		3: "RE",
		4: "TLE",
		5: "MLE",
		6: "UE",
	}
	State_value = map[string]int32{
		"CE":  0,
		"AC":  1,
		"WA":  2,
		"RE":  3,
		"TLE": 4,
		"MLE": 5,
		"UE":  6,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_sandbox_state_enum_proto_enumTypes[0].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_sandbox_state_enum_proto_enumTypes[0]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_sandbox_state_enum_proto_rawDescGZIP(), []int{0}
}

var File_sandbox_state_enum_proto protoreflect.FileDescriptor

var file_sandbox_state_enum_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f,
	0x65, 0x6e, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x76, 0x31, 0x2e, 0x73,
	0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x2a, 0x41, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x06, 0x0a, 0x02, 0x43, 0x45, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x41, 0x43, 0x10, 0x01, 0x12,
	0x06, 0x0a, 0x02, 0x57, 0x41, 0x10, 0x02, 0x12, 0x06, 0x0a, 0x02, 0x52, 0x45, 0x10, 0x03, 0x12,
	0x07, 0x0a, 0x03, 0x54, 0x4c, 0x45, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x4d, 0x4c, 0x45, 0x10,
	0x05, 0x12, 0x06, 0x0a, 0x02, 0x55, 0x45, 0x10, 0x06, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x73, 0x71, 0x74, 0x74, 0x2f, 0x73, 0x62,
	0x2d, 0x6a, 0x75, 0x64, 0x67, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x2f, 0x76,
	0x31, 0x2f, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x3b, 0x70, 0x62, 0x5f, 0x73, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sandbox_state_enum_proto_rawDescOnce sync.Once
	file_sandbox_state_enum_proto_rawDescData = file_sandbox_state_enum_proto_rawDesc
)

func file_sandbox_state_enum_proto_rawDescGZIP() []byte {
	file_sandbox_state_enum_proto_rawDescOnce.Do(func() {
		file_sandbox_state_enum_proto_rawDescData = protoimpl.X.CompressGZIP(file_sandbox_state_enum_proto_rawDescData)
	})
	return file_sandbox_state_enum_proto_rawDescData
}

var file_sandbox_state_enum_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_sandbox_state_enum_proto_goTypes = []interface{}{
	(State)(0), // 0: v1.sandbox.State
}
var file_sandbox_state_enum_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sandbox_state_enum_proto_init() }
func file_sandbox_state_enum_proto_init() {
	if File_sandbox_state_enum_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sandbox_state_enum_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sandbox_state_enum_proto_goTypes,
		DependencyIndexes: file_sandbox_state_enum_proto_depIdxs,
		EnumInfos:         file_sandbox_state_enum_proto_enumTypes,
	}.Build()
	File_sandbox_state_enum_proto = out.File
	file_sandbox_state_enum_proto_rawDesc = nil
	file_sandbox_state_enum_proto_goTypes = nil
	file_sandbox_state_enum_proto_depIdxs = nil
}
