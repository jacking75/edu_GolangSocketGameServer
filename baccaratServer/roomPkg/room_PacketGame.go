package roomPkg


import (
	"go.uber.org/zap"

	. "golang_socketGameServer_codelab/gohipernetFake"

	"golang_socketGameServer_codelab/baccaratServer/connectedSessions"
	"golang_socketGameServer_codelab/baccaratServer/protocol"
)


func (room *baseRoom) _packetProcess_GameStart(user *roomUser, packet protocol.Packet) int16 {
	NTELIB_LOG_DEBUG("[[Room _packetProcess_GameStart ]")

	errorCode := (int16)(protocol.ERROR_CODE_NONE)
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId

	// 방의 상태가 NONE 인가?
	if room.IsStateNone() == false {
		errorCode = protocol.ERROR_CODE_ROOM_GAME_START_INVALID_ROOM_STATE
		goto CheckError
	}

	// 유저의 최소 수가 2명 이상인가?
	if room.getCurUserCount() < 2 {
		errorCode = protocol.ERROR_CODE_ROOM_GAME_START_NOT_ENOUGH_MEMBERS
		goto CheckError
	}

	// 시작 요청은 방장이 하는가?
	if room.getMasterSessionUniqueId() != sessionUniqueId {
		errorCode = protocol.ERROR_CODE_ROOM_GAME_START_NOT_MASTER
		goto CheckError
	}

	// 게임을 시작한다.
	// 방 상태를 배팅 대기로 바꾼다
	room.changeState(ROOM_STATE_GAME_WAIT_BATTING)

	_sendRoomGameStartResult(sessionIndex, sessionUniqueId, errorCode)
	_sendRoomGameStartNotify();
	return errorCode

CheckError:
	_sendRoomGameStartResult(sessionIndex, sessionUniqueId, errorCode)
	return errorCode
}

func _sendRoomGameStartResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	//TODO 올바르게 고쳐야 한다
	response := protocol.RoomLeaveResPacket{ result }
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}

func _sendRoomGameStartNotify() {
	//TODO 구현하기
}



func (room *baseRoom) _packetProcess_GameBatting(user *roomUser, packet protocol.Packet) int16 {
	NTELIB_LOG_DEBUG("[[Room _packetProcess_GameBatting ]")

	errorCode := (int16)(protocol.ERROR_CODE_NONE)
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId


	// 방의 상태가 배팅 기다림인가?
	if room.IsStateGameBattingWait() == false {
		errorCode = protocol.ERROR_CODE_ROOM_GAME_BATTING_INVALID_ROOM_STATE
		goto CheckError
	}

	// 배팅 타입이 유효한가 ?

	// 배팅하고, 모두에게 알린다.

	// 배팅을 다 했다면 카드를 날린다. 여기에 결과까지 다 들어간다.
	//   이 함수는 대기 타임 완료에 따라 자동으로 호출될 수 있으므로 함수로 뺀다
	return protocol.ERROR_CODE_NONE

CheckError:
	_sendRoomGameStartResult(sessionIndex, sessionUniqueId, errorCode)
	return errorCode
}





func (room *baseRoom) _packetProcess_LeaveUser11(user *roomUser, packet protocol.Packet) int16 {
	NTELIB_LOG_DEBUG("[[Room _packetProcess_LeaveUser ]")

	room._leaveUserProcess(user)

	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId
	_sendRoomLeaveResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_NONE)
	return protocol.ERROR_CODE_NONE
}

func (room *baseRoom) _leaveUserProcess1(user *roomUser) {
	NTELIB_LOG_DEBUG("[[Room _leaveUserProcess]")

	//TODO 방의 상태가 ROOM_STATE_NOE, ROOM_STATE_GAME_RESULT 일 때만 나갈 수 있다.

	//TODO 유저가 접속이 끊어져서 나가는 경우라면 게임이 끝날 때까지 유저 정보 들고 있다가
	//  ROOM_STATE_GAME_RESULT 상태일 때 제거한다.

	roomUserUniqueId := user.RoomUniqueId
	userSessionIndex := user.netSessionIndex
	userSessionUniqueId := user.netSessionUniqueId

	room._removeUser(user)

	room._sendRoomLeaveUserNotify(roomUserUniqueId, userSessionUniqueId)

	curTime := NetLib_GetCurrnetUnixTime()
	connectedSessions.SetRoomNumber(userSessionIndex, userSessionUniqueId, -1, curTime)
}


func _sendRoomLeaveResult111(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomLeaveResPacket{ result }
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}

func (room *baseRoom) _sendRoomLeaveUserNotify111(roomUserUniqueId uint64, userSessionUniqueId uint64) {
	NTELIB_LOG_DEBUG("Room _sendRoomLeaveUserNotify", zap.Uint64("userSessionUniqueId", userSessionUniqueId), zap.Int32("RoomIndex", room._index))

	notifyPacket := protocol.RoomLeaveUserNtfPacket{roomUserUniqueId	}
	sendBuf, packetSize := notifyPacket.EncodingPacket()
	room.broadcastPacket(int16(packetSize), sendBuf, userSessionUniqueId)
}
