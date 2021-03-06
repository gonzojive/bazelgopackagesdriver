// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.20.0--rc2
// source: proto/gopackagesdriver.proto

package gopackagesdriverpb

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

type CheckStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckStatusRequest) Reset() {
	*x = CheckStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gopackagesdriver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckStatusRequest) ProtoMessage() {}

func (x *CheckStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gopackagesdriver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckStatusRequest.ProtoReflect.Descriptor instead.
func (*CheckStatusRequest) Descriptor() ([]byte, []int) {
	return file_proto_gopackagesdriver_proto_rawDescGZIP(), []int{0}
}

type CheckStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DebugMessage string `protobuf:"bytes,1,opt,name=debug_message,json=debugMessage,proto3" json:"debug_message,omitempty"`
}

func (x *CheckStatusResponse) Reset() {
	*x = CheckStatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gopackagesdriver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckStatusResponse) ProtoMessage() {}

func (x *CheckStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gopackagesdriver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckStatusResponse.ProtoReflect.Descriptor instead.
func (*CheckStatusResponse) Descriptor() ([]byte, []int) {
	return file_proto_gopackagesdriver_proto_rawDescGZIP(), []int{1}
}

func (x *CheckStatusResponse) GetDebugMessage() string {
	if x != nil {
		return x.DebugMessage
	}
	return ""
}

type LoadPackagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queries   []string   `protobuf:"bytes,1,rep,name=queries,proto3" json:"queries,omitempty"`
	LoadMode  uint64     `protobuf:"varint,2,opt,name=load_mode,json=loadMode,proto3" json:"load_mode,omitempty"`
	EnvParams *EnvParams `protobuf:"bytes,3,opt,name=env_params,json=envParams,proto3" json:"env_params,omitempty"`
}

func (x *LoadPackagesRequest) Reset() {
	*x = LoadPackagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gopackagesdriver_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadPackagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadPackagesRequest) ProtoMessage() {}

func (x *LoadPackagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gopackagesdriver_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadPackagesRequest.ProtoReflect.Descriptor instead.
func (*LoadPackagesRequest) Descriptor() ([]byte, []int) {
	return file_proto_gopackagesdriver_proto_rawDescGZIP(), []int{2}
}

func (x *LoadPackagesRequest) GetQueries() []string {
	if x != nil {
		return x.Queries
	}
	return nil
}

func (x *LoadPackagesRequest) GetLoadMode() uint64 {
	if x != nil {
		return x.LoadMode
	}
	return 0
}

func (x *LoadPackagesRequest) GetEnvParams() *EnvParams {
	if x != nil {
		return x.EnvParams
	}
	return nil
}

type EnvParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RulesGoRepositoryName string   `protobuf:"bytes,1,opt,name=rules_go_repository_name,json=rulesGoRepositoryName,proto3" json:"rules_go_repository_name,omitempty"`
	BazelBin              string   `protobuf:"bytes,2,opt,name=bazel_bin,json=bazelBin,proto3" json:"bazel_bin,omitempty"`
	BazelFlags            []string `protobuf:"bytes,3,rep,name=bazel_flags,json=bazelFlags,proto3" json:"bazel_flags,omitempty"`
	BazelQueryFlags       []string `protobuf:"bytes,4,rep,name=bazel_query_flags,json=bazelQueryFlags,proto3" json:"bazel_query_flags,omitempty"`
	BazelQueryScope       string   `protobuf:"bytes,5,opt,name=bazel_query_scope,json=bazelQueryScope,proto3" json:"bazel_query_scope,omitempty"`
	BazelBuildFlags       []string `protobuf:"bytes,6,rep,name=bazel_build_flags,json=bazelBuildFlags,proto3" json:"bazel_build_flags,omitempty"`
	WorkspaceRoot         string   `protobuf:"bytes,7,opt,name=workspace_root,json=workspaceRoot,proto3" json:"workspace_root,omitempty"`
}

func (x *EnvParams) Reset() {
	*x = EnvParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gopackagesdriver_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvParams) ProtoMessage() {}

func (x *EnvParams) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gopackagesdriver_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvParams.ProtoReflect.Descriptor instead.
func (*EnvParams) Descriptor() ([]byte, []int) {
	return file_proto_gopackagesdriver_proto_rawDescGZIP(), []int{3}
}

