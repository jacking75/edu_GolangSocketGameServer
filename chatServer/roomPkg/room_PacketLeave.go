package roomPkg

import (
	"go.uber.org/zap"

	. "golang_socketGameServer_codelab/gohipernetFake"

	"golang_socketGameServer_codelab/chatServer/connectedSessions"
	"golang_socketGameServer_codelab/chatServer/protocol"
)


func (room *baseRoom) _packetProcess_LeaveUser(user *roomUser, packet protocol.Packet) int16 {
	NTELIB_LOG_DEBUG("[[Room _packetProcess_LeaveUser ]")

	room._leaveUserProcess(user)

	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId
	_sendRoomLeaveResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_NONE)
	return protocol.ERROR_CODE_NONE
}

func (room *baseRoom) _leaveUserProcess(user *roomUser) {
	NTELIB_LOG_DEBUG("[[Room _leaveUserProcess]")

	roomUserUniqueId := user.RoomUniqueId
	userSessionIndex := user.netSessionIndex
	userSessionUniqueId := user.netSessionUniqueId

	room._removeUser(user)

	room._sendRoomLeaveUserNotify(roomUserUniqueId, userSessionUniqueId)

	curTime := NetLib_GetCurrnetUnixTime()
	connectedSessions.SetRoomNumber(userSessionIndex, userSessionUniqueId, -1, curTime)
}


func _sendRoomLeaveResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomLeaveResPacket{ result }
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}

func (room *baseRoom) _sendRoomLeaveUserNotify(roomUserUniqueId uint64, userSessionUniqueId uint64) {
	NTELIB_LOG_DEBUG("Room _sendRoomLeaveUserNotify", zap.Uint64("userSessionUniqueId", userSessionUniqueId), zap.Int32("RoomIndex", room._index))

	notifyPacket := protocol.RoomLeaveUserNtfPacket{roomUserUniqueId	}
	sendBuf, packetSize := notifyPacket.EncodingPacket()
	room.broadcastPacket(int16(packetSize), sendBuf, userSessionUniqueId)
}