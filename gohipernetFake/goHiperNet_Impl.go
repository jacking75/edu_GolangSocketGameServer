package gohipernetFake

import (
	"log"
	"net"
	"sync/atomic"
)


func start_Network_Impl(clientConfig *NetworkConfig, networkFunctor SessionNetworkFunctors) {
	defer PrintPanicStack()

	_InitNetworkSendFunction()

	// 아래 함수가 호출되면 무한 대기에 들어간다
	_tcpSessionManager = newClientSessionManager(clientConfig, networkFunctor)
	_start_TCPServer_block(clientConfig, networkFunctor)
}

func stopListen_impl() {
	_ = _mClientListener.Close()
}

func _start_TCPServer_block(config *NetworkConfig, networkFunctor SessionNetworkFunctors) {
	defer PrintPanicStack()
	logInfo("", 0, "tcpServerStart - Start")

	var err error
	tcpAddr, _ := net.ResolveTCPAddr("tcp", config.BindAddress)
	_mClientListener, err = net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer _mClientListener.Close()

	log.Println("Server Listen ...")

	for {
		conn, err := _mClientListener.Accept()
		if err != nil {
			// NetLibStopListen() 등으로 리스너가 닫힌 경우 Accept()는 계속 에러를 반환하므로 루프를 종료한다.
			logError("", 0, "Accept error: "+err.Error())
			break
		}

		client := &TcpSession{
			SeqIndex:       SeqNumIncrement(),
			TcpConn:        conn,
			NetworkFunctor: networkFunctor,
		}

		// 세션 인덱스 풀이 소진된 경우(MaxSessionCount 초과) addSession은 false를 반환하며
		// client.Index를 세팅하지 않는다. 반환값을 확인하지 않으면 그 연결이 index 0인 것처럼
		// 동작해 다른 정상 세션의 인덱스를 나중에 훔쳐 반환하는 문제로 이어진다.
		if _tcpSessionManager.addSession(client) == false {
			logError("", 0, "addSession failed. session pool exhausted")
			_ = conn.Close()
			continue
		}

		go client.handleTcpRead(networkFunctor)
	}

	logInfo("", 0, "tcpServerStart - End")
}

// 보내기 함수(선언만 있는. 일종의 인터페이스)에 실제 동작함수를 연결한다
func _InitNetworkSendFunction() {
	NetLibISendToClient = sendToClient
	NetLibISendToAllClient = sendToAllClient
	NetLibIPostSendToAllClient = postSendToAllClient
	NetLibIPostSendToClient = postSendToClient

	logInfo("", 0, "call _InitNetworkSendFunction")
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
