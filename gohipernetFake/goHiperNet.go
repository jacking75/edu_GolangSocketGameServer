// 애플리케이션에서 네트워크 라이브러리에 접근할 함수는 모두 여기에만 정의한다.
package gohipernetFake



func NetLibInitLog(loglevel int, logFunc func(int, string, uint64, string)) {
	_logLevel = loglevel

	if logFunc != nil {
		OutPutLog = logFunc
	}
}

// 네트워크 시작
func NetLibStartNetwork(clientConfig *NetworkConfig, networkFunctor SessionNetworkFunctors) {
	start_Network_Impl(clientConfig, networkFunctor)
}

func NetLibStopListen() {
	stopListen_impl()
}

// 특정 클라이언트에게 데이터를 보낸다
var NetLibISendToClient func(int32, uint64, []byte) bool
// 접속된 모든 클라이언트에게 데이터를 보낸다.
var NetLibISendToAllClient func([]byte)

var NetLibIPostSendToAllClient func([]byte)
var NetLibIPostSendToClient func(int32, uint64, []byte) bool


// 클라이언트 접속을 강제로 짜른다.
func NetLibForceDisconnectClient(sessionIndex int32, sessionUnqiueID uint64) {
	_tcpSessionManager.forceDisconnectClient(sessionIndex, sessionUnqiueID)
}