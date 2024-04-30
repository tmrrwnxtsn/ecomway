// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.6
// source: api/proto/shared/shared.proto

package shared

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OperationType int32

const (
	OperationType_PAYMENT OperationType = 0
	OperationType_PAYOUT  OperationType = 1
)

// Enum value maps for OperationType.
var (
	OperationType_name = map[int32]string{
		0: "PAYMENT",
		1: "PAYOUT",
	}
	OperationType_value = map[string]int32{
		"PAYMENT": 0,
		"PAYOUT":  1,
	}
)

func (x OperationType) Enum() *OperationType {
	p := new(OperationType)
	*p = x
	return p
}

func (x OperationType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OperationType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_shared_shared_proto_enumTypes[0].Descriptor()
}

func (OperationType) Type() protoreflect.EnumType {
	return &file_api_proto_shared_shared_proto_enumTypes[0]
}

func (x OperationType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OperationType.Descriptor instead.
func (OperationType) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{0}
}

type OperationExternalStatus int32

const (
	OperationExternalStatus_UNKNOWN OperationExternalStatus = 0
	OperationExternalStatus_PENDING OperationExternalStatus = 1
	OperationExternalStatus_SUCCESS OperationExternalStatus = 2
	OperationExternalStatus_FAILED  OperationExternalStatus = 3
)

// Enum value maps for OperationExternalStatus.
var (
	OperationExternalStatus_name = map[int32]string{
		0: "UNKNOWN",
		1: "PENDING",
		2: "SUCCESS",
		3: "FAILED",
	}
	OperationExternalStatus_value = map[string]int32{
		"UNKNOWN": 0,
		"PENDING": 1,
		"SUCCESS": 2,
		"FAILED":  3,
	}
)

func (x OperationExternalStatus) Enum() *OperationExternalStatus {
	p := new(OperationExternalStatus)
	*p = x
	return p
}

func (x OperationExternalStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OperationExternalStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_shared_shared_proto_enumTypes[1].Descriptor()
}

func (OperationExternalStatus) Type() protoreflect.EnumType {
	return &file_api_proto_shared_shared_proto_enumTypes[1]
}

func (x OperationExternalStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OperationExternalStatus.Descriptor instead.
func (OperationExternalStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{1}
}

type CommissionType int32

const (
	CommissionType_PERCENT  CommissionType = 0
	CommissionType_FIXED    CommissionType = 1
	CommissionType_COMBINED CommissionType = 2
	CommissionType_TEXT     CommissionType = 3
)

// Enum value maps for CommissionType.
var (
	CommissionType_name = map[int32]string{
		0: "PERCENT",
		1: "FIXED",
		2: "COMBINED",
		3: "TEXT",
	}
	CommissionType_value = map[string]int32{
		"PERCENT":  0,
		"FIXED":    1,
		"COMBINED": 2,
		"TEXT":     3,
	}
)

func (x CommissionType) Enum() *CommissionType {
	p := new(CommissionType)
	*p = x
	return p
}

func (x CommissionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommissionType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_shared_shared_proto_enumTypes[2].Descriptor()
}

func (CommissionType) Type() protoreflect.EnumType {
	return &file_api_proto_shared_shared_proto_enumTypes[2]
}

func (x CommissionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommissionType.Descriptor instead.
func (CommissionType) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{2}
}

type ToolType int32

const (
	ToolType_BANK_CARD ToolType = 0
	ToolType_WALLET    ToolType = 1
)

// Enum value maps for ToolType.
var (
	ToolType_name = map[int32]string{
		0: "BANK_CARD",
		1: "WALLET",
	}
	ToolType_value = map[string]int32{
		"BANK_CARD": 0,
		"WALLET":    1,
	}
)

func (x ToolType) Enum() *ToolType {
	p := new(ToolType)
	*p = x
	return p
}

func (x ToolType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ToolType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_shared_shared_proto_enumTypes[3].Descriptor()
}

func (ToolType) Type() protoreflect.EnumType {
	return &file_api_proto_shared_shared_proto_enumTypes[3]
}

func (x ToolType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ToolType.Descriptor instead.
func (ToolType) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{3}
}

type ToolStatus int32

const (
	ToolStatus_ACTIVE ToolStatus = 0
)

// Enum value maps for ToolStatus.
var (
	ToolStatus_name = map[int32]string{
		0: "ACTIVE",
	}
	ToolStatus_value = map[string]int32{
		"ACTIVE": 0,
	}
)

func (x ToolStatus) Enum() *ToolStatus {
	p := new(ToolStatus)
	*p = x
	return p
}

func (x ToolStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ToolStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_shared_shared_proto_enumTypes[4].Descriptor()
}

func (ToolStatus) Type() protoreflect.EnumType {
	return &file_api_proto_shared_shared_proto_enumTypes[4]
}

func (x ToolStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ToolStatus.Descriptor instead.
func (ToolStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{4}
}

type Limits struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MinAmount int64 `protobuf:"varint,1,opt,name=min_amount,json=minAmount,proto3" json:"min_amount,omitempty"`
	MaxAmount int64 `protobuf:"varint,2,opt,name=max_amount,json=maxAmount,proto3" json:"max_amount,omitempty"`
}

func (x *Limits) Reset() {
	*x = Limits{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_shared_shared_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Limits) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Limits) ProtoMessage() {}

func (x *Limits) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_shared_shared_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Limits.ProtoReflect.Descriptor instead.
func (*Limits) Descriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{0}
}

