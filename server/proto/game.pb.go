// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: game.proto

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

type ClientActionType int32

const (
	ClientActionType_ACTION_PING       ClientActionType = 0
	ClientActionType_ACTION_CONNECT    ClientActionType = 1
	ClientActionType_ACTION_DISCONNECT ClientActionType = 2
	ClientActionType_ACTION_MOVE       ClientActionType = 3
	ClientActionType_ACTION_ATTACK     ClientActionType = 4
	ClientActionType_ACTION_DAMAGE     ClientActionType = 5
	ClientActionType_ACTION_DEATH      ClientActionType = 6
	ClientActionType_ACTION_RESPAWN    ClientActionType = 7
)

// Enum value maps for ClientActionType.
var (
	ClientActionType_name = map[int32]string{
		0: "ACTION_PING",
		1: "ACTION_CONNECT",
		2: "ACTION_DISCONNECT",
		3: "ACTION_MOVE",
		4: "ACTION_ATTACK",
		5: "ACTION_DAMAGE",
		6: "ACTION_DEATH",
		7: "ACTION_RESPAWN",
	}
	ClientActionType_value = map[string]int32{
		"ACTION_PING":       0,
		"ACTION_CONNECT":    1,
		"ACTION_DISCONNECT": 2,
		"ACTION_MOVE":       3,
		"ACTION_ATTACK":     4,
		"ACTION_DAMAGE":     5,
		"ACTION_DEATH":      6,
		"ACTION_RESPAWN":    7,
	}
)

func (x ClientActionType) Enum() *ClientActionType {
	p := new(ClientActionType)
	*p = x
	return p
}

func (x ClientActionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClientActionType) Descriptor() protoreflect.EnumDescriptor {
	return file_game_proto_enumTypes[0].Descriptor()
}

func (ClientActionType) Type() protoreflect.EnumType {
	return &file_game_proto_enumTypes[0]
}

func (x ClientActionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ClientActionType.Descriptor instead.
func (ClientActionType) EnumDescriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{0}
}

type ServerMessageType int32

const (
	ServerMessageType_MESSAGE_PING ServerMessageType = 0
	ServerMessageType_MESSAGE_JOIN ServerMessageType = 1
	ServerMessageType_MESSAGE_TICK ServerMessageType = 2
)

// Enum value maps for ServerMessageType.
var (
	ServerMessageType_name = map[int32]string{
		0: "MESSAGE_PING",
		1: "MESSAGE_JOIN",
		2: "MESSAGE_TICK",
	}
	ServerMessageType_value = map[string]int32{
		"MESSAGE_PING": 0,
		"MESSAGE_JOIN": 1,
		"MESSAGE_TICK": 2,
	}
)

func (x ServerMessageType) Enum() *ServerMessageType {
	p := new(ServerMessageType)
	*p = x
	return p
}

func (x ServerMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServerMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_game_proto_enumTypes[1].Descriptor()
}

func (ServerMessageType) Type() protoreflect.EnumType {
	return &file_game_proto_enumTypes[1]
}

func (x ServerMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ServerMessageType.Descriptor instead.
func (ServerMessageType) EnumDescriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{1}
}

type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X uint32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	Z uint32 `protobuf:"varint,2,opt,name=z,proto3" json:"z,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{0}
}

func (x *Location) GetX() uint32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Location) GetZ() uint32 {
	if x != nil {
		return x.Z
	}
	return 0
}

type ClientAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    int32  `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Payload string `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ClientAction) Reset() {
	*x = ClientAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientAction) ProtoMessage() {}

func (x *ClientAction) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientAction.ProtoReflect.Descriptor instead.
func (*ClientAction) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{1}
}

func (x *ClientAction) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *ClientAction) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

type Connect struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId uint64    `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"`
	Color    string    `protobuf:"bytes,2,opt,name=color,proto3" json:"color,omitempty"`
	Spawn    *Location `protobuf:"bytes,3,opt,name=spawn,proto3" json:"spawn,omitempty"`
}

func (x *Connect) Reset() {
	*x = Connect{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Connect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Connect) ProtoMessage() {}

func (x *Connect) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Connect.ProtoReflect.Descriptor instead.
func (*Connect) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{2}
}

