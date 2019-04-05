package roomPkg

import (
	"time"

	"go.uber.org/zap"

	"golang_socketGameServer_codelab/chatServer2/protocol"
	. "golang_socketGameServer_codelab/gohipernetFake"
)

type roomPacketPipe struct {
	_roomCount int
	_roomRefList []*baseRoom

	_channelPacketCount    int32
	_chanPacket            chan protocol.Packet
	_chanInternalPacket    chan protocol.InternalPacket
	_chanMemebrPacket      chan RoomMemberPacket
	_onDoneTerminateNotify chan struct{}
}

func NewRoomPacketPipe(roomCount int, config RoomConfig) *roomPacketPipe {
	packetPipe := new(roomPacketPipe)

	packetPipe._roomRefList = make([]*baseRoom, roomCount)
	packetPipe._chanPacket = make(chan protocol.Packet, config.ChanPacketBufferCount)
	packetPipe._chanInternalPacket = make(chan protocol.InternalPacket, config.InternalPacketChanBufferCount)
	packetPipe._chanMemebrPacket = make(chan RoomMemberPacket, config.MaxUserCount)
	packetPipe._onDoneTerminateNotify = make(chan struct{})

	return packetPipe
}

func (packetPipe *roomPacketPipe) setRoom(room *baseRoom) {
	packetPipe._roomRefList[packetPipe._roomCount] = room
	packetPipe._roomCount++
}

func (packetPipe *roomPacketPipe) getRoomCount() int {
	return packetPipe._roomCount
}

func (packetPipe *roomPacketPipe) getRoomByRoomNum(roomNum int32) *baseRoom {
	for i := 0; i < packetPipe._roomCount; i++ {
		if packetPipe._roomRefList[i].getNumber() == roomNum {
			return packetPipe._roomRefList[i]
		}
	}

	return nil
}

func (packetPipe *roomPacketPipe) Stop() {
	close(packetPipe._onDoneTerminateNotify)
}


func (packetPipe *roomPacketPipe) roomProcess_goroutine() {
	NTELIB_LOG_INFO("start Park rooms process goroutine")

	for {
		if packetPipe.roomProcess_goroutine_Impl() {
			NTELIB_LOG_INFO("Wanted Stop rooms process goroutine")
			break
		}
	}

	NTELIB_LOG_INFO("Stop rooms process goroutine")
}


func (packetPipe *roomPacketPipe) roomProcess_goroutine_Impl() bool {
	IsWantedTermination := false
	defer PrintPanicStack()

	secondTimeticker := time.NewTicker(time.Second)
	defer secondTimeticker.Stop()

loop:
	for {
		select {
		// select 안의 _chanPacket 과 _chanInternalPacket 채널을 하나로 통합한다(성능에 더 좋음)
		// for packet := range packetQueue
		case packet := <-packetPipe._chanPacket:
			{
				packetPipe._packetProcess(packet)
			}
		case internalPacket := <-packetPipe._chanInternalPacket:
			{
				packetPipe._internalPacketProcess(internalPacket)
			}
		case _ = <-secondTimeticker.C:
			{
				packetPipe._secondTimeEvent()
			}
		case <-packetPipe._onDoneTerminateNotify:
			IsWantedTermination = true
			break loop
		}
	}

	return IsWantedTermination
}

func (packetPipe *roomPacketPipe) _packetProcess(packet protocol.Packet) int16 {
	NTELIB_LOG_DEBUG("[[Room - _packetProcess]]", zap.Int16("PacketID", packet.Id))

	room := packetPipe.getRoomByRoomNum(packet.RoomNumber)
	if room == nil {
		protocol.NotifyErrorPacket(packet.UserSessionIndex, packet.UserSessionUniqueId,
						protocol.ERROR_CODE_ROOM_INVALIDE_NUMBER)
		return protocol.ERROR_CODE_ROOM_INVALIDE_NUMBER
	}

	user := room.getUser(packet.UserSessionUniqueId)

	if user == nil && packet.Id != protocol.PACKET_ID_ROOM_ENTER_REQ {
		protocol.NotifyErrorPacket(packet.UserSessionIndex, packet.UserSessionUniqueId,
							protocol.ERROR_CODE_ROOM_NOT_IN_USER)
		return protocol.ERROR_CODE_ROOM_NOT_IN_USER
	}

	funcCount := len(room._funcPackeIdlist)
	for i := 0; i < funcCount; i++ {
		if room._funcPackeIdlist[i] != packet.Id {
			continue
		}

		result := room._funclist[i](user, packet)
		if result != protocol.ERROR_CODE_NONE {
			NTELIB_LOG_DEBUG("[[Room - _packetProcess - Fail]]",
				zap.Int16("PacketID", packet.Id), zap.Int16("Error", result))
		}

		return protocol.ERROR_CODE_NONE
	}

	NTELIB_LOG_DEBUG("[[Room - _packetProcess - Fail(Not Registered)]]", 								zap.Int16("PacketID", packet.Id))
	return protocol.ERROR_CODE_ROOM_NOT_REGISTED_PACKET_ID
}

func (packetPipe *roomPacketPipe) _internalPacketProcess(internalPacket protocol.InternalPacket) {
	roomList := packetPipe._roomRefList

	if internalPacket.RoomIndex == -1 { // 모두에게 보낼 경우
		for i := 0; i < packetPipe._roomCount; i++ {
			roomList[i]._internalPacketProcess(internalPacket)
		}
	} else {
		for i := 0; i < packetPipe._roomCount; i++ {
			if roomList[i]._index == internalPacket.RoomIndex {
				roomList[i]._internalPacketProcess(internalPacket)
				break
			}
		}
	}
}

func (packetPipe *roomPacketPipe) _secondTimeEvent() {
	for i := 0; i < packetPipe._roomCount; i++ {
		packetPipe._roomRefList[i].secondTimeEvent()
	}
}
