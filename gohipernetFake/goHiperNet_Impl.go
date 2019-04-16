package gohipernetFake

import (
	"log"
	"net"
	//"os"
	"sync/atomic"

	//"go.uber.org/zap"

)


func init_Network_Impl(clientHeaderSize int16, serverHeaderSize int16) {
	defer PrintPanicStack()

	_InitNetworkSendFunction()
}


func start_Network_Impl(clientConfig *NetworkConfig, networkFunctor SessionNetworkFunctors) {
	defer PrintPanicStack()

	// 아래 함수가 호출되면 무한 대기에 들어간다
	_tcpSessionManager = newClientSessionManager(clientConfig, networkFunctor)
	_start_TCPServer_block(clientConfig, networkFunctor)
}

// 대기하다가 채널에 통지가 오면 listen 소켓을 종료한다
func _stopTCPServerAccept_goroutine(onDone <-chan struct{}) {
	Logger.Info("Stop TCPServer Accept")
	IExportLog("INFO", "Stop TCPServer Accept")
}

func _start_TCPServer_block(config *NetworkConfig, networkFunctor SessionNetworkFunctors) {
	defer PrintPanicStack()
	Logger.Info("tcpServerStart - Start")
	IExportLog("Info", "tcpServerStart - Start")

	var err error
	tcpAddr, _ := net.ResolveTCPAddr("tcp", config.BindAddress)
	_mClientListener, err = net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer _mClientListener.Close()

	log.Println("Server Listen ...")

	for {
		conn, _ := _mClientListener.Accept()
		client := &TcpSession{
			SeqIndex:       SeqNumIncrement(),
			TcpConn:        conn,
			NetworkFunctor: networkFunctor,
		}

		_tcpSessionManager.addSession(client)

		go client.handleTcpRead(networkFunctor)
	}

	Logger.Info("tcpServerStart - End")
	IExportLog("Info", "tcpServerStart - End")
}

// 보내기 함수(선언만 있는. 일종의 인터페이스)에 실제 동작함수를 연결한다
func _InitNetworkSendFunction() {
	NetLibISendToClient = sendToClient
	NetLibISendToAllClient = sendToAllClient
	NetLibIPostSendToAllClient = postSendToAllClient
	NetLibIPostSendToClient = postSendToClient

	Logger.Info("call _InitNetworkSendFunction")
}

func sendToClient(sessionIndex int32, sessionUniqueID uint64, data []byte) bool {
	result := _tcpSessionManager.sendPacket(sessionIndex, sessionUniqueID, data)
	return result
}

func sendToAllClient(sendData []byte) {
	_tcpSessionManager.sendPacketAllClient(sendData)
}

func postSendToClient(sessionIndex int32, sessionUniqueID uint64, data []byte) bool {
	return sendToClient(sessionIndex, sessionUniqueID, data)
}

func postSendToAllClient(sendData []byte) {
	_tcpSessionManager.sendPacketAllClient(sendData)
}

func sendPacketToServer(sessionIndex int32, data []byte) bool {
	return false
}

func postSendPacketToServer(sessionIndex int32, data []byte) bool {
	return false
}



var _seqNumber uint64 // 절대 이것을 바로 사용하면 안 된다!!!

func SeqNumIncrement() uint64 {
	newValue := atomic.AddUint64(&_seqNumber, 1)
	return newValue
}

var _tcpSessionManager *tcpClientSessionManager
var _mClientListener *net.TCPListener
