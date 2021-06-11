package roomPkg

import (
	"github.com/vmihailenco/msgpack/v4"
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
	roomUserUniqueId := user.RoomUniqueId
	userSessionIndex := user.netSessionIndex
	userSessionUniqueId := user.netSessionUniqueId

	room._removeUser(user)

	room._sendRoomLeaveUserNotify(roomUserUniqueId, userSessionUniqueId)

	curTime := time.Now().Unix()
	connectedSessions.SetRoomNumber(userSessionIndex, userSessionUniqueId, -1, curTime)
}

func _sendRoomLeaveResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	response := protocol.RoomLeaveResPacket{int64(result)}
	bodyData, err := msgpack.Marshal(response)
	if err != nil {
		panic(err)
	}

	sendPacket := protocol.EncodingPacketHeaderInfo(bodyData, uint16(protocol.PACKET_ID_ROOM_LEAVE_RES), 0)
	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)
}

func (room *baseRoom) _sendRoomLeaveUserNotify(roomUserUniqueId uint64, userSessionUniqueId uint64) {
	notifyPacket := protocol.RoomLeaveUserNtfPacket{roomUserUniqueId}

	bodyData, err := msgpack.Marshal(notifyPacket)
	if err != nil {
		return
	}

	sendPacket := protocol.EncodingPacketHeaderInfo(bodyData, uint16(protocol.PACKET_ID_ROOM_LEAVE_USER_NTF), 0)
	packetSize := uint16(len(sendPacket))

	room.broadcastPacket(packetSize, sendPacket, userSessionUniqueId)
}