func (x *EnvParams) GetRulesGoRepositoryName() string {
	if x != nil {
		return x.RulesGoRepositoryName
	}
	return ""
}

func (x *EnvParams) GetBazelBin() string {
	if x != nil {
		return x.BazelBin
	}
	return ""
}

func (x *EnvParams) GetBazelFlags() []string {
	if x != nil {
		return x.BazelFlags
	}
	return nil
}

func (x *EnvParams) GetBazelQueryFlags() []string {
	if x != nil {
		return x.BazelQueryFlags
	}
	return nil
}

func (x *EnvParams) GetBazelQueryScope() string {
	if x != nil {
		return x.BazelQueryScope
	}
	return ""
}

func (x *EnvParams) GetBazelBuildFlags() []string {
	if x != nil {
		return x.BazelBuildFlags
	}
	return nil
}

func (x *EnvParams) GetWorkspaceRoot() string {
	if x != nil {
		return x.WorkspaceRoot
	}
	return ""
}

type LoadPackagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RawJson string `protobuf:"bytes,1,opt,name=raw_json,json=rawJson,proto3" json:"raw_json,omitempty"`
}

func (x *LoadPackagesResponse) Reset() {
	*x = LoadPackagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gopackagesdriver_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadPackagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadPackagesResponse) ProtoMessage() {}

func (x *LoadPackagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gopackagesdriver_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadPackagesResponse.ProtoReflect.Descriptor instead.
func (*LoadPackagesResponse) Descriptor() ([]byte, []int) {
	return file_proto_gopackagesdriver_proto_rawDescGZIP(), []int{4}
}

func (x *LoadPackagesResponse) GetRawJson() string {
	if x != nil {
		return x.RawJson
	}
	return ""
}

type ProcessInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pid         int64  `protobuf:"varint,1,opt,name=pid,proto3" json:"pid,omitempty"`
	GrpcAddress string `protobuf:"bytes,2,opt,name=grpc_address,json=grpcAddress,proto3" json:"grpc_address,omitempty"`
}

func (x *ProcessInfo) Reset() {
	*x = ProcessInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gopackagesdriver_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessInfo) ProtoMessage() {}

func (x *ProcessInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gopackagesdriver_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessInfo.ProtoReflect.Descriptor instead.
func (*ProcessInfo) Descriptor() ([]byte, []int) {
	return file_proto_gopackagesdriver_proto_rawDescGZIP(), []int{5}
}

func (x *ProcessInfo) GetPid() int64 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *ProcessInfo) GetGrpcAddress() string {
	if x != nil {
		return x.GrpcAddress
	}
	return ""
}

var File_proto_gopackagesdriver_proto protoreflect.FileDescriptor

var file_proto_gopackagesdriver_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67,
	0x65, 0x73, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10,
	0x67, 0x6f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72,
	0x22, 0x14, 0x0a, 0x12, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x0a, 0x13, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a,
	0x0d, 0x64, 0x65, 0x62, 0x75, 0x67, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65, 0x62, 0x75, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x88, 0x01, 0x0a, 0x13, 0x4c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75,
	0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x71, 0x75, 0x65,
	0x72, 0x69, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6d, 0x6f, 0x64,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x6f, 0x64,
	0x65, 0x12, 0x3a, 0x0a, 0x0a, 0x65, 0x6e, 0x76, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67,
	0x65, 0x73, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6e, 0x76, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x52, 0x09, 0x65, 0x6e, 0x76, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0xad, 0x02,
	0x0a, 0x09, 0x45, 0x6e, 0x76, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x37, 0x0a, 0x18, 0x72,
	0x75, 0x6c, 0x65, 0x73, 0x5f, 0x67, 0x6f, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f,
	0x72, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72,
	0x75, 0x6c, 0x65, 0x73, 0x47, 0x6f, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x5f, 0x62, 0x69,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x42, 0x69,
	0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x46, 0x6c, 0x61,
	0x67, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x5f, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0f, 0x62,
	0x61, 0x7a, 0x65, 0x6c, 0x51, 0x75, 0x65, 0x72, 0x79, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x2a,
	0x0a, 0x11, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x62, 0x61, 0x7a, 0x65, 0x6c,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x62, 0x61,
	0x7a, 0x65, 0x6c, 0x5f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0f, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x42, 0x75, 0x69, 0x6c,
	0x64, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x22, 0x31, 0x0a,
	0x14, 0x4c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x61, 0x77, 0x5f, 0x6a, 0x73, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x61, 0x77, 0x4a, 0x73, 0x6f, 0x6e,
	0x22, 0x42, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x70, 0x69,
	0x64, 0x12, 0x21, 0x0a, 0x0c, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x67, 0x72, 0x70, 0x63, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x32, 0xd8, 0x01, 0x0a, 0x17, 0x47, 0x6f, 0x50, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x73, 0x44, 0x72, 0x69, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x5f, 0x0a, 0x0c, 0x4c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73,
	0x12, 0x25, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x64, 0x72, 0x69,
	0x76, 0x65, 0x72, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x63, 0x6b,
	0x61, 0x67, 0x65, 0x73, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x50,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x5c, 0x0a, 0x0b, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x24, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x64, 0x72, 0x69,
	0x76, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x67, 0x6f, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x73, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f,
	0x6e, 0x7a, 0x6f, 0x6a, 0x69, 0x76, 0x65, 0x2f, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x67, 0x6f, 0x70,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x64, 0x72,
	0x69, 0x76, 0x65, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_gopackagesdriver_proto_rawDescOnce sync.Once
	file_proto_gopackagesdriver_proto_rawDescData = file_proto_gopackagesdriver_proto_rawDesc
)

