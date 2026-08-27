package roomPkg

import (
	"main/protocol"
)

type RoomManager struct {
	_roomStartNum int32
	_maxRoomCount int32
	_roomList     []baseRoom
}

func NewRoomManager(config RoomConfig) *RoomManager {
	roomManager := new(RoomManager)
	roomManager._initialize(config)
	return roomManager
}

func (roomMgr *RoomManager) _initialize(config RoomConfig) {
	roomMgr._roomStartNum = config.StartRoomNumber
	roomMgr._maxRoomCount = config.MaxRoomCount
	roomMgr._roomList = make([]baseRoom, config.MaxRoomCount)

	for i := int32(0); i < roomMgr._maxRoomCount; i++ {
		roomMgr._roomList[i].initialize(i, config)
		roomMgr._roomList[i].settingPacketFunction()
	}
}

func (roomMgr *RoomManager) getRoomByNumber(roomNumber int32) *baseRoom {
	roomIndex := roomNumber - roomMgr._roomStartNum

	// roomIndex < 0 검사가 없으면 RoomStartNum을 양수로 설정했을 때 roomNumber가 작은 값이라도
	// roomIndex가 음수가 되어 아래 슬라이스 접근에서 음수 인덱스 패닉이 발생할 수 있다.
	if roomIndex < 0 || roomIndex >= roomMgr._maxRoomCount {
		return nil
	}

	return &roomMgr._roomList[roomIndex]
}

func (roomMgr *RoomManager) PacketProcess(roomNumber int32, packet protocol.Packet) {
	isRoomEnterReq := false

	if roomNumber == -1 && packet.Id == protocol.PACKET_ID_ROOM_ENTER_REQ {
		isRoomEnterReq = true

		var requestPacket protocol.RoomEnterReqPacket
		(&requestPacket).Decoding(packet.Data)

		roomNumber = requestPacket.RoomNumber
	}

	room := roomMgr.getRoomByNumber(roomNumber)
	if room == nil {
		protocol.NotifyErrorPacket(packet.UserSessionIndex, packet.UserSessionUniqueId,
			protocol.ERROR_CODE_ROOM_INVALIDE_NUMBER)
		return
	}

	user := room.getUser(packet.UserSessionUniqueId)
	if user == nil && isRoomEnterReq == false {
		protocol.NotifyErrorPacket(packet.UserSessionIndex, packet.UserSessionUniqueId,
			protocol.ERROR_CODE_ROOM_NOT_IN_USER)
		return
	}

	funcCount := len(room._funcPackeIdlist)
	for i := 0; i < funcCount; i++ {
		if room._funcPackeIdlist[i] != packet.Id {
			continue
		}

		room._funclist[i](user, packet)
		return
	}
}

