package gohipernetFake

import (
	"net"
	"sync"
)



type TcpSession struct {
	Index          int32
	SeqIndex       uint64
	TcpConn        net.Conn
	NetworkFunctor SessionNetworkFunctors

	closeOnce sync.Once
}

func (session *TcpSession) handleTcpRead(networkFunctor SessionNetworkFunctors) {
	// 이 고루틴 안에서 패닉이 나도(예: 조작된 패킷) 이 세션만 정리하고 다른 접속자에게는 영향이 없도록 한다.
	defer PrintPanicStack()
	defer session.closeProcess()

	session.NetworkFunctor.OnConnect(session.Index, session.SeqIndex)


	var startRecvPos int16
	var result int
	recviveBuff := make([]byte, MAX_RECEIVE_BUFFER_SIZE)

	for {
		recvBytes, err := session.TcpConn.Read(recviveBuff[startRecvPos:])
		if err != nil {
			//TODO 끊는 이유 남기기
			return
		}

		if recvBytes < PACKET_HEADER_SIZE {
			//TODO 끊는 이유 남기기
			return
		}

		readAbleByte := int16(startRecvPos) + int16(recvBytes)
		startRecvPos, result = session.makePacket(readAbleByte, recviveBuff)
		if result != NET_ERROR_NONE {
			//TODO 끊는 이유 남기기
			return
		}

	}
}

func (session *TcpSession) makePacket(readAbleByte int16, recviveBuff []byte) (int16, int) {
	sessionIndex := session.Index
	sessionUnique := session.SeqIndex

	PacketHeaderSize := session.NetworkFunctor.PacketHeaderSize
	PacketTotalSizeFunc := session.NetworkFunctor.PacketTotalSizeFunc
	var startRecvPos int16 = 0
	var readPos int16

	for {
		if readAbleByte < PacketHeaderSize {
			break
		}

		requireDataSize := PacketTotalSizeFunc(recviveBuff[readPos:])

		// 클라이언트가 보낸 패킷 길이가 헤더 크기보다 작으면(0 포함) 유효하지 않은 패킷이다.
		// uint16 길이 필드가 int16으로 캐스팅되며 음수로 뒤집힌 경우도 여기서 함께 걸러진다.
		// 검증 없이 진행하면 아래 readAbleByte/readPos가 줄지 않아 무한 루프에 빠지거나,
		// 음수 상한으로 슬라이싱하여 패닉이 발생한다.
		if requireDataSize < PacketHeaderSize {
			return startRecvPos, NET_ERROR_RECV_MAKE_PACKET_INVALID_PACKET_SIZE
		}

		if requireDataSize > readAbleByte {
			break
		}

		if requireDataSize > MAX_PACKET_SIZE {
			return startRecvPos, NET_ERROR_RECV_MAKE_PACKET_TOO_LARGE_PACKET_SIZE
		}

		ltvPacket := recviveBuff[readPos:(readPos + requireDataSize)]
		readPos += requireDataSize
		readAbleByte -= requireDataSize


		session.NetworkFunctor.OnReceive(sessionIndex, sessionUnique, ltvPacket)
	}


	if readAbleByte > 0 {
		copy(recviveBuff, recviveBuff[readPos:(readPos+readAbleByte)])
	}

	startRecvPos = readAbleByte
	return startRecvPos, NET_ERROR_NONE
}

// closeProcess는 읽기 루프의 에러 경로와 NetLibForceDisconnectClient 양쪽에서
// 동시에 호출될 수 있다. sync.Once로 감싸 세션 인덱스가 풀에 이중 반환(double free)되거나
// OnClose 콜백이 두 번 호출되는 것을 막는다.
func (session *TcpSession) closeProcess() {
	session.closeOnce.Do(func() {
		session.TcpConn.Close()
		session.NetworkFunctor.OnClose(session.Index, session.SeqIndex)

		_tcpSessionManager.removeSession(session.Index, session.SeqIndex)
	})
}

// Send bytes to client
func (session *TcpSession) sendPacket(b []byte) error {
	_, err := session.TcpConn.Write(b)
	return err
}

func (session *TcpSession) close() error {
	return session.TcpConn.Close()
}