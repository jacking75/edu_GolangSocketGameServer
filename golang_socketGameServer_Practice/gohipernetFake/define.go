package gohipernetFake


const (
	MAX_RECEIVE_BUFFER_SIZE = 8012
	PACKET_HEADER_SIZE      = 5
	MAX_PACKET_SIZE         = 1024
)

const (
	NET_ERROR_NONE = 0
	NET_ERROR_RECV_MAKE_PACKET_TOO_LARGE_PACKET_SIZE = 1

)
const (
	NET_CLOSE_REMOTE = 1
	NET_CLOSE_RECV_TOO_SMALL_RECV_DATA = 2
)



type SessionNetworkFunctors struct {
	OnConnect func(int32, uint64)

	OnClose func(int32, uint64)

	// 데이터 도착 이벤트
	OnReceive func(int32, uint64, []byte) bool

	// 데이터 도착 이벤트. []byte가 링버퍼에 저장되어 있다
	OnReceiveBufferedData func(int32, uint64, []byte) bool


	// 데이터를 분석하여 패킷 크기를 반환한다.
	PacketTotalSizeFunc func([]byte) int16

	// 패킷 헤더의 크기
	PacketHeaderSize int16

	// true 이면 client와 연결한 세션이다.
	IsClientSession bool
}