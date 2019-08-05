package gohipernetFake

import (
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

type tcpClientSessionManager struct {
	_networkFunctor SessionNetworkFunctors

	_sessionMap      sync.Map
	_curSessionCount int32 // 멀티스레드에서 호출된다

	sessionIndexPool *Deque
}

func newClientSessionManager(config *NetworkConfig,
	networkFunctor SessionNetworkFunctors) *tcpClientSessionManager {
	sessionMgr := new(tcpClientSessionManager)
	sessionMgr._networkFunctor = networkFunctor
	sessionMgr._sessionMap = sync.Map{}

	sessionMgr._createSessionIndexPool(config.MaxSessionCount)

	return sessionMgr
}

func (sessionMgr *tcpClientSessionManager) _createSessionIndexPool(poolSize int) {
	sessionMgr.sessionIndexPool = NewCappedDeque(poolSize)

	for i := 0; i < poolSize; i++ {
		sessionMgr.sessionIndexPool.Append(int32(i))
	}
}

func (sessionMgr *tcpClientSessionManager) _allocSessionIndex() int32 {
	index := sessionMgr.sessionIndexPool.Shift()

	if index == nil {
		return -1
	}

	return index.(int32)
}

func (sessionMgr *tcpClientSessionManager) _freeSessionIndex(sessionIndex int32) {
	sessionMgr.sessionIndexPool.Append(sessionIndex)
}

func (sessionMgr *tcpClientSessionManager) addSession(session *TcpSession) bool {
	sessionUniqueId := session.SeqIndex
	sessionIndex := sessionMgr._allocSessionIndex()

	if sessionIndex == -1 {
		return false
	}

	_, result := sessionMgr._findSession(sessionIndex, sessionUniqueId)
	if result {
		return false
	}

	Logger.Info("SessionManager- addSession", zap.Int32("index", sessionIndex),
		zap.Uint64("unique", sessionUniqueId))

	session.Index = sessionIndex
	sessionMgr._sessionMap.Store(sessionUniqueId, session)
	return true
}

func (sessionMgr *tcpClientSessionManager) removeSession(sessionIndex int32, sessionUniqueId uint64) {
	Logger.Info("SessionManager- removeSession", zap.Int32("index", sessionIndex), zap.Uint64("unique", sessionUniqueId))
	sessionMgr._freeSessionIndex(sessionIndex)
	sessionMgr._sessionMap.Delete(sessionUniqueId)
}

func (sessionMgr *tcpClientSessionManager) sendPacket(sessionIndex int32,
	sessionUniqueId uint64,
	sendData []byte) bool {
	session, result := sessionMgr._findSession(sessionIndex, sessionUniqueId)
	if result == false {
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
