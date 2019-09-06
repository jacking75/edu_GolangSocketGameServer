package roomPkg

import "main/protocol"

type RoomConfig struct {
	StartRoomNumber int32
	MaxRoomCount    int32
	MaxUserCount    int32
}

type roomUser struct {
	netSessionIndex    int32
	netSessionUniqueId uint64

	// <<< 다른 유저에게 알려줘야 하는 정보
	RoomUniqueId uint64
	IDLen        int8
	ID           [protocol.MAX_USER_ID_BYTE_LENGTH]byte
	// >>> 다른 유저에게 알려줘야 하는 정보
	packetDataSize int16 // 다른 유저에게 알려줘야 하는 정보 의 크기
}

func (user *roomUser) init(userID []byte, uniqueId uint64) {
	idlen := len(userID)

	user.IDLen = int8(idlen)
	copy(user.ID[:], userID)

	user.RoomUniqueId = uniqueId
}

func (user *roomUser) SetNetworkInfo(sessionIndex int32, sessionUniqueId uint64) {
	user.netSessionIndex = sessionIndex
	user.netSessionUniqueId = sessionUniqueId
}

func (user *roomUser) PacketDataSize() int16 {
	return int16(1) + int16(user.IDLen) + 8
}

type addRoomUserInfo struct {
	userID []byte

	netSessionIndex    int32
	netSessionUniqueId uint64
}
