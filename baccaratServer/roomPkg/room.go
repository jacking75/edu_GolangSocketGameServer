package roomPkg

import (
	"sync"
	"time"

	. "gohipernetFake"
	"main/protocol"
)

type baseRoom struct {
	_index  int32
	_number int32 // 채널 고유 번호
	_config RoomConfig

	_curState int

	_gameLogic baccaratGame

	_roomUserUnqieuIdSeq uint64

	_userPool *sync.Pool

	//자료구조를 배열로 바꾸는 것이 좋음
	_userSessionUniqueIdMap map[uint64]*roomUser //range 순회 시 복사 비용이 발생해서 포인터 값을 사용한다.

	_masterUserSessionUniqueId uint64 // 방장의 네트워크 유니크ID

	_funcPackeIdlist []int16
	_funclist        []func(*roomUser, protocol.Packet) int16

	enterUserNotify func(int64, int32)
	leaveUserNotify func(int64)
}

func (room *baseRoom) getIndex() int32 {
	return room._index
}

func (room *baseRoom) getNumber() int32 {
	return room._number
}

func (room *baseRoom) isStateNone() bool {
	return room._curState == ROOM_STATE_NOE
}

func (room *baseRoom) isStateGameResult() bool {
	return room._curState == ROOM_STATE_GAME_RESULT
}

func (room *baseRoom) isStateGameBattingWait() bool {
	return room._curState == ROOM_STATE_GAME_WAIT_BATTING
}

func (room *baseRoom) getCurUserCount() int {
	count := len(room._userSessionUniqueIdMap)
	return count
}

func (room *baseRoom) changeState(state int) {
	room._curState = state

	if room._curState == ROOM_STATE_NOE {
		room._gameLogic.clear()
	} else if room._curState == ROOM_STATE_GAME_WAIT_BATTING {
		room._gameLogic.setBattingWaitTime(time.Now().Unix())
	}
}

func (room *baseRoom) getMasterSessionUniqueId() uint64 {
	return room._masterUserSessionUniqueId
}

func (room *baseRoom) getMasterUser(userSessionUniqueId uint64) *roomUser {
	return room.getUser(userSessionUniqueId)
}

func (room *baseRoom) _settingMasterUser() {
	if room.getCurUserCount() < 1 {
		room._masterUserSessionUniqueId = 0
	}

	count := (uint64)(0)
	masterUniqueId := (uint64)(0)

	for _, user := range room._userSessionUniqueIdMap {
		if user.RoomUniqueId > count {
			count = user.RoomUniqueId
			masterUniqueId = user.netSessionUniqueId
		}
	}

	room._masterUserSessionUniqueId = masterUniqueId
}

func (room *baseRoom) generateUserUniqueId() uint64 {
	room._roomUserUnqieuIdSeq++
	uniqueId := room._roomUserUnqieuIdSeq
	return uniqueId
}

func (room *baseRoom) initialize(index int32, config RoomConfig) {
	room._initialize(index, config)

	room._initUserPool()
	room._userSessionUniqueIdMap = make(map[uint64]*roomUser)

	room._gameLogic.init()
}

func (room *baseRoom) _initialize(index int32, config RoomConfig) {
	room._number = config.StartRoomNumber + index
	room._index = index
	room._config = config
	room._curState = ROOM_STATE_NOE
	room._masterUserSessionUniqueId = 0
}

func (room *baseRoom) EnableEnterUser() bool {
	if room._IsFullUser() {
		return false
	}

	return true
}

func (room *baseRoom) settingPacketFunction() {
	maxFuncListCount := 16
	room._funclist = make([]func(*roomUser, protocol.Packet) int16, 0, maxFuncListCount)
	room._funcPackeIdlist = make([]int16, 0, maxFuncListCount)

	room._addPacketFunction(protocol.PACKET_ID_ROOM_ENTER_REQ, room._packetProcess_EnterUser)
	room._addPacketFunction(protocol.PACKET_ID_ROOM_LEAVE_REQ, room._packetProcess_LeaveUser)
	room._addPacketFunction(protocol.PACKET_ID_ROOM_CHAT_REQ, room._packetProcess_Chat)
	room._addPacketFunction(protocol.PACKET_ID_ROOM_RELAY_REQ, room._packetProcess_Relay)
	room._addPacketFunction(protocol.PACKET_ID_GAME_START_REQ, room._packetProcess_GameStart)
	room._addPacketFunction(protocol.PACKET_ID_GAME_BATTING_REQ, room._packetProcess_GameBatting)
}

func (room *baseRoom) _addPacketFunction(packetID int16,
	packetFunc func(*roomUser, protocol.Packet,
	) int16) {
	room._funclist = append(room._funclist, packetFunc)
	room._funcPackeIdlist = append(room._funcPackeIdlist, packetID)
}

func (room *baseRoom) _initUserPool() {
	room._userPool = &sync.Pool{
		New: func() interface{} {
			user := new(roomUser)
			return user
		},
	}
}

func (room *baseRoom) _getUserObject() *roomUser {
	userObject := room._userPool.Get().(*roomUser)
	return userObject
}

func (room *baseRoom) _putUserObject(user *roomUser) {
	room._userPool.Put(user)
}

