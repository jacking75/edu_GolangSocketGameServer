package connectedSessions

import (
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	. "gohipernetFake"
)

// 스레드 세이프 해야 한다.
type Manager struct {
	_UserIDsessionMap *sync.Map

	_config CheckSessionStateConfig
	_maxSessionCount int32
	_sessionList     []*session

	_maxUserCount int32

	_currentLoginUserCount int32
}

var _manager Manager

func Init(maxSessionCount int32, maxUserCount int32, config CheckSessionStateConfig) bool {
	_manager._UserIDsessionMap = new(sync.Map)
	_manager._maxUserCount = maxUserCount

	_manager._maxSessionCount = maxSessionCount
	_manager._sessionList = make([]*session, maxSessionCount)

	for i := (int32)(0); i < maxSessionCount; i++ {
		_manager._sessionList[i] = new(session)
		_manager._sessionList[i].Init(i)
	}

	_manager._config = config
	return true
}

func Start() {
	go startCheckConnectedSessionsState_goroutine(_manager._config)
}

func AddSession(sessionIndex int32, sessionUniqueID uint64) bool {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return false
	}

	if _manager._sessionList[sessionIndex].GetConnectTimeSec() > 0 {
		NTELIB_LOG_ERROR("already connected session", zap.Int32("sessionIndex", sessionIndex))
		return false
	}

	// 방어적인 목적으로 한번 더 Clear 한다
	_manager._sessionList[sessionIndex].Clear()

	var curTime int64 = time.Now().Unix()
	_manager._sessionList[sessionIndex].SetRecentlyReceivedTimeSec(curTime)
	_manager._sessionList[sessionIndex].SetConnectTimeSec(curTime, sessionUniqueID)

	return true
}

func RemoveSession(sessionIndex int32, isLoginedUser bool) bool {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return false
	}

	if isLoginedUser {
		atomic.AddInt32(&_manager._currentLoginUserCount, -1)

		userID := string(_manager._sessionList[sessionIndex].getUserID())
		_manager._UserIDsessionMap.Delete(userID)
	}

	_manager._sessionList[sessionIndex].Clear()

	return true
}

func CurrentLoginUserCount() int32 {
	count := atomic.LoadInt32(&_manager._currentLoginUserCount)
	return count
}

func _validSessionIndex(index int32) bool {
	if index < 0 || index >= _manager._maxSessionCount {
		return false
	}
	return true
}

func GetNetworkUniqueID(sessionIndex int32) uint64 {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return 0
	}

	return _manager._sessionList[sessionIndex].GetNetworkUniqueID()
}

func GetUserID(sessionIndex int32) ([]byte, bool) {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return nil, false
	}

	return _manager._sessionList[sessionIndex].getUserID(), true
}

func GetUserIDLength(sessionIndex int32) int8 {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return -1
	}

	return _manager._sessionList[sessionIndex].getUserIDLength()
}

func isConnectedSession(sessionIndex int32) bool {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return false
	}

	if _manager._sessionList[sessionIndex].GetConnectTimeSec() == 0 {
		return false
	}

	return true
}

func GetConnectTimeSec(sessionIndex int32) int64 {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return 0
	}

	return _manager._sessionList[sessionIndex].GetConnectTimeSec()
}

func EnableLogin(sessionIndex int32) bool {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return false
	}
	return _manager._sessionList[sessionIndex].IsAuth() == false
}

func SetLogin(sessionIndex int32, sessionUniqueId uint64, userID []byte, curTimeSec int64) bool {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return false
	}

	newUserID := string(userID)
	if _, ok := _manager._UserIDsessionMap.Load(newUserID); ok {
		return false
	}

	_manager._sessionList[sessionIndex].SetUser(sessionUniqueId, userID, curTimeSec)
	_manager._UserIDsessionMap.Store(newUserID, _manager._sessionList[sessionIndex])

	atomic.AddInt32(&_manager._currentLoginUserCount, 1)
	return true
}

func IsLoginUser(sessionIndex int32) bool {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return false
	}

	return _manager._sessionList[sessionIndex].IsAuth()
}

func SetRoomEntering(sessionIndex int32, roomNum int32) bool {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return false
	}

	return _manager._sessionList[sessionIndex].setRoomEntering(roomNum)
}

func SetRoomNumber(sessionIndex int32, sessionUniqueId uint64, roomNum int32, curTimeSec int64) bool {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return false
	}

	return _manager._sessionList[sessionIndex].setRoomNumber(sessionUniqueId, roomNum, curTimeSec)
}

func GetRoomNumber(sessionIndex int32) (int32, int32) {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return -1, -1
	}
	return _manager._sessionList[sessionIndex].getRoomNumber()
}

func SetRecentlyReceivedTimeSec(sessionIndex int32, time int64) {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
	}

	_manager._sessionList[sessionIndex].SetRecentlyReceivedTimeSec(time)
}

func GetRecentlyReceivedTimeSec(sessionIndex int32) int64 {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
		return -1
	}
	return _manager._sessionList[sessionIndex].GetRecentlyReceivedTimeSec()
}

func SetDisConnectWaitStartTimeSec(sessionIndex int32, timeSec int64) {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
	}

	_manager._sessionList[sessionIndex].SetDisConnectWaitStartTimeSec(timeSec)
}

func GetDisConnectWaitStartTimeSec(sessionIndex int32) int64 {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
	}

	return _manager._sessionList[sessionIndex].GetDisConnectWaitStartTimeSec()
}

func AddRequestPerSecondTime(sessionIndex int32, curTime int64) int32 {
	if _validSessionIndex(sessionIndex) == false {
		NTELIB_LOG_ERROR("Invalid sessionIndex", zap.Int32("sessionIndex", sessionIndex))
	}

	return _manager._sessionList[sessionIndex].AddRequestPerSecondTime(curTime)
}
