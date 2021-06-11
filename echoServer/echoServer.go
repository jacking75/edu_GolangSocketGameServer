package main

import (
	"fmt"
	. "gohipernetFake"
	"strconv"
	"strings"
)


type EchoServer struct {
	ServerIndex int
	IP          string
	Port        int
}

func createServer(netConfig NetworkConfig) {
	OutPutLog(LOG_LEVEL_INFO,"", 0,"CreateServer !!!")

	var server EchoServer

	if server.setIPAddress(netConfig.BindAddress) == false {
		OutPutLog(LOG_LEVEL_ERROR,"", 0,"fail. server address")
		return
	}


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

func (server *EchoServer) setIPAddress(ipAddress string) bool {
	results := strings.Split(ipAddress, ":")
	if len(results) != 2 {
		return false
	}

	server.IP = results[0]
	server.Port, _ = strconv.Atoi(results[1])

	return true
}

func (server *EchoServer) OnConnect(sessionIndex int32, sessionUniqueID uint64) {
	OutPutLog(LOG_LEVEL_INFO,"", 0,fmt.Sprintf("[OnConnect] sessionIndex: %d", sessionIndex))
}

func (server *EchoServer) OnReceive(sessionIndex int32, sessionUniqueID uint64, data []byte) bool {
	NetLibISendToClient(sessionIndex, sessionUniqueID, data)
	return true
}

func (server *EchoServer) OnClose(sessionIndex int32, sessionUniqueID uint64) {
	OutPutLog(LOG_LEVEL_INFO,"", 0,fmt.Sprintf("[OnClose] sessionIndex: %d", sessionIndex))
}


