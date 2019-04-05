package gohipernetFake

import (
	"go.uber.org/zap"
	"net"
)



type TcpSession struct {
	Index          int32
	SeqIndex       uint64
	conn           net.Conn
	NetworkFunctor SessionNetworkFunctors
}

func (session *TcpSession) handleTcpRead(networkFunctor SessionNetworkFunctors) {
	session.NetworkFunctor.OnConnect(session.Index, session.SeqIndex)


	var startRecvPos int16
	var result int
	recviveBuff := make([]byte, MAX_RECEIVE_BUFFER_SIZE)

	for {
		recvBytes, err := session.conn.Read(recviveBuff[startRecvPos:])
		if err != nil {
			//TODO 끊는 이유 남기기
			session.closeProcess()
			return
		}

		if recvBytes < PACKET_HEADER_SIZE {
			//TODO 끊는 이유 남기기
			session.closeProcess()
			return
		}

		readAbleByte := int16(startRecvPos) + int16(recvBytes)
		startRecvPos, result = session.makePacket(readAbleByte, recviveBuff)
		if result != NET_ERROR_NONE {
			//TODO 끊는 이유 남기기
			session.closeProcess()
			return
		}

	}
}

func (session *TcpSession) makePacket(readAbleByte int16, recviveBuff []byte) (int16, int) {
	sessionIndex := session.Index
	sessionUnique := session.SeqIndex

	var startRecvPos int16 = 0
	var readPos int16

	for {
		if readAbleByte < PACKET_HEADER_SIZE {
			break
		}

		requireDataSize := packetTotalSize(recviveBuff[readPos:])

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

func (session *TcpSession) closeProcess() {
	Logger.Info("closeProcess", zap.Int32("sessionIndex", session.Index), zap.Uint64("SeqIndex", session.SeqIndex))

	session.conn.Close()
	session.NetworkFunctor.OnClose(session.Index, session.SeqIndex)

	_tcpSessionManager.removeSession(session.SeqIndex)
}

// Send bytes to client
func (session *TcpSession) sendPacket(b []byte) error {
	_, err := session.conn.Write(b)
	return err
}

func (session *TcpSession) close() error {
	return session.conn.Close()
}