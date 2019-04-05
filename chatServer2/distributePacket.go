package main

import (
	"bytes"

	"go.uber.org/zap"

	"golang_socketGameServer_codelab/chatServer2/connectedSessions"
	"golang_socketGameServer_codelab/chatServer2/protocol"
	"golang_socketGameServer_codelab/chatServer2/roomPkg"
	. "golang_socketGameServer_codelab/gohipernetFake"
)

func (server *ChatServer) _distributePacket(sessionIndex int32,
	sessionUniqueId uint64,
	packetData []byte,
	) {
	packetID := protocol.PeekPacketID(packetData)
	bodySize, bodyData := protocol.PeekPacketBody(packetData)
	NTELIB_LOG_DEBUG("_distributePacket", zap.Int32("sessionIndex", sessionIndex), zap.Uint64("sessionUniqueId", sessionUniqueId), zap.Int16("PacketID", packetID))

	// 여기에서 바로 처리하는 패킷
	if server.ProcessPacket(sessionIndex, sessionUniqueId, packetID, bodySize, bodyData) {
		return
	}

	if connectedSessions.IsLoginUser(sessionIndex) == false {
		protocol.NotifyErrorPacket(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_PACKET_NOT_LOGIN_USER)
		return
	}


	// 룸 번호를 얻는다.
	userRoomNum, errorCode := _rommNumberFromPacketOrSessionInfo(packetID, bodyData, sessionIndex, sessionUniqueId )
	if errorCode != protocol.ERROR_CODE_NONE {
		return;
	}

	server._packetToRoom(packetID, bodySize, bodyData, sessionIndex, sessionUniqueId, userRoomNum)
}

func (server *ChatServer) ProcessPacket(sessionIndex int32,
	sessionUniqueId uint64,
	packetID int16,
	bodySize int16,
	bodyData []byte,
	) bool {
	isProcessed := false

	switch packetID {
	case protocol.PACKET_ID_LOGIN_REQ:
		{
			isProcessed = true

			//DB와 연동하지 않으므로 중복 로그인만 아니면 다 성공으로 한다
			var request protocol.LoginReqPacket
			if (&request).Decoding(bodyData) == false {
				_sendLoginResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_PACKET_DECODING_FAIL)
				break
			}

			userID := bytes.Trim(request.UserID[:], "\x00");
			//gohipernet.LOG_DEBUG("PACKET_ID_LOGIN_REQ", zap.String("userID", userID))

			if len(userID) <= 0 {
				_sendLoginResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_LOGIN_USER_INVALID_ID)
				break
			}

			curTime := NetLib_GetCurrnetUnixTime()

			if connectedSessions.SetLogin(sessionIndex, sessionUniqueId, userID, curTime) == false {
				_sendLoginResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_LOGIN_USER_DUPLICATION)
				break
			}

			_sendLoginResult(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_NONE)
		}

	case protocol.PACKET_ID_DEV_ECHO_REQ:
		{

		}
	}

	return isProcessed
}

func _sendLoginResult(sessionIndex int32, sessionUniqueId uint64, result int16) {
	var response protocol.LoginResPacket
	response.Result = result
	sendPacket, _ := response.EncodingPacket()

	NetLibIPostSendToClient(sessionIndex, sessionUniqueId, sendPacket)

	NTELIB_LOG_DEBUG("SendLoginResult", zap.Int32("sessionIndex", sessionIndex), zap.Int16("result", result))
}

func _rommNumberFromPacketOrSessionInfo(packetID int16,
	bodyData []byte,
	sessionIndex int32,
	sessionUniqueId uint64,
	) (int32, int16) {
	userRoomNum, userRoomNumOfEntering := connectedSessions.GetRoomNumber(sessionIndex)

	if packetID == protocol.PACKET_ID_ROOM_ENTER_REQ {
		// 이미 룸에 있다면 에러
		if userRoomNum >= 0 || userRoomNumOfEntering >= 0 {
			protocol.NotifyErrorPacket(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_ENTER_ROOM_ALREADY)
			return -1, protocol.ERROR_CODE_ENTER_ROOM_ALREADY
		}

		roomEnterReq := protocol.RoomEnterReqPacket{}
		(&roomEnterReq).Decoding(bodyData)

		if roomEnterReq.RoomNumber >= 0 {
			userRoomNum = roomEnterReq.RoomNumber
		} else {
			//TODO: 자동을 원하면 적당한 곳에 넣는다.
			NTELIB_LOG_ERROR("_distributePacket - Room.  룸입장 자동 번호 할당은 미 구현", zap.Int32("sessionIndex", sessionIndex))
			return -1, protocol.ERROR_CODE_ENTER_ROOM_AUTO_ROOM_NUMBER
		}

		if connectedSessions.SetRoomEntering(sessionIndex, userRoomNum) == false {
			protocol.NotifyErrorPacket(sessionIndex, sessionUniqueId, protocol.ERROR_CODE_ENTER_ROOM_PREV_WORKING)
			return -1, protocol.ERROR_CODE_ENTER_ROOM_PREV_WORKING
		}
	}

	if userRoomNum < 0 {
		return -1, protocol.ERROR_CODE_USER_NOT_IN_ROOM
	}

	return userRoomNum, protocol.ERROR_CODE_NONE
}

func (server *ChatServer) _packetToRoom(packetID int16,
	bodySize int16,
	bodyData []byte,
	sessionIndex int32,
	sessionUniqueId uint64,
	userRoomNum int32,
	) {
	packet := protocol.Packet{Id: packetID}
	packet.UserSessionIndex = sessionIndex
	packet.UserSessionUniqueId = sessionUniqueId
	packet.RoomNumber = userRoomNum
	packet.Id = packetID
	packet.DataSize = bodySize
	packet.Data = make([]byte, packet.DataSize)
	copy(packet.Data, bodyData)

	if packetID == protocol.PACKET_ID_ROOM_ENTER_REQ {
		userID, _ := connectedSessions.GetUserID(sessionIndex)
		userIDLen := connectedSessions.GetUserIDLength(sessionIndex)
		packet.UserID = make([]byte, userIDLen)
		copy(packet.UserID, userID)
	}

	startRoomNumber := int32(server.appConfig.RoomStartNum)
	roomIndex := roomPkg.RoomNumberToIndex(startRoomNumber, userRoomNum)
	server._roomMgr.PushPacket(roomIndex, packet)

	NTELIB_LOG_DEBUG("_distributePacket - Room", zap.Int32("sessionIndex", sessionIndex), zap.Int32("RoomNumber", userRoomNum), zap.Int32("RoomIndex", roomIndex))
}