func (x *Connect) GetPlayerId() uint64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *Connect) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *Connect) GetSpawn() *Location {
	if x != nil {
		return x.Spawn
	}
	return nil
}

type Disconnect struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId uint64 `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"`
}

func (x *Disconnect) Reset() {
	*x = Disconnect{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Disconnect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Disconnect) ProtoMessage() {}

func (x *Disconnect) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Disconnect.ProtoReflect.Descriptor instead.
func (*Disconnect) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{3}
}

func (x *Disconnect) GetPlayerId() uint64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

type Move struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId uint64    `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"`
	Target   *Location `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
}

func (x *Move) Reset() {
	*x = Move{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Move) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Move) ProtoMessage() {}

func (x *Move) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Move.ProtoReflect.Descriptor instead.
func (*Move) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{4}
}

func (x *Move) GetPlayerId() uint64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *Move) GetTarget() *Location {
	if x != nil {
		return x.Target
	}
	return nil
}

type Attack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourcePlayerId       uint64    `protobuf:"varint,1,opt,name=sourcePlayerId,proto3" json:"sourcePlayerId,omitempty"`
	SourcePlayerLocation *Location `protobuf:"bytes,2,opt,name=sourcePlayerLocation,proto3" json:"sourcePlayerLocation,omitempty"`
	TargetPlayerId       uint64    `protobuf:"varint,3,opt,name=targetPlayerId,proto3" json:"targetPlayerId,omitempty"`
	TargetPlayerLocation *Location `protobuf:"bytes,4,opt,name=targetPlayerLocation,proto3" json:"targetPlayerLocation,omitempty"`
}

func (x *Attack) Reset() {
	*x = Attack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attack) ProtoMessage() {}

func (x *Attack) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attack.ProtoReflect.Descriptor instead.
func (*Attack) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{5}
}

func (x *Attack) GetSourcePlayerId() uint64 {
	if x != nil {
		return x.SourcePlayerId
	}
	return 0
}

func (x *Attack) GetSourcePlayerLocation() *Location {
	if x != nil {
		return x.SourcePlayerLocation
	}
	return nil
}

func (x *Attack) GetTargetPlayerId() uint64 {
	if x != nil {
		return x.TargetPlayerId
	}
	return 0
}

func (x *Attack) GetTargetPlayerLocation() *Location {
	if x != nil {
		return x.TargetPlayerLocation
	}
	return nil
}

type Damage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourcePlayerId uint64 `protobuf:"varint,1,opt,name=sourcePlayerId,proto3" json:"sourcePlayerId,omitempty"`
	TargetPlayerId uint64 `protobuf:"varint,2,opt,name=targetPlayerId,proto3" json:"targetPlayerId,omitempty"`
	DamageDealt    int32  `protobuf:"varint,3,opt,name=damageDealt,proto3" json:"damageDealt,omitempty"`
}

func (x *Damage) Reset() {
	*x = Damage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Damage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Damage) ProtoMessage() {}

func (x *Damage) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Damage.ProtoReflect.Descriptor instead.
func (*Damage) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{6}
}

func (x *Damage) GetSourcePlayerId() uint64 {
	if x != nil {
		return x.SourcePlayerId
	}
	return 0
}

func (x *Damage) GetTargetPlayerId() uint64 {
	if x != nil {
		return x.TargetPlayerId
	}
	return 0
}

func (x *Damage) GetDamageDealt() int32 {
	if x != nil {
		return x.DamageDealt
	}
	return 0
}

type Death struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId uint64    `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"`
	Location *Location `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *Death) Reset() {
	*x = Death{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Death) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Death) ProtoMessage() {}

func (x *Death) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Death.ProtoReflect.Descriptor instead.
func (*Death) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{7}
}

func (x *Death) GetPlayerId() uint64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *Death) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

type Respawn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId uint64    `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"`
	Spawn    *Location `protobuf:"bytes,2,opt,name=spawn,proto3" json:"spawn,omitempty"`
}

func (x *Respawn) Reset() {
	*x = Respawn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Respawn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Respawn) ProtoMessage() {}

func (x *Respawn) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Respawn.ProtoReflect.Descriptor instead.
func (*Respawn) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{8}
}

