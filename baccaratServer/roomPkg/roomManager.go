package roomPkg

import (
	"main/protocol"
)


type RoomManager struct {
	_roomStartNum  int32
	_maxRoomCount  int32
	_roomCountList []int16
	_roomList      []baseRoom
}

func NewRoomManager(config RoomConfig) *RoomManager {
	roomManager := new(RoomManager)
	roomManager._initialize(config)
	return roomManager
}

func (roomMgr *RoomManager) _initialize(config RoomConfig) {
	roomMgr._roomStartNum = config.StartRoomNumber
	roomMgr._maxRoomCount = config.MaxRoomCount
	roomMgr._roomCountList = make([]int16, config.MaxRoomCount)
	roomMgr._roomList = make([]baseRoom, config.MaxRoomCount)

	for i := int32(0); i < roomMgr._maxRoomCount; i++ {
		roomMgr._roomList[i].initialize(i, config)
		roomMgr._roomList[i].settingPacketFunction()
	}
}

func (roomMgr *RoomManager) GetAllChannelUserCount() []int16 {
	maxRoomCount := roomMgr._maxRoomCount
	for i := int32(0); i < maxRoomCount; i++ {
		roomMgr._roomCountList[i] = (int16)(roomMgr._getRoomUserCount(i))
	}

	return roomMgr._roomCountList
}

func (roomMgr *RoomManager) getRoomByNumber(roomNumber int32) *baseRoom {
	roomIndex := roomNumber - roomMgr._roomStartNum

	if roomNumber < 0 || roomIndex >= roomMgr._maxRoomCount {
		return nil
	}

	return &roomMgr._roomList[roomIndex]
}

//  이 함수를 호출할 때의 채널 인덱스는 꼭 호출자가 유효범위인 것을 보증해야 한다.
func (roomMgr *RoomManager) _getRoomUserCount(roomId int32) int {
	return roomMgr._roomList[roomId].getCurUserCount()
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

func (roomMgr *RoomManager) CheckRoomState(curTimeMilliSec int64) {
	//TODO 한번에 모든 방을 다 조사할 필요가 없다. 밀리세컨드 단위로 타이머를 돌게 하고 그룹 단위로 방을 조사한다

	for i := 0; i < (int)(roomMgr._maxRoomCount); i++ {
		roomMgr._roomList[i].checkState(curTimeMilliSec)
	}
}
