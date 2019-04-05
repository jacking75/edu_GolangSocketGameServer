package connectedSessions

import (
	"time"

	"go.uber.org/zap"

	"golang_socketGameServer_codelab/chatServer2/protocol"
	. "golang_socketGameServer_codelab/gohipernetFake"
)

type CheckSessionStateConfig struct {
	CheckCountAtOnce 		int
	CheckPeriodMillSec		int

	LoginWaitTimeSec          int
	DisConnectWaitTimeSec     int
	RoomEnterWaitTimeSec      int
	PingWaitTimeSec           int
	MaxRequestCountPerSecond  int //클라이언트의 초당 최대 요청 수. 이것 보다 많으면 블랙 유저
}

func startCheckConnectedSessionsState_goroutine(config CheckSessionStateConfig) {
	NTELIB_LOG_INFO("Start CheckConnectedSessionsState goroutine")

	for {
		if _CheckSessionState_goroutine_Impl(config) {
			NTELIB_LOG_INFO("Wanted Stop CheckConnectedSessionsState goroutine")
			break
		}
	}

	NTELIB_LOG_INFO("Stop CheckConnectedSessionsState goroutine")
}

func _CheckSessionState_goroutine_Impl(config CheckSessionStateConfig) bool {
	IsWantedTermination := false

	interval := time.Duration(config.CheckPeriodMillSec) * time.Millisecond
	ticker := time.NewTicker(interval)

	defer PrintPanicStack()
	defer ticker.Stop()

	checkStartIndex := int32(0)
	maxIndex := _manager._maxSessionCount

	for {
		_ = <-ticker.C

		if NetLib_IsRunningServer() == false {
			IsWantedTermination = true
			break
		}

		nextStartIndex := _checkSessionState(checkStartIndex, maxIndex, config)
		checkStartIndex = nextStartIndex
	}

	return IsWantedTermination
}

func _checkSessionState(startIndex int32, maxIndex int32, config CheckSessionStateConfig) (nextIndex int32) {
	curTime := NetLib_GetCurrnetUnixTime()

	if startIndex < 0 || startIndex >= maxIndex {
		startIndex = 0
	}

	checkMaxIndex := startIndex + int32(config.CheckCountAtOnce)
	if checkMaxIndex > maxIndex {
		checkMaxIndex = maxIndex
	}

	for i := startIndex; i < checkMaxIndex; i++ {
		if isConnectedSession(i) == false {
			continue
		}

		if _checkHeartBeat(i, int64(config.PingWaitTimeSec), curTime) {
			continue
		}

		if IsLoginUser(i) {
			if _checkDisConnectWait(i, int64(config.DisConnectWaitTimeSec), curTime) {
				continue
			}
		} else {
			if _checkNotLogIn(i, int64(config.LoginWaitTimeSec), curTime) {
				continue
			}
		}
	}

	nextIndex = checkMaxIndex
	return
}

func _checkDisConnectWait(sessionIndex int32, disConnectWaitTimeSec int64, curTime int64) bool {
	disConnectTime := GetDisConnectWaitStartTimeSec(sessionIndex)
	maxWaitTimeSec := disConnectTime + disConnectWaitTimeSec

	if disConnectTime !=0 && maxWaitTimeSec <= curTime {
		_forceDisconnect(sessionIndex, "Waiting LogOut")
		return true
	}

	return false
}

func _checkNotLogIn(sessionIndex int32,
	loginWaitTimeSec int64,
	curTime int64,
	) bool {
	connectedTime := GetConnectTimeSec(sessionIndex)
	disConnectTime := GetDisConnectWaitStartTimeSec(sessionIndex)
	maxWaitTimeSec := connectedTime + loginWaitTimeSec

	if disConnectTime == 0 && maxWaitTimeSec <= curTime {
		_disablePacketProcess(sessionIndex, protocol.ERROR_CODE_DISCONNECT_UNAUTHENTICATED_USER, "Not Login")
		return true
	}

	return false
}

func _checkHeartBeat(sessionIndex int32,
	pingWaitTimeSec int64,
	curTime int64,
	) bool {
	pingTime := GetRecentlyReceivedTimeSec(sessionIndex)
	maxWaitTimeSec := pingTime + pingWaitTimeSec

	if maxWaitTimeSec <= curTime {
		_forceDisconnect(sessionIndex, "HeartBeat")
		return true
	}

	return false
}

func _forceDisconnect(sessionIndex int32, reason string) {
	sessionUniqueId := GetNetworkUniqueID(sessionIndex)
	NTELIB_LOG_INFO("ForceDisconnectClient", zap.String("Reason", reason), zap.Int32("sessionIndex", sessionIndex), zap.Uint64("sessionUniqueId", sessionUniqueId))

	NetLibForceDisconnectClient(sessionIndex, sessionUniqueId)
}

func _disablePacketProcess(sessionIndex int32,
	errorCode int16,
	reason string,
	) {
	sessionUniqueId := GetNetworkUniqueID(sessionIndex)
	NTELIB_LOG_INFO("DisablePacketProcessClient", zap.String("Reason", reason), zap.Int32("sessionIndex", sessionIndex), zap.Uint64("sessionUniqueId", sessionUniqueId))


	protocol.NotifyErrorPacket(sessionIndex, sessionUniqueId, errorCode)

	SetDisConnectWaitStartTimeSec(sessionIndex, NetLib_GetCurrnetUnixTime())

	//NetLibDisablePacketProcessClient(sessionIndex, sessionUniqueId)
}

