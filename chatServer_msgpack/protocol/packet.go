package protocol

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack/v4"
	"reflect"

	. "gohipernetFake"
)

const (
	PACKET_TYPE_NORMAL   = 0
	PACKET_TYPE_COMPRESS = 1
	PACKET_TYPE_SECURE   = 2
)

const (
	MAX_USER_ID_BYTE_LENGTH      = 16
	MAX_USER_PW_BYTE_LENGTH      = 16
	MAX_CHAT_MESSAGE_BYTE_LENGTH = 126
)

type Header struct {
	TotalSize  int16
	ID         int16
	PacketType int8 // 비트 필드로 데이터 설정. 0 이면 Normal, 1번 비트 On(압축), 2번 비트 On(암호화)
}

type Packet struct {
	UserSessionIndex    int32
	UserSessionUniqueId uint64
	Id                  uint16
	DataSize            uint16
	Data                []byte
}

func (packet Packet) GetSessionInfo() (int32, uint64) {
	return packet.UserSessionIndex, packet.UserSessionUniqueId
}

var _clientSessionHeaderSize uint16
var _ServerSessionHeaderSize uint16

func Init_packet() {
	_clientSessionHeaderSize = protocolInitHeaderSize()
	_ServerSessionHeaderSize = protocolInitHeaderSize()
}

func ClientHeaderSize() uint16 {
	return _clientSessionHeaderSize
}
func ServerHeaderSize() uint16 {
	return _ServerSessionHeaderSize
}

func protocolInitHeaderSize() uint16 {
	var packetHeader Header
	headerSize := Sizeof(reflect.TypeOf(packetHeader))
	return (uint16)(headerSize)
}

// Header의 PacketID만 읽는다
func PeekPacketID(rawData []byte) uint16 {
	packetID := binary.LittleEndian.Uint16(rawData[2:])
	return uint16(packetID)
}

// 보디데이터의 참조만 가져간다
func PeekPacketBody(rawData []byte) (bodySize uint16, refBody []byte) {
	headerSize := ClientHeaderSize()
	totalSize := binary.LittleEndian.Uint16(rawData)
	bodySize = totalSize - headerSize

	if bodySize > 0 {
		refBody = rawData[headerSize:]
	}

	return bodySize, refBody
}

func DecodingPacketHeader(header *Header, data []byte) {
	reader := MakeReader(data, true)
	header.TotalSize, _ = reader.ReadS16()
	header.ID, _ = reader.ReadS16()
	header.PacketType, _ = reader.ReadS8()
}

func EncodingPacketHeader(writer *RawPacketData, totalSize uint16, pktId uint16, packetType int8) {
	writer.WriteU16(totalSize)
	writer.WriteU16(pktId)
	writer.WriteS8(packetType)
}
func EncodingPacketHeaderInfo(bodyData []byte, pktId uint16, packetType int8) []byte {
	totalSize := ClientHeaderSize() + uint16(len(bodyData))
	sendBuf := make([]byte, totalSize)

	writer := MakeWriter(sendBuf, true)
	writer.WriteU16(totalSize)
	writer.WriteU16(pktId)
	writer.WriteS8(packetType)
	writer.WriteBytes(bodyData)
	return sendBuf
}
///<<< 패킷 인코딩/디코딩
// C# 클라이언트가 msgpack 라이브러리로 SimpleMsgPack.Net을 사용하고 있는데 정수는 int64만 지원하므로 정수는 int64(혹은 uint64) 타입만 지원한다

// [[[ 로그인 ]]] PACKET_ID_LOGIN_REQ
type LoginReqPacket struct {
	UserID string
	UserPW string
}

type LoginResPacket struct {
	Result int64
}


// [[[  ]]]   PACKET_ID_ERROR_NTF
type ErrorNtfPacket struct {
	ErrorCode int64
}


/// [ 방 입장 ]
type RoomEnterReqPacket struct {
	RoomNumber int64
}

type RoomEnterResPacket struct {
	Result           int64
	RoomUserUniqueId uint64
}

///  방에 있는 있는 유저 리스트 통보(채널에 들어오는 유저에게)
type RoomUserInfoPktData struct {
	UniqueId uint64
	ID       string
}

type RoomUserListNtfPacket struct {
	UserCount int64
	UniqueId []uint64
	ID       []string
}


// 채널에 있는 유저들에게 새 유저의 정보를 알려준다
type RoomNewUserNtfPacket struct {
	UniqueId uint64
	ID       string
}


//<<< 방에서 나가기
type RoomLeaveResPacket struct {
	Result int64
}

type RoomLeaveUserNtfPacket struct {
	UserUniqueId uint64
}


/// [ 방 채팅 ]]
type RoomChatReqPacket struct {
	Msg      string
}

type RoomChatResPacket struct {
	Result int64
}

type RoomChatNtfPacket struct {
	UserUniqueId uint64
	Msg          string
}


///<<< Room Relay
type RoomRelayReqPacket struct {
	Data []byte
}

func (request RoomRelayReqPacket) EncodingPacket(size int16) ([]byte, uint16) {
	totalSize := _clientSessionHeaderSize + uint16(len(request.Data))
	sendBuf := make([]byte, totalSize)
	writer := MakeWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_RELAY_REQ, 0)

	writer.WriteBytes(request.Data)
	return sendBuf, totalSize
}

func (request *RoomRelayReqPacket) Decoding(bodyData []byte) bool {
	reader := MakeReader(bodyData, true)
	request.Data = reader.ReadBytes(len(bodyData))
	return true
}

type RoomRelayNtfPacket struct {
	RoomUserUniqueId uint64
	Data             []byte
}

func (notify RoomRelayNtfPacket) EncodingPacket(size int16) ([]byte, uint16) {
	totalSize := _clientSessionHeaderSize + 8 + uint16(len(notify.Data))
	sendBuf := make([]byte, totalSize)
	writer := MakeWriter(sendBuf, true)
	EncodingPacketHeader(&writer, totalSize, PACKET_ID_ROOM_RELAY_NTF, 0)

	writer.WriteU64(notify.RoomUserUniqueId)
	writer.WriteBytes(notify.Data)
	return sendBuf, totalSize
}

func (notify *RoomRelayNtfPacket) Decoding(bodyData []byte) bool {
	reader := MakeReader(bodyData, true)
	notify.RoomUserUniqueId, _ = reader.ReadU64()
	notify.Data = reader.ReadBytes(len(bodyData) - 8)
	return true
}

func NotifyErrorPacket(sessionIndex int32, sessionUniqueId uint64, errorCode int16) {
	var response ErrorNtfPacket
	response.ErrorCode = int64(errorCode)

	bodyData, err := msgpack.Marshal(response)
	if err != nil {
		return
	}

	sendPacket := EncodingPacketHeaderInfo(bodyData, uint16(PACKET_ID_ERROR_NTF), 0)
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}
