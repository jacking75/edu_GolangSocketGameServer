package main

import (
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	. "golang_socketGameServer_codelab/gohipernetFake"
)


type EchoServer struct {
	ServerIndex int
	IP          string
	Port        int
}

func createServer(netConfig NetworkConfig) {
	NTELIB_LOG_INFO("CreateServer !!!")

	var server EchoServer

	if server.setIPAddress(netConfig.BindAddress) == false {
		NTELIB_LOG_ERROR("fail. server address")
		return
	}


	networkFunctor := SessionNetworkFunctors{}
	networkFunctor.OnConnect = server.OnConnect
	networkFunctor.OnReceive = server.OnReceive
	networkFunctor.OnReceiveBufferedData = nil
	networkFunctor.OnClose = server.OnClose
	networkFunctor.PacketTotalSizeFunc = nil
	networkFunctor.PacketHeaderSize = PACKET_HEADER_SIZE
	networkFunctor.IsClientSession = true

	NetLibInitNetwork(PACKET_HEADER_SIZE, PACKET_HEADER_SIZE)

	NetLibStartNetwork(&netConfig, networkFunctor)

	server.Stop()
}

func (server *EchoServer) setIPAddress(ipAddress string) bool {
	results := strings.Split(ipAddress, ":")
	if len(results) != 2 {
		return false
	}

	server.IP = results[0]
	server.Port, _ = strconv.Atoi(results[1])

	NTELIB_LOG_INFO("Server Address", zap.String("IP", server.IP), zap.Int("Port", server.Port))
	return true
}

func (server *EchoServer) Stop() {
	NTELIB_LOG_INFO("chatServer Stop !!!")

	NetLib_StopServer() // 이 함수가 꼭 제일 먼저 호출 되어야 한다.

	NTELIB_LOG_INFO("chatServer Stop Waitting...")
	time.Sleep(1 * time.Second)
}


func (server *EchoServer) OnConnect(sessionIndex int32, sessionUniqueID uint64) {
	NTELIB_LOG_INFO("client OnConnect", zap.Int32("sessionIndex",sessionIndex), zap.Uint64("sessionUniqueId",sessionUniqueID))
}

func (server *EchoServer) OnReceive(sessionIndex int32, sessionUniqueID uint64, data []byte) bool {
	NTELIB_LOG_DEBUG("OnReceive", zap.Int32("sessionIndex", sessionIndex),
					zap.Uint64("sessionUniqueID", sessionUniqueID),
					zap.Int("packetSize", len(data)))

	NetLibISendToClient(sessionIndex, sessionUniqueID, data)
	return true
}

func (server *EchoServer) OnClose(sessionIndex int32, sessionUniqueID uint64) {
	NTELIB_LOG_INFO("client OnCloseClientSession", zap.Int32("sessionIndex", sessionIndex), zap.Uint64("sessionUniqueId", sessionUniqueID))
}