func (x *Respawn) GetPlayerId() uint64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *Respawn) GetSpawn() *Location {
	if x != nil {
		return x.Spawn
	}
	return nil
}

type ServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    int32  `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Payload string `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ServerMessage) Reset() {
	*x = ServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage) ProtoMessage() {}

func (x *ServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage.ProtoReflect.Descriptor instead.
func (*ServerMessage) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{9}
}

func (x *ServerMessage) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *ServerMessage) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

type GameTickAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  uint32 `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GameTickAction) Reset() {
	*x = GameTickAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameTickAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameTickAction) ProtoMessage() {}

func (x *GameTickAction) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameTickAction.ProtoReflect.Descriptor instead.
func (*GameTickAction) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{10}
}

func (x *GameTickAction) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *GameTickAction) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type GameTick struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tick    uint32            `protobuf:"varint,1,opt,name=tick,proto3" json:"tick,omitempty"`
	Actions []*GameTickAction `protobuf:"bytes,2,rep,name=actions,proto3" json:"actions,omitempty"`
}

func (x *GameTick) Reset() {
	*x = GameTick{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameTick) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameTick) ProtoMessage() {}

func (x *GameTick) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameTick.ProtoReflect.Descriptor instead.
func (*GameTick) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{11}
}

func (x *GameTick) GetTick() uint32 {
	if x != nil {
		return x.Tick
	}
	return 0
}

func (x *GameTick) GetActions() []*GameTickAction {
	if x != nil {
		return x.Actions
	}
	return nil
}

type JoinGameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId uint64     `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"`
	Color    string     `protobuf:"bytes,2,opt,name=color,proto3" json:"color,omitempty"`
	Spawn    *Location  `protobuf:"bytes,3,opt,name=spawn,proto3" json:"spawn,omitempty"`
	Others   []*Connect `protobuf:"bytes,4,rep,name=others,proto3" json:"others,omitempty"`
}

func (x *JoinGameResponse) Reset() {
	*x = JoinGameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinGameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinGameResponse) ProtoMessage() {}

func (x *JoinGameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinGameResponse.ProtoReflect.Descriptor instead.
func (*JoinGameResponse) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{12}
}

func (x *JoinGameResponse) GetPlayerId() uint64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *JoinGameResponse) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *JoinGameResponse) GetSpawn() *Location {
	if x != nil {
		return x.Spawn
	}
	return nil
}

func (x *JoinGameResponse) GetOthers() []*Connect {
	if x != nil {
		return x.Others
	}
	return nil
}

var File_game_proto protoreflect.FileDescriptor

var file_game_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x61,
	0x6d, 0x65, 0x22, 0x26, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0c,
	0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01,
	0x7a, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x01, 0x7a, 0x22, 0x3c, 0x0a, 0x0c, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x61, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x24, 0x0a, 0x05, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x22, 0x28, 0x0a, 0x0a, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4a, 0x0a, 0x04, 0x4d, 0x6f, 0x76, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x06, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x67, 0x61, 0x6d, 0x65,
	0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x22, 0xe0, 0x01, 0x0a, 0x06, 0x41, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x12, 0x26, 0x0a, 0x0e,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x42, 0x0a, 0x14, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x14, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x0e, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x42, 0x0a, 0x14, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x14,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x7a, 0x0a, 0x06, 0x44, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x26,
	0x0a, 0x0e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x44, 0x65, 0x61, 0x6c, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0b, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x44, 0x65, 0x61, 0x6c, 0x74,
	0x22, 0x4f, 0x0a, 0x05, 0x44, 0x65, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x4b, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x05, 0x73, 0x70, 0x61, 0x77,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x22, 0x3d,
	0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x3a, 0x0a,
	0x0e, 0x47, 0x61, 0x6d, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x4e, 0x0a, 0x08, 0x47, 0x61, 0x6d,
	0x65, 0x54, 0x69, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x63, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x69, 0x63, 0x6b, 0x12, 0x2e, 0x0a, 0x07, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x91, 0x01, 0x0a, 0x10, 0x4a, 0x6f,
	0x69, 0x6e, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x6c, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72,
	0x12, 0x24, 0x0a, 0x05, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x05, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x12, 0x25, 0x0a, 0x06, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x06, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x73, 0x2a, 0xab, 0x01,
	0x0a, 0x10, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x50, 0x49, 0x4e,
	0x47, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x4f,
	0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x43, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x44, 0x49, 0x53, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x02, 0x12, 0x0f,
	0x0a, 0x0b, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x4f, 0x56, 0x45, 0x10, 0x03, 0x12,
	0x11, 0x0a, 0x0d, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x54, 0x54, 0x41, 0x43, 0x4b,
	0x10, 0x04, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x41, 0x4d,
	0x41, 0x47, 0x45, 0x10, 0x05, 0x12, 0x10, 0x0a, 0x0c, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x44, 0x45, 0x41, 0x54, 0x48, 0x10, 0x06, 0x12, 0x12, 0x0a, 0x0e, 0x41, 0x43, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x41, 0x57, 0x4e, 0x10, 0x07, 0x2a, 0x49, 0x0a, 0x11, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x50, 0x49, 0x4e, 0x47,
	0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x4a, 0x4f,
	0x49, 0x4e, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f,
	0x54, 0x49, 0x43, 0x4b, 0x10, 0x02, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_game_proto_rawDescOnce sync.Once
	file_game_proto_rawDescData = file_game_proto_rawDesc
)

