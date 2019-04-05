// 클라이언트 세션에 대한 네트워크 이벤트 처리
package main

import (
	"go.uber.org/zap"

	"golang_socketGameServer_codelab/chatServer2/connectedSessions"
	. "golang_socketGameServer_codelab/gohipernetFake"
)

func (server *ChatServer) OnConnect(sessionIndex int32, sessionUniqueID uint64) {
	NTELIB_LOG_INFO("client OnConnect", zap.Int32("sessionIndex",sessionIndex), zap.Uint64("sessionUniqueId",sessionUniqueID))

	connectedSessions.AddSession(sessionIndex, sessionUniqueID)
}

func (server *ChatServer) OnReceive(sessionIndex int32, sessionUniqueID uint64, data []byte) bool {
	server._distributePacket(sessionIndex, sessionUniqueID, data)
	return true
}

func (server *ChatServer) OnClose(sessionIndex int32, sessionUniqueID uint64) {
	NTELIB_LOG_INFO("client OnCloseClientSession", zap.Int32("sessionIndex", sessionIndex), zap.Uint64("sessionUniqueId", sessionUniqueID))

	server.disConnectClient(sessionIndex, sessionUniqueID)
}



func (server *ChatServer) disConnectClient(sessionIndex int32, sessionUniqueId uint64) {
	// 로그인도 안한 유저라면 그냥 여기서 처리한다.
	if connectedSessions.IsLoginUser(sessionIndex) == false {
		NTELIB_LOG_INFO("DisConnectClient - Not Login User", zap.Int32("sessionIndex", sessionIndex))
		connectedSessions.RemoveSession(sessionIndex, false)
		return
	}

	connectedTime := connectedSessions.GetConnectTimeSec(sessionIndex)
	if connectedTime == 0 {
		NTELIB_LOG_INFO("DisConnectClient - getConnectTimeSec is 0!")
	} else {
		//DB.WriteUserLastLoginTime(connectedTime, userID)
	}


	roomNum, roomNumOfEntering := connectedSessions.GetRoomNumber(sessionIndex)

	connectedSessions.RemoveSession(sessionIndex, true)

	if roomNum > -1 {
		server._roomMgr.DisConnectedUser(sessionUniqueId, roomNum)
	}

	// 현재 방에 들어가고 있는 중이므로 방에서 유저를 뺀다.
	if roomNumOfEntering > -1 && roomNum != roomNumOfEntering {
		server._roomMgr.DisConnectedUser(sessionUniqueId, roomNumOfEntering)
	}

	// 해당 세션을 초기화 한다. 메모리를 지우고, redis도 지운다
	NTELIB_LOG_INFO("DisConnectClient - Login User", zap.Int32("sessionIndex", sessionIndex))
}