func (x *Limits) GetMinAmount() int64 {
	if x != nil {
		return x.MinAmount
	}
	return 0
}

func (x *Limits) GetMaxAmount() int64 {
	if x != nil {
		return x.MaxAmount
	}
	return 0
}

type Commission struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     CommissionType    `protobuf:"varint,1,opt,name=type,proto3,enum=shared.CommissionType" json:"type,omitempty"`
	Currency *string           `protobuf:"bytes,2,opt,name=currency,proto3,oneof" json:"currency,omitempty"`
	Percent  *float64          `protobuf:"fixed64,3,opt,name=percent,proto3,oneof" json:"percent,omitempty"`
	Absolute *float64          `protobuf:"fixed64,4,opt,name=absolute,proto3,oneof" json:"absolute,omitempty"`
	Message  map[string]string `protobuf:"bytes,5,rep,name=message,proto3" json:"message,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Commission) Reset() {
	*x = Commission{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_shared_shared_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commission) ProtoMessage() {}

func (x *Commission) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_shared_shared_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commission.ProtoReflect.Descriptor instead.
func (*Commission) Descriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{1}
}

func (x *Commission) GetType() CommissionType {
	if x != nil {
		return x.Type
	}
	return CommissionType_PERCENT
}

func (x *Commission) GetCurrency() string {
	if x != nil && x.Currency != nil {
		return *x.Currency
	}
	return ""
}

func (x *Commission) GetPercent() float64 {
	if x != nil && x.Percent != nil {
		return *x.Percent
	}
	return 0
}

func (x *Commission) GetAbsolute() float64 {
	if x != nil && x.Absolute != nil {
		return *x.Absolute
	}
	return 0
}

func (x *Commission) GetMessage() map[string]string {
	if x != nil {
		return x.Message
	}
	return nil
}

type Method struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string             `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	DisplayedName  map[string]string  `protobuf:"bytes,2,rep,name=displayed_name,json=displayedName,proto3" json:"displayed_name,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ExternalSystem string             `protobuf:"bytes,3,opt,name=external_system,json=externalSystem,proto3" json:"external_system,omitempty"`
	ExternalMethod string             `protobuf:"bytes,4,opt,name=external_method,json=externalMethod,proto3" json:"external_method,omitempty"`
	Limits         map[string]*Limits `protobuf:"bytes,5,rep,name=limits,proto3" json:"limits,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Commission     *Commission        `protobuf:"bytes,6,opt,name=commission,proto3" json:"commission,omitempty"`
}

func (x *Method) Reset() {
	*x = Method{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_shared_shared_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Method) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Method) ProtoMessage() {}

func (x *Method) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_shared_shared_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Method.ProtoReflect.Descriptor instead.
func (*Method) Descriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{2}
}

func (x *Method) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Method) GetDisplayedName() map[string]string {
	if x != nil {
		return x.DisplayedName
	}
	return nil
}

func (x *Method) GetExternalSystem() string {
	if x != nil {
		return x.ExternalSystem
	}
	return ""
}

func (x *Method) GetExternalMethod() string {
	if x != nil {
		return x.ExternalMethod
	}
	return ""
}

func (x *Method) GetLimits() map[string]*Limits {
	if x != nil {
		return x.Limits
	}
	return nil
}

func (x *Method) GetCommission() *Commission {
	if x != nil {
		return x.Commission
	}
	return nil
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Group string `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
	Code  string `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_shared_shared_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_shared_shared_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{3}
}

