// Code generated by protoc-gen-go.
// source: msp/identities.proto
// DO NOT EDIT!

/*
Package msp is a generated protocol buffer package.

It is generated from these files:
	msp/identities.proto
	msp/msp_config.proto
	msp/msp_principal.proto

It has these top-level messages:
	SerializedIdentity
	MSPConfig
	FabricMSPConfig
	SigningIdentityInfo
	KeyInfo
	FabricOUIdentifier
	MSPPrincipal
	OrganizationUnit
	MSPRole
*/
package msp

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// This struct represents an Identity
// (with its MSP identifier) to be used
// to serialize it and deserialize it
type SerializedIdentity struct {
	// The identifier of the associated membership service provider
	Mspid string `protobuf:"bytes,1,opt,name=mspid" json:"mspid,omitempty"`
	// the Identity, serialized according to the rules of its MPS
	IdBytes []byte `protobuf:"bytes,2,opt,name=id_bytes,json=idBytes,proto3" json:"id_bytes,omitempty"`
}

func (m *SerializedIdentity) Reset()                    { *m = SerializedIdentity{} }
func (m *SerializedIdentity) String() string            { return proto.CompactTextString(m) }
func (*SerializedIdentity) ProtoMessage()               {}
func (*SerializedIdentity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*SerializedIdentity)(nil), "msp.SerializedIdentity")
}

func init() { proto.RegisterFile("msp/identities.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x4c, 0xcd, 0x3d, 0x0b, 0xc2, 0x30,
	0x10, 0xc6, 0x71, 0xaa, 0xf8, 0x16, 0x9c, 0x42, 0x87, 0xba, 0x55, 0xa7, 0x4e, 0xc9, 0xe0, 0x37,
	0x28, 0x38, 0x38, 0xb8, 0xd4, 0xcd, 0x45, 0x9a, 0xe6, 0x6c, 0x0f, 0x1a, 0x73, 0xe4, 0xe2, 0x50,
	0x3f, 0xbd, 0x68, 0x16, 0xb7, 0xe7, 0x0f, 0x3f, 0x78, 0x44, 0xee, 0x98, 0x34, 0x5a, 0x78, 0x46,
	0x8c, 0x08, 0xac, 0x28, 0xf8, 0xe8, 0xe5, 0xdc, 0x31, 0x1d, 0x4e, 0x42, 0x5e, 0x21, 0x60, 0x3b,
	0xe2, 0x1b, 0xec, 0x39, 0x91, 0x49, 0xe6, 0x62, 0xe1, 0x98, 0xd0, 0x16, 0x59, 0x99, 0x55, 0x9b,
	0x26, 0x85, 0xdc, 0x89, 0x35, 0xda, 0xbb, 0x99, 0x22, 0x70, 0x31, 0x2b, 0xb3, 0x6a, 0xdb, 0xac,
	0xd0, 0xd6, 0xdf, 0xac, 0x2f, 0x62, 0xef, 0x43, 0xaf, 0x86, 0x89, 0x20, 0x8c, 0x60, 0x7b, 0x08,
	0xea, 0xd1, 0x9a, 0x80, 0x5d, 0xfa, 0x62, 0xe5, 0x98, 0x6e, 0x55, 0x8f, 0x71, 0x78, 0x19, 0xd5,
	0x79, 0xa7, 0xff, 0xa4, 0x4e, 0x52, 0x27, 0xa9, 0x1d, 0x93, 0x59, 0xfe, 0xf6, 0xf1, 0x13, 0x00,
	0x00, 0xff, 0xff, 0xaf, 0xfd, 0x80, 0x9c, 0xb9, 0x00, 0x00, 0x00,
}
