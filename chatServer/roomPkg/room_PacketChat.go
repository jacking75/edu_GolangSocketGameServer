package roomPkg

import (
	. "gohipernetFake"

	"main/protocol"
)

func (room *baseRoom) _packetProcess_Chat(user *roomUser, packet protocol.Packet) int16 {
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId

	var chatPacket protocol.RoomChatReqPacket
	if chatPacket.Decoding(packet.Data) == false {
		_sendRoomChatResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_PACKET_DECODING_FAIL)
		return protocol.ERROR_CODE_PACKET_DECODING_FAIL
	}

	// 채팅 최대길이 제한
	msgLen := len(chatPacket.Msgs)
	if msgLen < 1 || msgLen > protocol.MAX_CHAT_MESSAGE_BYTE_LENGTH {
		_sendRoomChatResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN)
		return protocol.ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN
	}

	var chatNotifyResponse protocol.RoomChatNtfPacket
	chatNotifyResponse.RoomUserUniqueId = user.RoomUniqueId
	chatNotifyResponse.MsgLen = int16(msgLen)
	chatNotifyResponse.Msg = chatPacket.Msgs
	notifySendBuf, packetSize := chatNotifyResponse.EncodingPacket()
	room.broadcastPacket(packetSize, notifySendBuf, 0)

	_sendRoomChatResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_NONE)

	return protocol.ERROR_CODE_NONE
}

func _sendRoomChatResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomChatResPacket{result}
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}
