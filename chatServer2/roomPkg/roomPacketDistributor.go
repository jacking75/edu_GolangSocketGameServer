package roomPkg

import (
	"go.uber.org/zap"

	"golang_socketGameServer_codelab/chatServer2/protocol"
	. "golang_socketGameServer_codelab/gohipernetFake"
)

type roomPacketDistributor struct {
	_roomCount                      int32
	_roomIndexListRoomPacketPipeRef []*roomPacketPipe //채널과 1:1로 패킷파이프를 할당한다.

	_roomPacketPipeInstCount int
	_roomPacketPipeInstList  []*roomPacketPipe // 채널 수보다 패킷파이프 수가 더 작을 수 있으므로 실제 생성된 것을 따로 관리한다.
}

func (distributor *roomPacketDistributor) init(roomCount int32, roomCountByGoroutine int32) {
	distributor._roomCount = roomCount
	distributor._roomIndexListRoomPacketPipeRef = make([]*roomPacketPipe, roomCount)

	passinglyPacketPipeInstCount := (roomCount / roomCountByGoroutine) + 2
	distributor._roomPacketPipeInstList = make([]*roomPacketPipe, passinglyPacketPipeInstCount)
}

func (distributor *roomPacketDistributor) setPacketPipe(roomIndex int, pipe *roomPacketPipe) {
	distributor._roomIndexListRoomPacketPipeRef[roomIndex] = pipe
}

func (distributor *roomPacketDistributor) addPacketPipeInst(pipe *roomPacketPipe) {
	distributor._roomPacketPipeInstList[distributor._roomPacketPipeInstCount] = pipe
	distributor._roomPacketPipeInstCount++
}

func (distributor *roomPacketDistributor) startAllRoomProcess() {
	for i := 0; i < distributor._roomPacketPipeInstCount; i++ {
		go distributor._roomPacketPipeInstList[i].roomProcess_goroutine()
	}
}

func (distributor *roomPacketDistributor) notifyStopRoomProcess() {
	for i := 0; i < distributor._roomPacketPipeInstCount; i++ {
		distributor._roomPacketPipeInstList[i].Stop()
	}
}

func (distributor *roomPacketDistributor) ValidRoomIndex(roomIndex int32) bool {
	if roomIndex < 0 || roomIndex >= distributor._roomCount {
		return false
	}

	return true
}

func (distributor *roomPacketDistributor) PushMemberPacket(roomIndex int32, packet RoomMemberPacket) {
	if distributor.ValidRoomIndex(roomIndex) == false {
		NTELIB_LOG_ERROR("fail pushMemberPacket", zap.Int32("roomIndex", roomIndex))
		return
	}

	distributor._roomIndexListRoomPacketPipeRef[roomIndex]._chanMemebrPacket <- packet
}

func (distributor *roomPacketDistributor) PushPacket(roomIndex int32, packet protocol.Packet) {
	if distributor.ValidRoomIndex(roomIndex) == false {
		NTELIB_LOG_ERROR("fail PushPacket", zap.Int32("roomIndex", roomIndex), zap.Int16("PacketID", packet.Id))
		return
	}

	//gohipernet.LOG_DEBUG("PushPacket", zap.Int32("roomIndex", roomIndex), zap.Int16("packetID", packet.Id))
	distributor._roomIndexListRoomPacketPipeRef[roomIndex]._chanPacket <- packet
}

func (distributor *roomPacketDistributor) PushInternalPacketAllRooms(packet protocol.InternalPacket) {
	for i := 0; i < distributor._roomPacketPipeInstCount; i++ {
		distributor._roomPacketPipeInstList[i]._chanInternalPacket <- packet
	}
}

func (distributor *roomPacketDistributor) PushInternalPacket(packet protocol.InternalPacket) {
	roomIndex := packet.RoomIndex

	if distributor.ValidRoomIndex(roomIndex) == false {
		NTELIB_LOG_ERROR("fail pushMemberPacket", zap.Int32("roomIndex", roomIndex))
		return
	}

	distributor._roomIndexListRoomPacketPipeRef[roomIndex]._chanInternalPacket <- packet
}

func (distributor *roomPacketDistributor) PushInternalPacketRange(roomStartIndex int32,
	roomEndIndex int32,
	packet protocol.InternalPacket,
	) {
	if distributor.ValidRoomIndex(roomStartIndex) == false {
		NTELIB_LOG_ERROR("fail pushMemberPacket", zap.Int32("roomStartIndex", roomStartIndex))
		return
	}

	if distributor.ValidRoomIndex(roomEndIndex) == false {
		NTELIB_LOG_ERROR("ChannelEndIndex Over", zap.Int32("roomEndIndex", roomEndIndex), zap.Int32("roomCount", distributor._roomCount))
		return
	}

	for i := roomStartIndex; i <= roomEndIndex; i++ {
		packet.RoomIndex = i
		distributor._roomIndexListRoomPacketPipeRef[i]._chanInternalPacket <- packet
	}
}