func file_proto_gopackagesdriver_proto_rawDescGZIP() []byte {
	file_proto_gopackagesdriver_proto_rawDescOnce.Do(func() {
		file_proto_gopackagesdriver_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_gopackagesdriver_proto_rawDescData)
	})
	return file_proto_gopackagesdriver_proto_rawDescData
}

var file_proto_gopackagesdriver_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_gopackagesdriver_proto_goTypes = []interface{}{
	(*CheckStatusRequest)(nil),   // 0: gopackagesdriver.CheckStatusRequest
	(*CheckStatusResponse)(nil),  // 1: gopackagesdriver.CheckStatusResponse
	(*LoadPackagesRequest)(nil),  // 2: gopackagesdriver.LoadPackagesRequest
	(*EnvParams)(nil),            // 3: gopackagesdriver.EnvParams
	(*LoadPackagesResponse)(nil), // 4: gopackagesdriver.LoadPackagesResponse
	(*ProcessInfo)(nil),          // 5: gopackagesdriver.ProcessInfo
}
var file_proto_gopackagesdriver_proto_depIdxs = []int32{
	3, // 0: gopackagesdriver.LoadPackagesRequest.env_params:type_name -> gopackagesdriver.EnvParams
	2, // 1: gopackagesdriver.GoPackagesDriverService.LoadPackages:input_type -> gopackagesdriver.LoadPackagesRequest
	0, // 2: gopackagesdriver.GoPackagesDriverService.CheckStatus:input_type -> gopackagesdriver.CheckStatusRequest
	4, // 3: gopackagesdriver.GoPackagesDriverService.LoadPackages:output_type -> gopackagesdriver.LoadPackagesResponse
	1, // 4: gopackagesdriver.GoPackagesDriverService.CheckStatus:output_type -> gopackagesdriver.CheckStatusResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_gopackagesdriver_proto_init() }
func file_proto_gopackagesdriver_proto_init() {
	if File_proto_gopackagesdriver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_gopackagesdriver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckStatusRequest); i {
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
		file_proto_gopackagesdriver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckStatusResponse); i {
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
		file_proto_gopackagesdriver_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadPackagesRequest); i {
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
		file_proto_gopackagesdriver_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvParams); i {
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
		file_proto_gopackagesdriver_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadPackagesResponse); i {
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
		file_proto_gopackagesdriver_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessInfo); i {
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
			RawDescriptor: file_proto_gopackagesdriver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_gopackagesdriver_proto_goTypes,
		DependencyIndexes: file_proto_gopackagesdriver_proto_depIdxs,
		MessageInfos:      file_proto_gopackagesdriver_proto_msgTypes,
	}.Build()
	File_proto_gopackagesdriver_proto = out.File
	file_proto_gopackagesdriver_proto_rawDesc = nil
	file_proto_gopackagesdriver_proto_goTypes = nil
	file_proto_gopackagesdriver_proto_depIdxs = nil
}
