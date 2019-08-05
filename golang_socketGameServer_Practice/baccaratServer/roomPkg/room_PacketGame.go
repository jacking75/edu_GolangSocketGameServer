package roomPkg

import (
	"bytes"
	. "golang_socketGameServer_codelab/gohipernetFake"

	"golang_socketGameServer_codelab/baccaratServer/protocol"
)


func (room *baseRoom) _packetProcess_GameStart(user *roomUser, packet protocol.Packet) int16 {
	NTELIB_LOG_DEBUG("[[Room _packetProcess_GameStart ]")

	errorCode := (int16)(protocol.ERROR_CODE_NONE)
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId

	// 방의 상태가 NONE 인가?
	if room.isStateNone() == false {
		errorCode = (int16)(protocol.ERROR_CODE_ROOM_GAME_START_INVALID_ROOM_STATE)
		_sendRoomGameStartResult(sessionIndex, sessionUniqueId, errorCode)
		return errorCode
	}

	// 유저의 최소 수. 현재는 테스트를 위해 일단 1명이 최소 수
	if room.getCurUserCount() < 1 {
		errorCode = (int16)(protocol.ERROR_CODE_ROOM_GAME_START_NOT_ENOUGH_MEMBERS)
		_sendRoomGameStartResult(sessionIndex, sessionUniqueId, errorCode)
		return errorCode
	}

	// 시작 요청은 방장이 하는가?
	if room.getMasterSessionUniqueId() != sessionUniqueId {
		errorCode = (int16)(protocol.ERROR_CODE_ROOM_GAME_START_NOT_MASTER)
		_sendRoomGameStartResult(sessionIndex, sessionUniqueId, errorCode)
		return errorCode
	}

	// 게임을 시작한다.
	room.changeState(ROOM_STATE_GAME_WAIT_BATTING)

	_sendRoomGameStartResult(sessionIndex, sessionUniqueId, errorCode)
	_sendRoomGameStartNotify(room);

	for _, user := range room._userSessionUniqueIdMap {
		protocol.Update("Plays = Plays+1",bytes.Trim(user.ID[:], "\x00"))
	}


	return errorCode
}

func _sendRoomGameStartResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomGameStartResPacket{ result }
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}

func _sendRoomGameStartNotify(room *baseRoom) {
	notify := protocol.RoomGameStartNtfPacket{room._gameLogic._statusChangeCompletionMillSec}
	notifySendBuf, packetSize := notify.EncodingPacket()
	room.broadcastPacket(packetSize, notifySendBuf, 0)
}



func (room *baseRoom) _packetProcess_GameBatting(user *roomUser, packet protocol.Packet) int16 {
	NTELIB_LOG_DEBUG("[[Room _packetProcess_GameBatting ]")

	errorCode := (int16)(protocol.ERROR_CODE_NONE)
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId
	var battingPacket protocol.RoomGameBattingReqPacket

	// 방의 상태가 배팅 기다림인가?
	if room.isStateGameBattingWait() == false {
		errorCode = (int16)(protocol.ERROR_CODE_ROOM_GAME_BATTING_INVALID_ROOM_STATE)
		goto CheckError
	}

	if battingPacket.Decoding(packet.Data) == false {
		errorCode = (int16)(protocol.ERROR_CODE_ROOM_GAME_BATTING_FAIL_PACKET)
		goto CheckError
	}

	if battingPacket.SelectSide < BATTING_SELECT_PLAYER || battingPacket.SelectSide > BATTING_SELECT_TIE {
		errorCode = (int16)(protocol.ERROR_CODE_ROOM_GAME_BATTING_INVALID_BAT_SELECT)
		goto CheckError
	}

	if battingPacket.SelectSide == user.selectBat {
		errorCode = (int16)(protocol.ERROR_CODE_ROOM_GAME_BATTING_SAME_BAT_SELECT)
		goto CheckError
	}


	// 배팅하고, 모두에게 알린다.
	user.selectBatting(battingPacket.SelectSide)

	_sendRoomGameBattingResult(sessionIndex, sessionUniqueId, errorCode)

	_sendRoomGameBattingNotify(room, user.RoomUniqueId, battingPacket.SelectSide)


	if room.isAllUserBatting() {
		room.endGame()
	}

	return protocol.ERROR_CODE_NONE

CheckError:
	_sendRoomGameBattingResult(sessionIndex, sessionUniqueId, errorCode)
	return errorCode
}

func _sendRoomGameBattingResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomGameBattingResPacket{ result }
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}

func _sendRoomGameBattingNotify(room *baseRoom, roomUserUniqueId uint64, selectSide int8) {
	notify := protocol.RoomGameBattingNtfPacket{ roomUserUniqueId, selectSide }
	notifySendBuf, packetSize := notify.EncodingPacket()
	room.broadcastPacket(packetSize, notifySendBuf, 0)

	notifyToIndividual := protocol.RoomGameBattingNtfPacket{ roomUserUniqueId, 0 }
	notifySendBuf, packetSize = notifyToIndividual.EncodingPacket()
	room.deliveryPacket(roomUserUniqueId, packetSize, notifySendBuf)
}