func (room *baseRoom) addUser(userInfo addRoomUserInfo) (*roomUser, int16) {
	if room._IsFullUser() {
		return nil, protocol.ERROR_CODE_ENTER_ROOM_USER_FULL
	}

	if room.getUser(userInfo.netSessionUniqueId) != nil {
		return nil, protocol.ERROR_CODE_ENTER_ROOM_DUPLCATION_USER
	}

	user := room._getUserObject()
	user.init(userInfo.userID, room.generateUserUniqueId())
	user.SetNetworkInfo(userInfo.netSessionIndex, userInfo.netSessionUniqueId)
	user.packetDataSize = user.PacketDataSize()

	room._userSessionUniqueIdMap[user.netSessionUniqueId] = user

	if room.getCurUserCount() == 1 {
		room._masterUserSessionUniqueId = userInfo.netSessionUniqueId
	}

	return user, protocol.ERROR_CODE_NONE
}

func (room *baseRoom) _IsFullUser() bool {
	if room.getCurUserCount() == (int)(room._config.MaxUserCount) {
		return true
	}

	return false
}

func (room *baseRoom) _removeUser(user *roomUser) {
	delete(room._userSessionUniqueIdMap, user.netSessionUniqueId)
	room._removeUserObject(user)
}

func (room *baseRoom) _removeUserObject(user *roomUser) {
	room._putUserObject(user)
	room._settingMasterUser()
}

func (room *baseRoom) getUser(sessionUniqueId uint64) *roomUser {
	if user, ok := room._userSessionUniqueIdMap[sessionUniqueId]; ok {
		return user
	}

	return nil
}

// 함수 이름 바꾸는 것이 좋을 듯
func (room *baseRoom) allocAllUserInfo(exceptSessionUniqueId uint64) (userCount int8, dataSize int16, dataBuffer []byte) {
	for _, user := range room._userSessionUniqueIdMap {
		if user.netSessionUniqueId == exceptSessionUniqueId {
			continue
		}

		userCount++
		dataSize += user.packetDataSize
	}

	dataBuffer = make([]byte, dataSize)
	writer := MakeWriter(dataBuffer, true)

	for _, user := range room._userSessionUniqueIdMap {
		if user.netSessionUniqueId == exceptSessionUniqueId {
			continue
		}

		_writeUserInfo(&writer, user)
	}

	return userCount, dataSize, dataBuffer
}

// 유저 하나에게 보낼 때는 통으로 보낸다
func (room *baseRoom) _allocUserInfo(user *roomUser) (dataSize int16, dataBuffer []byte) {
	dataSize = user.packetDataSize
	dataBuffer = make([]byte, dataSize)
	writer := MakeWriter(dataBuffer, true)
	_writeUserInfo(&writer, user)

	return dataSize, dataBuffer
}

func _writeUserInfo(writer *RawPacketData, user *roomUser) {
	writer.WriteU64(user.RoomUniqueId)
	writer.WriteS8(user.IDLen)
	writer.WriteBytes(user.ID[0:user.IDLen])
}

func (room *baseRoom) _disConnectedUser(sessionUniqueId uint64) bool {
	user := room.getUser(sessionUniqueId)
	if user == nil {
		return false
	}

	room._leaveUserProcess(user)
	return true
}

func (room *baseRoom) secondTimeEvent() {
	//TODO 주기적으로 방의 유저가 연결 되어 있는지 확인 필요
}

func (room *baseRoom) broadcastPacket(packetSize int16,
	sendPacket []byte,
	exceptSessionUniqueId uint64,
) {

	for _, user := range room._userSessionUniqueIdMap {
		if user.netSessionUniqueId == exceptSessionUniqueId {
			continue
		}

		NetLibIPostSendToClient(user.netSessionIndex, user.netSessionUniqueId, sendPacket)
	}
}

func (room *baseRoom) disConnectedUser(sessionUniqueId uint64) int16 {
	if room._disConnectedUser(sessionUniqueId) == false {
		return protocol.ERROR_CODE_LEAVE_ROOM_INTERNAL_INVALID_USER
	}

	return protocol.ERROR_CODE_NONE
}

func (room *baseRoom) isAllUserBatting() bool {
	count := 0
	for _, user := range room._userSessionUniqueIdMap {
		if user.selectBat != BATTING_SELECT_NONE {
			count++
		}
	}

	return count == len(room._userSessionUniqueIdMap)
}

func (room *baseRoom) endGame() {
	gameResult := room._gameLogic.doBaccarat()

	notify := protocol.RoomGameResultNtfPacket{}
	copy(notify.CardsBanker[:], gameResult.cardsBanker[:])
	copy(notify.CardsPlayer[:], gameResult.cardsPlayer[:])
	notify.PlayerScore = gameResult.playerScore
	notify.BankerScore = gameResult.bankerScore
	notify.Result = gameResult.result

	notifySendBuf, packetSize := notify.EncodingPacket()
	room.broadcastPacket(packetSize, notifySendBuf, 0)

	room.changeState(ROOM_STATE_GAME_RESULT)
}

func (room *baseRoom) checkState(curTimeMilliSec int64) {
	if room.isStateNone() {
		return
	} else if room.isStateGameBattingWait() {
		if room._gameLogic.isTimeOver(curTimeMilliSec) {
			room.endGame()
		}
	} else if room.isStateGameResult() {
		if room._gameLogic.isTimeOver(curTimeMilliSec) {
			room.changeState(ROOM_STATE_NOE)
		}
	}

}
