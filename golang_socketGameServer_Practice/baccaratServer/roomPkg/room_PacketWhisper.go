package roomPkg

import (
	"go.uber.org/zap"

	. "golang_socketGameServer_codelab/gohipernetFake"

	"golang_socketGameServer_codelab/baccaratServer/protocol"
)


func (room *baseRoom) _packetProcess_Whisper(user *roomUser, packet protocol.Packet) int16 {
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId

	var whisperPacket protocol.RoomWhisperReqPacket
	if whisperPacket.Decoding(packet.Data) == false {
		_sendRoomChatResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_PACKET_DECODING_FAIL)
		return protocol.ERROR_CODE_PACKET_DECODING_FAIL
	}

	// 채팅 최대길이 제한
	msgLen := len(whisperPacket.Msg)
	if msgLen < 1 || msgLen > protocol.MAX_CHAT_MESSAGE_BYTE_LENGTH {
		_sendRoomWhisperResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN)
		return protocol.ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN
	}


	var whisperNotifyResponse protocol.RoomWhisperNtfPacket
	whisperNotifyResponse.SendUserUniqueId = user.RoomUniqueId
	whisperNotifyResponse.ReceiveUserUniqueId = whisperPacket.ReceiveUserUniqueId
	whisperNotifyResponse.MsgLen = int16(msgLen)
	whisperNotifyResponse.Msg = whisperPacket.Msg
	notifySendBuf, packetSize := whisperNotifyResponse.EncodingPacket()
	errcode := room.whisperPacket(whisperNotifyResponse.SendUserUniqueId, whisperPacket.ReceiveUserUniqueId, packetSize, notifySendBuf)

	if errcode != protocol.ERROR_CODE_NONE {
		_sendRoomWhisperResult(sessionIndex, sessionUniqueId, errcode)
		return errcode
	}
	_sendRoomWhisperResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_NONE)

	receiverID := room._userSessionUniqueIdMap[whisperPacket.ReceiveUserUniqueId].ID[:]

	NTELIB_LOG_DEBUG("ParkChannel Whisper Notify Function", zap.String("Sender", string(user.ID[:])),
		zap.String("Receiver", string(receiverID)), zap.String("Message", string(whisperPacket.Msg)))

	return protocol.ERROR_CODE_NONE
}

func _sendRoomWhisperResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomWhisperResPacket{ Result:result }
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}