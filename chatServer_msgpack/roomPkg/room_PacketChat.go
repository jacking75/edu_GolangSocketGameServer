package roomPkg

import (
	"github.com/vmihailenco/msgpack/v4"
	. "gohipernetFake"

	"main/protocol"
)

func (room *baseRoom) _packetProcess_Chat(user *roomUser, packet protocol.Packet) int16 {
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId

	var chatPacket protocol.RoomChatReqPacket
	err := msgpack.Unmarshal(packet.Data, &chatPacket)
	if err != nil {
		_sendRoomChatResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_PACKET_DECODING_FAIL)
		return protocol.ERROR_CODE_PACKET_DECODING_FAIL
	}

	// 채팅 최대길이 제한. MAX_CHAT_MESSAGE_BYTE_LENGTH는 정의만 되어 있고 검증이 빠져 있어서
	// 클라이언트가 임의 길이의 메시지를 그대로 방 전체에 브로드캐스트할 수 있었다.
	msgLen := len(chatPacket.Msg)
	if msgLen < 1 || msgLen > protocol.MAX_CHAT_MESSAGE_BYTE_LENGTH {
		_sendRoomChatResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN)
		return protocol.ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN
	}

	var chatNotifyResponse protocol.RoomChatNtfPacket
	chatNotifyResponse.UserUniqueId = user.RoomUniqueId
	chatNotifyResponse.Msg = chatPacket.Msg
	bodyData, err := msgpack.Marshal(chatNotifyResponse)
	if err != nil {
		return protocol.ERROR_CODE_PACKET_ENCODING_FAIL
	}

	notifySendBuf := protocol.EncodingPacketHeaderInfo(bodyData, uint16(protocol.PACKET_ID_ROOM_CHAT_NOTIFY), 0)
	packetSize := uint16(len(notifySendBuf))
	room.broadcastPacket(packetSize, notifySendBuf, 0)

	_sendRoomChatResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_NONE)

	return protocol.ERROR_CODE_NONE
}

func _sendRoomChatResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomChatResPacket{int64(result)}
	bodyData, err := msgpack.Marshal(response)
	if err != nil {
		panic(err)
	}

	sendPacket := protocol.EncodingPacketHeaderInfo(bodyData, uint16(protocol.PACKET_ID_ROOM_CHAT_RES), 0)

	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}
