package main

import (
	"strconv"
	"strings"

	. "gohipernetFake"

	"main/connectedSessions"
	"main/protocol"
	"main/roomPkg"
)

type configAppServer struct {
	GameName                   string

	RoomMaxCount               int32
	RoomStartNum               int32
	RoomMaxUserCount           int32
}

type ChatServer struct {
	ServerIndex int
	IP          string
	Port        int

	PacketChan	chan protocol.Packet

	RoomMgr *roomPkg.RoomManager
}

func createAnsStartServer(netConfig NetworkConfig, appConfig configAppServer) {
	OutPutLog(LOG_LEVEL_INFO, "", 0, "CreateServer !!!")

	var server ChatServer

	if server.setIPAddress(netConfig.BindAddress) == false {
		OutPutLog(LOG_LEVEL_INFO, "", 0, "fail. server address")
		return
	}

	protocol.Init_packet();

	maxUserCount := appConfig.RoomMaxCount * appConfig.RoomMaxUserCount
	connectedSessions.Init(netConfig.MaxSessionCount, maxUserCount)

	server.PacketChan = make(chan protocol.Packet, 256)

	roomConfig := roomPkg.RoomConfig{
		appConfig.RoomStartNum,
		appConfig.RoomMaxCount,
		appConfig.RoomMaxUserCount,
	}
	server.RoomMgr = roomPkg.NewRoomManager(roomConfig)


	go server.PacketProcess_goroutine()


	networkFunctor := SessionNetworkFunctors{}
	networkFunctor.OnConnect = server.OnConnect
	networkFunctor.OnReceive = server.OnReceive
	networkFunctor.OnReceiveBufferedData = nil
	networkFunctor.OnClose = server.OnClose
	networkFunctor.PacketTotalSizeFunc = PacketTotalSize
	networkFunctor.PacketHeaderSize = PACKET_HEADER_SIZE
	networkFunctor.IsClientSession = true


	NetLibStartNetwork(&netConfig, networkFunctor)
}

func (server *ChatServer) setIPAddress(ipAddress string) bool {
	results := strings.Split(ipAddress, ":")
	if len(results) != 2 {
		return false
	}

	server.IP = results[0]
	server.Port, _ = strconv.Atoi(results[1])

	return true
}




func (server *ChatServer) OnConnect(sessionIndex int32, sessionUniqueID uint64) {
	connectedSessions.AddSession(sessionIndex, sessionUniqueID)
}

func (server *ChatServer) OnReceive(sessionIndex int32, sessionUniqueID uint64, data []byte) bool {
	server.DistributePacket(sessionIndex, sessionUniqueID, data)
	return true
}

func (server *ChatServer) OnClose(sessionIndex int32, sessionUniqueID uint64) {
	server.disConnectClient(sessionIndex, sessionUniqueID)
}

func (server *ChatServer) disConnectClient(sessionIndex int32, sessionUniqueId uint64) {
	// 로그인도 안한 유저라면 그냥 여기서 처리한다.
	// 방 입장을 안한 유저라면 여기서 처리해도 괜찮지만 아래로 넘긴다.
	if connectedSessions.IsLoginUser(sessionIndex) == false {
		connectedSessions.RemoveSession(sessionIndex, false)
		return
	}


	packet := protocol.Packet {
		sessionIndex,
		sessionUniqueId,
		protocol.PACKET_ID_SESSION_CLOSE_SYS,
		0,
		nil,
	}

	server.PacketChan <- packet
}
