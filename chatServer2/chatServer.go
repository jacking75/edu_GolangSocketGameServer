package main

import (
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	. "gohipernetFake"
	"main/connectedSessions"
	"main/protocol"
	"main/roomPkg"
)

type ChatServer struct {
	ServerIndex int
	IP          string
	Port        int

	appConfig configAppServer

	_roomMgr         *roomPkg.RoomManager

	//_timerScheduler *TimerScheduler //현재는 사용하지 않음
}

func createServer(netConfig NetworkConfig, appConfig configAppServer) {
	NTELIB_LOG_INFO("CreateServer !!!")

	if _checkConfig(netConfig, appConfig) == false {
		return
	}

	var server ChatServer

	if server.setIPAddress(netConfig.BindAddress) == false {
		NTELIB_LOG_ERROR("fail. server address")
		return
	}

	protocol.Init_packet();
	
	server.init(int32(netConfig.MaxSessionCount), appConfig)

	networkFunctor := SessionNetworkFunctors{}
	networkFunctor.OnConnect = server.OnConnect
	networkFunctor.OnReceive = server.OnReceive
	networkFunctor.OnReceiveBufferedData = nil
	networkFunctor.OnClose = server.OnClose
	networkFunctor.PacketTotalSizeFunc = PacketTotalSize
	networkFunctor.PacketHeaderSize = PACKET_HEADER_SIZE
	networkFunctor.IsClientSession = true

	// 실습에서는 아래 코드 호출하여도 적용 되자 않음
	NetLibInitNetwork(protocol.ClientHeaderSize(), protocol.ServerHeaderSize())

	NetLibStartNetwork(&netConfig, networkFunctor)

	server.Stop()
}

func (server *ChatServer) setIPAddress(ipAddress string) bool {
	results := strings.Split(ipAddress, ":")
	if len(results) != 2 {
		return false
	}

	server.IP = results[0]
	server.Port, _ = strconv.Atoi(results[1])

	NTELIB_LOG_INFO("Server Address", zap.String("IP", server.IP), zap.Int("Port", server.Port))
	return true
}

func (server *ChatServer) init(maxClientSessionCount int32, appConfig configAppServer) bool {
	server.appConfig = appConfig

	maxUserCount := int32((appConfig.RoomMaxCount * appConfig.RoomMaxUserCount) + 3) // 3은 관리자 수
	checkStateConfig := connectedSessions.CheckSessionStateConfig{
		appConfig.CheckCountAtOnce,
		appConfig.CheckReriodMillSec,
		appConfig.LoginWaitTimeSec,
		appConfig.DisConnectWaitTimeSec,
		appConfig.RoomEnterWaitTimeSec,
		appConfig.PingWaitTimeSec,
		appConfig.MaxRequestCountPerSecond}
	connectedSessions.Init(maxClientSessionCount, maxUserCount, checkStateConfig)
	connectedSessions.Start()

	config := roomPkg.RoomConfig{int32(appConfig.RoomStartNum),
		int32(appConfig.RoomMaxCount),
		int32(appConfig.RoomMaxUserCount),
		int32(appConfig.RoomCountByGoroutine),
		int32(appConfig.RoomMaxProcessBufferCount),
		int32(appConfig.RoomInternalPacketChanBufferCount) }

	server._roomMgr = roomPkg.NewRoomManager(config)
	server._roomMgr.Start()

	//server._init_TimerScheduler()
	return true
}

/*func (server *ChatServer) _init_TimerScheduler() {
	NTELIB_LOG_INFO("Start main TimerScheduler")

	server._timerScheduler = new(TimerScheduler)
	server._timerScheduler.Start()
}*/

func (server *ChatServer) PacketTotalSize(data []byte) int16 {
	totalsize := binary.LittleEndian.Uint16(data)
	return int16(totalsize)
}

func (server *ChatServer) Stop() {
	NTELIB_LOG_INFO("chatServer Stop !!!")

	server._roomMgr.Stop()

	NTELIB_LOG_INFO("chatServer Stop Waitting...")
	time.Sleep(1 * time.Second)
}


func _checkConfig(netConfig NetworkConfig, appConfig configAppServer) bool {
	userCount := appConfig.RoomMaxUserCount * appConfig.RoomMaxCount

	if netConfig.MaxSessionCount < userCount {
		NTELIB_LOG_ERROR("userCount less than netConfig.MaxSessionCount", zap.Int("userCount", userCount), zap.Int("MaxSessionCount", netConfig.MaxSessionCount))
		return false
	}
	return true
}
