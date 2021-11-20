// Upvest Confidential
//
// Copyright 2020 - 2021 Upvest GmbH. All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains the property
// of Upvest GmbH. The intellectual and technical concepts contained herein
// are proprietary to Upvest GmbH and are protected by trade secret or
// copyright law. Dissemination of this information or reproduction of this
// material is strictly forbidden unless prior written permission is
// obtained from Upvest GmbH.

// registry.group_id: data-and-reporting
// registry.artifact_id: reportservice
// registry.title: Report Service
// registry.description: The Report service is responsible for generating reports
// registry.labels: domain

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: pdcl.proto

package sentinelpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PublishRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cid string `protobuf:"bytes,1,opt,name=cid,proto3" json:"cid,omitempty"`
}

func (x *PublishRequest) Reset() {
	*x = PublishRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sentinel_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishRequest) ProtoMessage() {}

func (x *PublishRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sentinel_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishRequest.ProtoReflect.Descriptor instead.
func (*PublishRequest) Descriptor() ([]byte, []int) {
	return file_sentinel_proto_rawDescGZIP(), []int{0}
}

func (x *PublishRequest) GetCid() string {
	if x != nil {
		return x.Cid
	}
	return ""
}

type PublishResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PublishResponse) Reset() {
	*x = PublishResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sentinel_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishResponse) ProtoMessage() {}

func (x *PublishResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sentinel_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishResponse.ProtoReflect.Descriptor instead.
func (*PublishResponse) Descriptor() ([]byte, []int) {
	return file_sentinel_proto_rawDescGZIP(), []int{1}
}

type Commit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Created           *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=created,proto3" json:"created,omitempty"`
	PreviousCommitCid string                 `protobuf:"bytes,2,opt,name=previous_commit_cid,json=previousCommitCid,proto3" json:"previous_commit_cid,omitempty"`
	MessagesCids      []string               `protobuf:"bytes,3,rep,name=messages_cids,json=messagesCids,proto3" json:"messages_cids,omitempty"`
}

func (x *Commit) Reset() {
	*x = Commit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sentinel_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commit) ProtoMessage() {}

func (x *Commit) ProtoReflect() protoreflect.Message {
	mi := &file_sentinel_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commit.ProtoReflect.Descriptor instead.
func (*Commit) Descriptor() ([]byte, []int) {
	return file_sentinel_proto_rawDescGZIP(), []int{2}
}

func (x *Commit) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *Commit) GetPreviousCommitCid() string {
	if x != nil {
		return x.PreviousCommitCid
	}
	return ""
}

func (x *Commit) GetMessagesCids() []string {
	if x != nil {
		return x.MessagesCids
	}
	return nil
}

var File_sentinel_proto protoreflect.FileDescriptor

var file_sentinel_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a, 0x0e, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x63, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x69, 0x64, 0x22,
	0x11, 0x0a, 0x0f, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x93, 0x01, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x34, 0x0a,
	0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x63, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x11, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x43, 0x69, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x5f,
	0x63, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x43, 0x69, 0x64, 0x73, 0x32, 0x4a, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x74,
	0x69, 0x6e, 0x65, 0x6c, 0x12, 0x3e, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x12,
	0x18, 0x2e, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x65, 0x6e, 0x74,
	0x69, 0x6e, 0x65, 0x6c, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x3b, 0x73, 0x65, 0x6e, 0x74, 0x69,
	0x6e, 0x65, 0x6c, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sentinel_proto_rawDescOnce sync.Once
	file_sentinel_proto_rawDescData = file_sentinel_proto_rawDesc
)

func file_sentinel_proto_rawDescGZIP() []byte {
	file_sentinel_proto_rawDescOnce.Do(func() {
		file_sentinel_proto_rawDescData = protoimpl.X.CompressGZIP(file_sentinel_proto_rawDescData)
	})
	return file_sentinel_proto_rawDescData
}

var (
	file_sentinel_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
	file_sentinel_proto_goTypes  = []interface{}{
		(*PublishRequest)(nil),        // 0: sentinel.PublishRequest
		(*PublishResponse)(nil),       // 1: sentinel.PublishResponse
		(*Commit)(nil),                // 2: sentinel.Commit
		(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	}
)

var file_sentinel_proto_depIdxs = []int32{
	3, // 0: sentinel.Commit.created:type_name -> google.protobuf.Timestamp
	0, // 1: sentinel.Sentinel.Publish:input_type -> sentinel.PublishRequest
	1, // 2: sentinel.Sentinel.Publish:output_type -> sentinel.PublishResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_sentinel_proto_init() }
func file_sentinel_proto_init() {
	if File_sentinel_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sentinel_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishRequest); i {
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
		file_sentinel_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishResponse); i {
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
		file_sentinel_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Commit); i {
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
			RawDescriptor: file_sentinel_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sentinel_proto_goTypes,
		DependencyIndexes: file_sentinel_proto_depIdxs,
		MessageInfos:      file_sentinel_proto_msgTypes,
	}.Build()
	File_sentinel_proto = out.File
	file_sentinel_proto_rawDesc = nil
	file_sentinel_proto_goTypes = nil
	file_sentinel_proto_depIdxs = nil
}
