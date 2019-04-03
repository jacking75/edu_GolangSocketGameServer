// 애플리케이션에서 네트워크 라이브러리에 접근할 함수는 모두 여기에만 정의한다.
package gohipernetFake



func NetLibInitLog() {
	init_Log()
	wrapLoggerFunc()
}
// 네트워크 초기화
func NetLibInitNetwork(clientHeaderSize int16, serverHeaderSize int16) {
	init_Network_Impl(clientHeaderSize, serverHeaderSize)
}

// 네트워크 시작
func NetLibStartNetwork(clientConfig *NetworkConfig, networkFunctor SessionNetworkFunctors) {
	start_Network_Impl(clientConfig, networkFunctor)
}


// Send Interface Function
var NetLibISendToClient func(int32, uint64, []byte) bool
var NetLibISendToAllClient func([]byte)
var NetLibIPostSendToAllClient func([]byte)
var NetLibIPostSendToClient func(int32, uint64, []byte) bool

// 지정한 클라이언트를 강제 종료 시킨다
func NetLibForceDisconnectClient(sessionIndex int32, sessionUnqiueID uint64) {
	//_tcpSessionManager.forceDisconnectClient(sessionIndex, sessionUnqiueID)
}