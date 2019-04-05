package gohipernetFake

import (
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)


type tcpClientSessionManager struct {
	_networkFunctor SessionNetworkFunctors

	_sessionMap sync.Map
	_curSessionCount int32 // 멀티스레드에서 호출된다
}

func newClientSessionManager(config *NetworkConfig,
							networkFunctor SessionNetworkFunctors) *tcpClientSessionManager {
	sessionMgr := new(tcpClientSessionManager)
	sessionMgr._networkFunctor = networkFunctor
	sessionMgr._sessionMap = sync.Map{}
	return sessionMgr
}

func (sessionMgr *tcpClientSessionManager) addSession(session* TcpSession) bool {
	sessionIndex := session.Index
	sessionUniqueId := session.SeqIndex

	session, resut := sessionMgr._findSession(sessionIndex, sessionUniqueId)
	if resut  {
		return false
	}

	Logger.Info("SessionManager- addSession", zap.Uint64("unique", sessionUniqueId))
	sessionMgr._sessionMap.Store(sessionUniqueId, session)
	return true
}

func (sessionMgr *tcpClientSessionManager) removeSession(sessionUniqueId uint64) {
	Logger.Info("SessionManager- removeSession", zap.Uint64("unique", sessionUniqueId))
	sessionMgr._sessionMap.Delete(sessionUniqueId)
}

func (sessionMgr *tcpClientSessionManager) sendPacket(sessionIndex int32,
			sessionUniqueId uint64,
			sendData []byte) bool {
	session, resut := sessionMgr._findSession(sessionIndex, sessionUniqueId)
	if resut == false {
		return false
	}

	session.sendPacket(sendData)
	return true
}

func (sessionMgr *tcpClientSessionManager) sendPacketAllClient(sendData []byte) {
	sessionMgr._sessionMap.Range(func(_, value interface{}) bool {
		value.(*TcpSession).sendPacket(sendData)
		return true
	})
}

func (sessionMgr *tcpClientSessionManager) _connectedessionCount() int32 {
	count := atomic.LoadInt32(&sessionMgr._curSessionCount)
	return count
}

func (sessionMgr *tcpClientSessionManager) _IncConnectedessionCount() {
	atomic.AddInt32(&sessionMgr._curSessionCount, 1)
}

func (sessionMgr *tcpClientSessionManager) _DecConnectedessionCount() {
	atomic.AddInt32(&sessionMgr._curSessionCount, -1)
}

func (sessionMgr *tcpClientSessionManager) _findSession(sessionIndex int32,
							sessionUniqueId uint64,
							) (*TcpSession, bool) {
	if session, ok := sessionMgr._sessionMap.Load(sessionUniqueId); ok {
		return session.(*TcpSession), true
	}

	return nil, false
}

func (sessionMgr *tcpClientSessionManager) forceDisconnectClient(sessionIndex int32,
	sessionUniqueId uint64) bool {
	Logger.Info("forceDisconnectClient")

	session, resut := sessionMgr._findSession(sessionIndex, sessionUniqueId)
	if resut == false {
		return false
	}

	session.closeProcess()
	return true
}

func (sessionMgr *tcpClientSessionManager) _forceCloseAllSession() {
	Logger.Info("_forceCloseAllSession - Start")

	sessionMgr._sessionMap.Range(func(_, value interface{}) bool {
		value.(*TcpSession).closeProcess()
		return true
	})

	Logger.Info("_forceCloseAllSession - End")
	IExportLog("Info", "_forceCloseAllSession - End")
}