func (x *Error) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *Error) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type ReturnURLs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Common  string  `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	Success *string `protobuf:"bytes,2,opt,name=success,proto3,oneof" json:"success,omitempty"`
	Fail    *string `protobuf:"bytes,3,opt,name=fail,proto3,oneof" json:"fail,omitempty"`
}

func (x *ReturnURLs) Reset() {
	*x = ReturnURLs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_shared_shared_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReturnURLs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReturnURLs) ProtoMessage() {}

func (x *ReturnURLs) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_shared_shared_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReturnURLs.ProtoReflect.Descriptor instead.
func (*ReturnURLs) Descriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{4}
}

func (x *ReturnURLs) GetCommon() string {
	if x != nil {
		return x.Common
	}
	return ""
}

func (x *ReturnURLs) GetSuccess() string {
	if x != nil && x.Success != nil {
		return *x.Success
	}
	return ""
}

func (x *ReturnURLs) GetFail() string {
	if x != nil && x.Fail != nil {
		return *x.Fail
	}
	return ""
}

type Tool struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId         int64            `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ExternalMethod string           `protobuf:"bytes,3,opt,name=external_method,json=externalMethod,proto3" json:"external_method,omitempty"`
	Type           *ToolType        `protobuf:"varint,4,opt,name=type,proto3,enum=shared.ToolType,oneof" json:"type,omitempty"`
	Details        *structpb.Struct `protobuf:"bytes,5,opt,name=details,proto3,oneof" json:"details,omitempty"`
	Displayed      string           `protobuf:"bytes,6,opt,name=displayed,proto3" json:"displayed,omitempty"`
	Name           string           `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	Status         ToolStatus       `protobuf:"varint,8,opt,name=status,proto3,enum=shared.ToolStatus" json:"status,omitempty"`
	Fake           bool             `protobuf:"varint,9,opt,name=fake,proto3" json:"fake,omitempty"`
	CreatedAt      int64            `protobuf:"varint,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt      int64            `protobuf:"varint,11,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Tool) Reset() {
	*x = Tool{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_shared_shared_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tool) ProtoMessage() {}

func (x *Tool) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_shared_shared_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tool.ProtoReflect.Descriptor instead.
func (*Tool) Descriptor() ([]byte, []int) {
	return file_api_proto_shared_shared_proto_rawDescGZIP(), []int{5}
}

func (x *Tool) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Tool) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Tool) GetExternalMethod() string {
	if x != nil {
		return x.ExternalMethod
	}
	return ""
}

func (x *Tool) GetType() ToolType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return ToolType_BANK_CARD
}

func (x *Tool) GetDetails() *structpb.Struct {
	if x != nil {
		return x.Details
	}
	return nil
}

func (x *Tool) GetDisplayed() string {
	if x != nil {
		return x.Displayed
	}
	return ""
}

func (x *Tool) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tool) GetStatus() ToolStatus {
	if x != nil {
		return x.Status
	}
	return ToolStatus_ACTIVE
}

func (x *Tool) GetFake() bool {
	if x != nil {
		return x.Fake
	}
	return false
}

func (x *Tool) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Tool) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

var File_api_proto_shared_shared_proto protoreflect.FileDescriptor

var file_api_proto_shared_shared_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x64, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x06, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x73, 0x12,
	0x1d, 0x0a, 0x0a, 0x6d, 0x69, 0x6e, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x6d, 0x69, 0x6e, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x6d, 0x61, 0x78, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x6d, 0x61, 0x78, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xb6, 0x02,
	0x0a, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x70, 0x65, 0x72,
	0x63, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x48, 0x01, 0x52, 0x07, 0x70, 0x65,
	0x72, 0x63, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x61, 0x62, 0x73, 0x6f,
	0x6c, 0x75, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x48, 0x02, 0x52, 0x08, 0x61, 0x62,
	0x73, 0x6f, 0x6c, 0x75, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x39, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x3a, 0x0a, 0x0c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x61, 0x62,
	0x73, 0x6f, 0x6c, 0x75, 0x74, 0x65, 0x22, 0xa9, 0x03, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x48, 0x0a, 0x0e, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x64, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0d, 0x64, 0x69,
	0x73, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x65,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x12, 0x27, 0x0a, 0x0f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x65,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x32, 0x0a,
	0x06, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x2e, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x73, 0x12, 0x32, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x40, 0x0a, 0x12, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79,
	0x65, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x49, 0x0a, 0x0b, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64,
	0x2e, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x31, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x71, 0x0a, 0x0a, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x55,
	0x52, 0x4c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x66, 0x61,
	0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x04, 0x66, 0x61, 0x69, 0x6c,
	0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x22, 0x80, 0x03, 0x0a, 0x04, 0x54, 0x6f, 0x6f,
	0x6c, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x65, 0x78,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x12, 0x29, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x10, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x54, 0x6f, 0x6f, 0x6c, 0x54,
	0x79, 0x70, 0x65, 0x48, 0x00, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12, 0x36,
	0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x48, 0x01, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x73, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x2e, 0x54, 0x6f, 0x6f, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x61, 0x6b, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x04, 0x66, 0x61, 0x6b, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x42,
	0x0a, 0x0a, 0x08, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2a, 0x28, 0x0a, 0x0d, 0x4f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07,
	0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x41, 0x59,
	0x4f, 0x55, 0x54, 0x10, 0x01, 0x2a, 0x4c, 0x0a, 0x17, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0b, 0x0a,
	0x07, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55,
	0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x03, 0x2a, 0x40, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x52, 0x43, 0x45, 0x4e, 0x54,
	0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x49, 0x58, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0c, 0x0a,
	0x08, 0x43, 0x4f, 0x4d, 0x42, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x54,
	0x45, 0x58, 0x54, 0x10, 0x03, 0x2a, 0x25, 0x0a, 0x08, 0x54, 0x6f, 0x6f, 0x6c, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0d, 0x0a, 0x09, 0x42, 0x41, 0x4e, 0x4b, 0x5f, 0x43, 0x41, 0x52, 0x44, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x10, 0x01, 0x2a, 0x18, 0x0a, 0x0a,
	0x54, 0x6f, 0x6f, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43,
	0x54, 0x49, 0x56, 0x45, 0x10, 0x00, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x6d, 0x72, 0x72, 0x77, 0x6e, 0x78, 0x74, 0x73, 0x6e, 0x2f,
	0x65, 0x63, 0x6f, 0x6d, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_shared_shared_proto_rawDescOnce sync.Once
	file_api_proto_shared_shared_proto_rawDescData = file_api_proto_shared_shared_proto_rawDesc
)