func file_game_proto_rawDescGZIP() []byte {
	file_game_proto_rawDescOnce.Do(func() {
		file_game_proto_rawDescData = protoimpl.X.CompressGZIP(file_game_proto_rawDescData)
	})
	return file_game_proto_rawDescData
}

var file_game_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_game_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_game_proto_goTypes = []interface{}{
	(ClientActionType)(0),    // 0: game.ClientActionType
	(ServerMessageType)(0),   // 1: game.ServerMessageType
	(*Location)(nil),         // 2: game.Location
	(*ClientAction)(nil),     // 3: game.ClientAction
	(*Connect)(nil),          // 4: game.Connect
	(*Disconnect)(nil),       // 5: game.Disconnect
	(*Move)(nil),             // 6: game.Move
	(*Attack)(nil),           // 7: game.Attack
	(*Damage)(nil),           // 8: game.Damage
	(*Death)(nil),            // 9: game.Death
	(*Respawn)(nil),          // 10: game.Respawn
	(*ServerMessage)(nil),    // 11: game.ServerMessage
	(*GameTickAction)(nil),   // 12: game.GameTickAction
	(*GameTick)(nil),         // 13: game.GameTick
	(*JoinGameResponse)(nil), // 14: game.JoinGameResponse
}
var file_game_proto_depIdxs = []int32{
	2,  // 0: game.Connect.spawn:type_name -> game.Location
	2,  // 1: game.Move.target:type_name -> game.Location
	2,  // 2: game.Attack.sourcePlayerLocation:type_name -> game.Location
	2,  // 3: game.Attack.targetPlayerLocation:type_name -> game.Location
	2,  // 4: game.Death.location:type_name -> game.Location
	2,  // 5: game.Respawn.spawn:type_name -> game.Location
	12, // 6: game.GameTick.actions:type_name -> game.GameTickAction
	2,  // 7: game.JoinGameResponse.spawn:type_name -> game.Location
	4,  // 8: game.JoinGameResponse.others:type_name -> game.Connect
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_game_proto_init() }
func file_game_proto_init() {
	if File_game_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_game_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
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
		file_game_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientAction); i {
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
		file_game_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Connect); i {
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
		file_game_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Disconnect); i {
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
		file_game_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Move); i {
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
		file_game_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attack); i {
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
		file_game_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Damage); i {
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
		file_game_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Death); i {
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
		file_game_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Respawn); i {
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
		file_game_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage); i {
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
		file_game_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameTickAction); i {
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
		file_game_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameTick); i {
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
		file_game_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinGameResponse); i {
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
			RawDescriptor: file_game_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_game_proto_goTypes,
		DependencyIndexes: file_game_proto_depIdxs,
		EnumInfos:         file_game_proto_enumTypes,
		MessageInfos:      file_game_proto_msgTypes,
	}.Build()
	File_game_proto = out.File
	file_game_proto_rawDesc = nil
	file_game_proto_goTypes = nil
	file_game_proto_depIdxs = nil
}
