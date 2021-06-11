package roomPkg

import (
	"github.com/vmihailenco/msgpack/v4"
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
	err := msgpack.Unmarshal(packet.Data, &requestPacket)
	if err != nil {
		return protocol.ERROR_CODE_PACKET_DECODING_FAIL
	}

	userID, ok := connectedSessions.GetUserID(sessionIndex)
	if ok == false {
		_sendRoomEnterResult(sessionIndex, sessionUniqueId, 0, protocol.ERROR_CODE_ENTER_ROOM_INVALID_USER_ID)
		return protocol.ERROR_CODE_ENTER_ROOM_INVALID_USER_ID
	}

	userInfo := addRoomUserInfo{
		userID,
		sessionIndex,
		sessionUniqueId,
	}
	newUser, addResult := room.addUser(userInfo)

	if addResult != protocol.ERROR_CODE_NONE {
		_sendRoomEnterResult(sessionIndex, sessionUniqueId, 0, addResult)
		return addResult
	}

	if connectedSessions.SetRoomNumber(sessionIndex, sessionUniqueId, room.getNumber(), curTime) == false {
		_sendRoomEnterResult(sessionIndex, sessionUniqueId, 0, protocol.ERROR_CODE_ENTER_ROOM_INVALID_SESSION_STATE)
		return protocol.ERROR_CODE_ENTER_ROOM_INVALID_SESSION_STATE
	}

	if room.getCurUserCount() > 1 {
		//룸의 다른 유저에게 통보한다.
		room._sendNewUserInfoPacket(newUser)

		// 지금 들어온 유저에게 이미 채널에 있는 유저들의 정보를 보낸다
		room._sendUserInfoListPacket(newUser)
	}

	_sendRoomEnterResult(sessionIndex, sessionUniqueId, newUser.RoomUniqueId, protocol.ERROR_CODE_NONE)
	return protocol.ERROR_CODE_NONE
}

func _sendRoomEnterResult(sessionIndex int32, sessionUniqueId uint64, userUniqueId uint64, result int16) {
	response := protocol.RoomEnterResPacket{
		int64(result),
		userUniqueId,
	}

	bodyData, err := msgpack.Marshal(response)
	if err != nil {
		return
	}

	sendPacket := protocol.EncodingPacketHeaderInfo(bodyData, uint16(protocol.PACKET_ID_ROOM_ENTER_RES), 0)
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}

func (room *baseRoom) _sendUserInfoListPacket(user *roomUser) {
	userCount, uniqueIdList, idList := room.allocAllUserInfo(user.netSessionUniqueId)

	var response protocol.RoomUserListNtfPacket
	response.UserCount = int64(userCount)
	response.UniqueId = uniqueIdList
	response.ID = idList
	/*response.UniqueId = make([]uint64, 0, 0)
	response.ID = make([]string, 0, 0)
	for i := 0; i < 2; i++ {
		response.UniqueId = append(response.UniqueId, 1)
		response.ID = append(response.ID, "test1")
	}*/

	bodyData, err := msgpack.Marshal(response)
	if err != nil {
		return
	}
	sendPacket := protocol.EncodingPacketHeaderInfo(bodyData, uint16(protocol.PACKET_ID_ROOM_USER_LIST_NTF), 0)

	NetLibIPostSendToClient(user.netSessionIndex, user.netSessionUniqueId, sendPacket)
}

func (room *baseRoom) _sendNewUserInfoPacket(user *roomUser) {
	var response protocol.RoomNewUserNtfPacket
	response.ID = string(user.ID[0:user.IDLen])
	response.UniqueId = user.RoomUniqueId

	bodyData, err := msgpack.Marshal(response)
	if err != nil {
		panic(err)
	}

	sendPacket := protocol.EncodingPacketHeaderInfo(bodyData, uint16(protocol.PACKET_ID_ROOM_NEW_USER_NTF), 0)
	packetSize := uint16(len(sendPacket))
	room.broadcastPacket(packetSize, sendPacket, user.netSessionUniqueId) // 자신을 제외하고 모든 유저에게 Send
}
