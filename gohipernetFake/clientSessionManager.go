package gohipernetFake

import (
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

	session.Index = sessionIndex
	sessionMgr._sessionMap.Store(sessionUniqueId, session)
	return true
}

func (sessionMgr *tcpClientSessionManager) removeSession(sessionIndex int32, sessionUniqueId uint64) {
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

	// 쓰기 실패(예: 상대가 이미 연결을 끊은 broken pipe)를 무시하면 죽은 세션이
	// 정리되지 않고 방치된다. 실패 시 해당 세션을 명시적으로 정리한다.
	if err := session.sendPacket(sendData); err != nil {
		session.closeProcess()
		return false
	}
	return true
}

func (sessionMgr *tcpClientSessionManager) sendPacketAllClient(sendData []byte) {
	sessionMgr._sessionMap.Range(func(_, value interface{}) bool {
		session := value.(*TcpSession)
		if err := session.sendPacket(sendData); err != nil {
			session.closeProcess()
		}
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

	session, resut := sessionMgr._findSession(sessionIndex, sessionUniqueId)
	if resut == false {
		return false
	}

	session.closeProcess()
	return true
}

func (sessionMgr *tcpClientSessionManager) _forceCloseAllSession() {
	sessionMgr._sessionMap.Range(func(_, value interface{}) bool {
		value.(*TcpSession).closeProcess()
		return true
	})
}
