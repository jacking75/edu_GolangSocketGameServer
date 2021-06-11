package roomPkg

import (
	"time"

	. "gohipernetFake"

	"main/connectedSessions"
	"main/protocol"
)

func (room *baseRoom) _packetProcess_EnterUser(inValidUser *roomUser, packet protocol.Packet) int16 {
	curTime := time.Now().Unix()
	sessionIndex := packet.UserSessionIndex
	sessionUniqueId := packet.UserSessionUniqueId

	var requestPacket protocol.RoomEnterReqPacket
	(&requestPacket).Decoding(packet.Data)

	//TODO 방의 상태가 NONE이 아니면 들어올 수 없다.

	userID, ok := connectedSessions.GetUserID(sessionIndex)
	if ok == false {
		_sendRoomEnterResult(sessionIndex, sessionUniqueId, 0,0, protocol.ERROR_CODE_ENTER_ROOM_INVALID_USER_ID)
		return protocol.ERROR_CODE_ENTER_ROOM_INVALID_USER_ID
	}

	userInfo := addRoomUserInfo{
		userID,
		sessionIndex,
		sessionUniqueId,
	}
	newUser, addResult := room.addUser(userInfo)

	if addResult != protocol.ERROR_CODE_NONE {
		_sendRoomEnterResult(sessionIndex, sessionUniqueId, 0, 0, addResult)
		return addResult
	}

	if connectedSessions.SetRoomNumber(sessionIndex, sessionUniqueId, room.getNumber(), curTime) == false {
		_sendRoomEnterResult(sessionIndex, sessionUniqueId, 0, 0, protocol.ERROR_CODE_ENTER_ROOM_INVALID_SESSION_STATE)
		return protocol.ERROR_CODE_ENTER_ROOM_INVALID_SESSION_STATE
	}

	if room.getCurUserCount() > 1 {
		//룸의 다른 유저에게 통보한다.
		room._sendNewUserInfoPacket(newUser)

		// 지금 들어온 유저에게 이미 채널에 있는 유저들의 정보를 보낸다
		room._sendUserInfoListPacket(newUser)
	}


	roomNumebr := room.getNumber()
	_sendRoomEnterResult(sessionIndex, sessionUniqueId, roomNumebr, newUser.RoomUniqueId, protocol.ERROR_CODE_NONE)
	return protocol.ERROR_CODE_NONE
}

func _sendRoomEnterResult(sessionIndex int32, sessionUniqueId uint64, roomNumber int32, userUniqueId uint64, result int16) {
	response := protocol.RoomEnterResPacket{
		result,
		roomNumber,
		userUniqueId,
	}
	sendPacket, _ := response.EncodingPacket()
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}


func (room *baseRoom) _sendUserInfoListPacket(user *roomUser) {
	userCount, userInfoListSize, userInfoListBuffer := room.allocAllUserInfo(user.netSessionUniqueId)

	var response protocol.RoomUserListNtfPacket
	response.UserCount = userCount
	response.UserList = userInfoListBuffer
	sendBuf, _ := response.EncodingPacket(userInfoListSize)
	NetLibIPostSendToClient(user.netSessionIndex, user.netSessionUniqueId, sendBuf)
}

func (room *baseRoom) _sendNewUserInfoPacket(user *roomUser) {
	userInfoSize, userInfoListBuffer := room._allocUserInfo(user)

	var response protocol.RoomNewUserNtfPacket
	response.User = userInfoListBuffer
	sendBuf, packetSize := response.EncodingPacket(userInfoSize)
	room.broadcastPacket(int16(packetSize), sendBuf, user.netSessionUniqueId) // 자신을 제외하고 모든 유저에게 Send
}


