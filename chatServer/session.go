package connectedSessions

import (
	"sync/atomic"

	"main/protocol"
)

type session struct {
	_index int32

	_networkUniqueID uint64 //네트워크 세션의 유니크 ID

	_userID       [protocol.MAX_USER_ID_BYTE_LENGTH]byte
	_userIDLength int8

	_connectTimeSec    int64 // 연결된 시간
	_RoomNum           int32 //
	_RoomNumOfEntering int32 // 현재 입장 중인 룸의 번호
}

func (session *session) Init(index int32) {
	session._index = index
	session.Clear()
}

func (session *session) _ClearUserId() {
	session._userIDLength = 0
}

func (session *session) Clear() {
	session._ClearUserId()
	session.setRoomNumber(0, -1, 0)
	session.SetConnectTimeSec(0, 0)
}

func (session *session) GetIndex() int32 {
	return session._index
}

func (session *session) GetNetworkUniqueID() uint64 {
	return atomic.LoadUint64(&session._networkUniqueID)
}

func (session *session) validNetworkUniqueID(uniqueId uint64) bool {
	return atomic.LoadUint64(&session._networkUniqueID) == uniqueId
}

func (session *session) GetNetworkInfo() (int32, uint64) {
	index := session.GetIndex()
	uniqueID := atomic.LoadUint64(&session._networkUniqueID)
	return index, uniqueID
}

func (session *session) setUserID(userID []byte) {
	session._userIDLength = int8(len(userID))
	copy(session._userID[:], userID)
}

func (session *session) getUserID() []byte {
	return session._userID[0:session._userIDLength]
}

func (session *session) getUserIDLength() int8 {
	return session._userIDLength
}

func (session *session) SetConnectTimeSec(timeSec int64, uniqueID uint64) {
	atomic.StoreInt64(&session._connectTimeSec, timeSec)
	atomic.StoreUint64(&session._networkUniqueID, uniqueID)
}

func (session *session) GetConnectTimeSec() int64 {
	return atomic.LoadInt64(&session._connectTimeSec)
}

func (session *session) SetUser(sessionUniqueId uint64,
	userID []byte,
	curTimeSec int64,
) {
	session.setUserID(userID)
	session.setRoomNumber(sessionUniqueId, -1, curTimeSec) // 방어적인 목적으로 채널 번호 초기화
}

func (session *session) IsAuth() bool {
	if session._userIDLength > 0 {
		return true
	}

	return false
}

func (session *session) setRoomEntering(roomNum int32) bool {
	if atomic.CompareAndSwapInt32(&session._RoomNumOfEntering, -1, roomNum) == false {
		return false
	}

	return true
}

func (session *session) setRoomNumber(sessionUniqueId uint64, roomNum int32, curTimeSec int64) bool {
	if roomNum == -1 {
		atomic.StoreInt32(&session._RoomNum, roomNum)
		atomic.StoreInt32(&session._RoomNumOfEntering, roomNum)
		return true
	}

	if sessionUniqueId != 0 && session.validNetworkUniqueID(sessionUniqueId) == false {
		return false

	}
	// 입력이 -1이 아닌경우 -1이 아닐 때만 compareswap으로 변경한다. 실패하면 채널 입장도 실패이다.
	if atomic.CompareAndSwapInt32(&session._RoomNum, -1, roomNum) == false {
		return false
	}

	atomic.StoreInt32(&session._RoomNumOfEntering, roomNum)
	return true
}

func (session *session) getRoomNumber() (int32, int32) {
	roomNum := atomic.LoadInt32(&session._RoomNum)
	roomNumOfEntering := atomic.LoadInt32(&session._RoomNum)
	return roomNum, roomNumOfEntering
}