func file_api_proto_shared_shared_proto_rawDescGZIP() []byte {
	file_api_proto_shared_shared_proto_rawDescOnce.Do(func() {
		file_api_proto_shared_shared_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_shared_shared_proto_rawDescData)
	})
	return file_api_proto_shared_shared_proto_rawDescData
}

var file_api_proto_shared_shared_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_api_proto_shared_shared_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_proto_shared_shared_proto_goTypes = []interface{}{
	(OperationType)(0),           // 0: shared.OperationType
	(OperationExternalStatus)(0), // 1: shared.OperationExternalStatus
	(CommissionType)(0),          // 2: shared.CommissionType
	(ToolType)(0),                // 3: shared.ToolType
	(ToolStatus)(0),              // 4: shared.ToolStatus
	(*Limits)(nil),               // 5: shared.Limits
	(*Commission)(nil),           // 6: shared.Commission
	(*Method)(nil),               // 7: shared.Method
	(*Error)(nil),                // 8: shared.Error
	(*ReturnURLs)(nil),           // 9: shared.ReturnURLs
	(*Tool)(nil),                 // 10: shared.Tool
	nil,                          // 11: shared.Commission.MessageEntry
	nil,                          // 12: shared.Method.DisplayedNameEntry
	nil,                          // 13: shared.Method.LimitsEntry
	(*structpb.Struct)(nil),      // 14: google.protobuf.Struct
}
var file_api_proto_shared_shared_proto_depIdxs = []int32{
	2,  // 0: shared.Commission.type:type_name -> shared.CommissionType
	11, // 1: shared.Commission.message:type_name -> shared.Commission.MessageEntry
	12, // 2: shared.Method.displayed_name:type_name -> shared.Method.DisplayedNameEntry
	13, // 3: shared.Method.limits:type_name -> shared.Method.LimitsEntry
	6,  // 4: shared.Method.commission:type_name -> shared.Commission
	3,  // 5: shared.Tool.type:type_name -> shared.ToolType
	14, // 6: shared.Tool.details:type_name -> google.protobuf.Struct
	4,  // 7: shared.Tool.status:type_name -> shared.ToolStatus
	5,  // 8: shared.Method.LimitsEntry.value:type_name -> shared.Limits
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_api_proto_shared_shared_proto_init() }
func file_api_proto_shared_shared_proto_init() {
	if File_api_proto_shared_shared_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_shared_shared_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Limits); i {
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
		file_api_proto_shared_shared_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Commission); i {
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
		file_api_proto_shared_shared_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Method); i {
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
		file_api_proto_shared_shared_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_api_proto_shared_shared_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReturnURLs); i {
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
		file_api_proto_shared_shared_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tool); i {
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
	file_api_proto_shared_shared_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_api_proto_shared_shared_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_api_proto_shared_shared_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_shared_shared_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_shared_shared_proto_goTypes,
		DependencyIndexes: file_api_proto_shared_shared_proto_depIdxs,
		EnumInfos:         file_api_proto_shared_shared_proto_enumTypes,
		MessageInfos:      file_api_proto_shared_shared_proto_msgTypes,
	}.Build()
	File_api_proto_shared_shared_proto = out.File
	file_api_proto_shared_shared_proto_rawDesc = nil
	file_api_proto_shared_shared_proto_goTypes = nil
	file_api_proto_shared_shared_proto_depIdxs = nil
}
