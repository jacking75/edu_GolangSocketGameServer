package roomPkg

import (
	"go.uber.org/zap"

	"golang_socketGameServer_codelab/chatServer2/protocol"
	. "golang_socketGameServer_codelab/gohipernetFake"
)


type RoomManager struct {
	_roomStartNum  int32
	_maxRoomCount  int32
	_roomCountList []int16
	_roomList      []baseRoom

	_roomPacketDistributor roomPacketDistributor
}

func NewRoomManager(config RoomConfig) *RoomManager {
	roomManager := new(RoomManager)
	roomManager._initialize(config)
	return roomManager
}

func (roomMgr *RoomManager) Start() {
	roomMgr._roomPacketDistributor.startAllRoomProcess()
}

func (roomMgr *RoomManager) Stop() {
	roomMgr._roomPacketDistributor.notifyStopRoomProcess()
}

func (roomMgr *RoomManager) _initialize(config RoomConfig) {
	roomMgr._roomStartNum = config.StartRoomNumber
	roomMgr._maxRoomCount = config.MaxRoomCount
	roomMgr._roomCountList = make([]int16, config.MaxRoomCount)
	roomMgr._roomList = make([]baseRoom, config.MaxRoomCount)
	roomMgr._roomPacketDistributor.init(config.MaxRoomCount, config.RoomCountByGoroutine)

	for i := int32(0); i < roomMgr._maxRoomCount; i++ {
		roomMgr._roomList[i].initialize(i, config)
		roomMgr._roomList[i].settingPacketFunction()
	}

	_log_StartRoomPacketProcess(config.MaxRoomCount, config)
	roomMgr._setRoomPacketPipe(config)

	NTELIB_LOG_INFO("[[[RoomManager initialize]]]", zap.Int32("_maxRoomCount", roomMgr._maxRoomCount))
}

func (roomMgr *RoomManager) _setRoomPacketPipe(config RoomConfig) {
	maxRoomCount := int(config.MaxRoomCount)
	roomCountByGoroutine := int(config.RoomCountByGoroutine)

	for i := 0; i < maxRoomCount; i += roomCountByGoroutine {
		roomlistCount := roomCountByGoroutine
		if maxRoomCount-i < roomCountByGoroutine {
			roomlistCount = maxRoomCount - i
		}

		packetPipe := NewRoomPacketPipe(roomlistCount, config)

		for n := 0; n < roomlistCount; n++ {
			index := i + n
			packetPipe.setRoom(&roomMgr._roomList[index])
			roomMgr._roomPacketDistributor.setPacketPipe(index, packetPipe)
		}

		roomMgr._roomPacketDistributor.addPacketPipeInst(packetPipe)
	}
}

func (roomMgr *RoomManager) GetAllChannelUserCount() []int16 {
	maxRoomCount := roomMgr._maxRoomCount
	for i := int32(0); i < maxRoomCount; i++ {
		roomMgr._roomCountList[i] = (int16)(roomMgr._getRoomUserCount(i))
	}

	return roomMgr._roomCountList
}

func (roomMgr *RoomManager) GetRoom(roomIndex int32) *baseRoom {
	if roomIndex < 0 || roomIndex >= roomMgr._maxRoomCount {
		return nil
	}

	return &roomMgr._roomList[roomIndex]
}

//  이 함수를 호출할 때의 채널 인덱스는 꼭 호출자가 유효범위인 것을 보증해야 한다.
func (roomMgr *RoomManager) _getRoomUserCount(roomId int32) int32 {
	return roomMgr._roomList[roomId].getCurUserCount()
}


func RoomNumberToIndex(roomStartNumber int32, roomNumber int32) int32 {
	return roomNumber - roomStartNumber
}

func (roomMgr *RoomManager) PushPacket(roomlIndex int32, packet protocol.Packet) {
	roomMgr._roomPacketDistributor.PushPacket(roomlIndex, packet)
}

func (roomMgr *RoomManager) PushInternalPacket(packet protocol.InternalPacket) {
	roomMgr._roomPacketDistributor.PushInternalPacket(packet)
}

func (roomMgr *RoomManager) PushInternalPacketAllRooms(packet protocol.InternalPacket) {
	roomMgr._roomPacketDistributor.PushInternalPacketAllRooms(packet)
}

func (roomMgr *RoomManager) PushInternalPacketRange(startIndex int32, endIndex int32, packet protocol.InternalPacket) {
	roomMgr._roomPacketDistributor.PushInternalPacketRange(startIndex, endIndex, packet)
}

func (roomMgr *RoomManager) DisConnectedUser(sessionUniqueId uint64, roomNum int32) {
	roomIndex := RoomNumberToIndex(roomMgr._roomStartNum, roomNum)

	internalPacket := protocol.InternalPacketDisConnectedUserToRoom{
		sessionUniqueId,
		roomNum,
	}

	bodyData, bodySize := internalPacket.Encoding()
	packet := protocol.InternalPacket{
		roomIndex,
		protocol.INTERNAL_PACKET_ID_DISCONNECTED_USER_TO_ROOM,
		bodySize,
		bodyData,
	}

	roomMgr.PushInternalPacket(packet)
}

func getRoomByRoomIndex(count int32, rooms []*baseRoom, roomIndex int32) *baseRoom {
	for i := int32(0); i < count; i++ {
		if rooms[i].getIndex() == roomIndex {
			return rooms[i]
		}
	}

	return nil
}



func _log_StartRoomPacketProcess(maxRoomCount int32, config RoomConfig) {
	NTELIB_LOG_INFO("[[[RoomManager _startRoomPacketProcess]]]",
		zap.Int32("maxRoomCount", maxRoomCount),
		zap.Int32("StartRoomNumber", config.StartRoomNumber),
		zap.Int32("MaxUserCount", config.MaxUserCount),
		zap.Int32("ChanPacketBufferCount", config.ChanPacketBufferCount),
		zap.Int32("InternalPacketChanBufferCount", config.InternalPacketChanBufferCount))
}
