package roomPkg

import (
	"time"

	. "gohipernetFake"

	"main/connectedSessions"
	"main/protocol"
)


func (room *baseRoom) _packetProcess_LeaveUser(user *roomUser, packet protocol.Packet) int16 {
	room._leaveUserProcess(user)

	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId
	_sendRoomLeaveResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_NONE)
	return protocol.ERROR_CODE_NONE
}

func (room *baseRoom) _leaveUserProcess(user *roomUser) {
	//TODO 방의 상태가 ROOM_STATE_NOE, ROOM_STATE_GAME_RESULT 일 때만 나갈 수 있다.

	//TODO 유저가 접속이 끊어져서 나가는 경우라면 게임이 끝날 때까지 유저 정보 들고 있다가
	//  ROOM_STATE_GAME_RESULT 상태일 때 제거한다.

	roomUserUniqueId := user.RoomUniqueId
	userSessionIndex := user.netSessionIndex
	userSessionUniqueId := user.netSessionUniqueId

	room._removeUser(user)

	room._sendRoomLeaveUserNotify(roomUserUniqueId, userSessionUniqueId)

	curTime := time.Now().Unix()
	connectedSessions.SetRoomNumber(userSessionIndex, userSessionUniqueId, -1, curTime)
}


func _sendRoomLeaveResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomLeaveResPacket{ result }
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}

func (room *baseRoom) _sendRoomLeaveUserNotify(roomUserUniqueId uint64, userSessionUniqueId uint64) {
	notifyPacket := protocol.RoomLeaveUserNtfPacket{roomUserUniqueId	}
	sendBuf, packetSize := notifyPacket.EncodingPacket()
	room.broadcastPacket(int16(packetSize), sendBuf, userSessionUniqueId)
}