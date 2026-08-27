package gohipernetFake

import (
	"encoding/binary"
	"testing"
	"time"
)

// 테스트용 PacketTotalSizeFunc: 버퍼의 첫 2바이트(리틀 엔디안)를 패킷 총 길이로 사용한다.
func testPacketTotalSizeFunc(data []byte) int16 {
	return int16(binary.LittleEndian.Uint16(data))
}

func newTestSession(onReceive func(int32, uint64, []byte) bool) *TcpSession {
	return &TcpSession{
		NetworkFunctor: SessionNetworkFunctors{
			PacketHeaderSize:    5,
			PacketTotalSizeFunc: testPacketTotalSizeFunc,
			OnReceive:           onReceive,
		},
	}
}

func TestMakePacket_SinglePacket(t *testing.T) {
	var received [][]byte
	session := newTestSession(func(_ int32, _ uint64, data []byte) bool {
		cp := append([]byte(nil), data...)
		received = append(received, cp)
		return true
	})

	buf := make([]byte, MAX_RECEIVE_BUFFER_SIZE)
	binary.LittleEndian.PutUint16(buf[0:2], 7) // totalSize = 5(header) + 2(body)
	buf[5], buf[6] = 0xAA, 0xBB

	leftover, result := session.makePacket(7, buf)
	if result != NET_ERROR_NONE {
		t.Fatalf("result = %d, want NET_ERROR_NONE", result)
	}
	if leftover != 0 {
		t.Fatalf("leftover = %d, want 0", leftover)
	}
	if len(received) != 1 || len(received[0]) != 7 {
		t.Fatalf("received = %v, want exactly one 7-byte packet", received)
	}
}

func TestMakePacket_PartialPacketWaitsForMoreData(t *testing.T) {
	var receiveCount int
	session := newTestSession(func(_ int32, _ uint64, data []byte) bool {
		receiveCount++
		return true
	})

	buf := make([]byte, MAX_RECEIVE_BUFFER_SIZE)
	binary.LittleEndian.PutUint16(buf[0:2], 7) // totalSize = 7, but only 4 bytes have arrived

	leftover, result := session.makePacket(4, buf)
	if result != NET_ERROR_NONE {
		t.Fatalf("result = %d, want NET_ERROR_NONE", result)
	}
	if leftover != 4 {
		t.Fatalf("leftover = %d, want 4 (packet not complete yet)", leftover)
	}
	if receiveCount != 0 {
		t.Fatalf("OnReceive called %d times, want 0", receiveCount)
	}
}

func TestMakePacket_TooLargePacketIsRejected(t *testing.T) {
	session := newTestSession(func(_ int32, _ uint64, data []byte) bool { return true })

	buf := make([]byte, MAX_RECEIVE_BUFFER_SIZE)
	binary.LittleEndian.PutUint16(buf[0:2], MAX_PACKET_SIZE+1)

	_, result := session.makePacket(MAX_PACKET_SIZE+1, buf)
	if result != NET_ERROR_RECV_MAKE_PACKET_TOO_LARGE_PACKET_SIZE {
		t.Fatalf("result = %d, want NET_ERROR_RECV_MAKE_PACKET_TOO_LARGE_PACKET_SIZE", result)
	}
}

// 패킷 길이 필드가 헤더 크기보다 작은 경우(0 포함), 예전 코드는 readAbleByte/readPos가 줄지 않아
// 무한 루프에 빠졌다. 이제는 즉시 에러를 반환해야 한다.
func TestMakePacket_TooSmallPacketSizeIsRejectedNotInfiniteLoop(t *testing.T) {
	session := newTestSession(func(_ int32, _ uint64, data []byte) bool { return true })

	buf := make([]byte, MAX_RECEIVE_BUFFER_SIZE)
	binary.LittleEndian.PutUint16(buf[0:2], 0) // totalSize = 0

	done := make(chan struct{})
	var result int
	go func() {
		_, result = session.makePacket(10, buf)
		close(done)
	}()

	select {
	case <-done:
		if result != NET_ERROR_RECV_MAKE_PACKET_INVALID_PACKET_SIZE {
			t.Fatalf("result = %d, want NET_ERROR_RECV_MAKE_PACKET_INVALID_PACKET_SIZE", result)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("makePacket did not return - infinite loop on a too-small packet size")
	}
}

// uint16 길이 필드가 32768 이상일 때 int16으로 캐스팅되며 음수로 뒤집히는 상황을 그대로 흉내낸다.
// 이 경우도 헤더 크기 미만으로 취급되어 안전하게 거부되어야 한다(음수 상한 슬라이싱 패닉 방지).
func TestMakePacket_NegativeSizeFromUint16OverflowIsRejected(t *testing.T) {
	session := newTestSession(func(_ int32, _ uint64, data []byte) bool { return true })

	buf := make([]byte, MAX_RECEIVE_BUFFER_SIZE)
	binary.LittleEndian.PutUint16(buf[0:2], 40000) // int16(40000) < 0

	_, result := session.makePacket(10, buf)
	if result != NET_ERROR_RECV_MAKE_PACKET_INVALID_PACKET_SIZE {
		t.Fatalf("result = %d, want NET_ERROR_RECV_MAKE_PACKET_INVALID_PACKET_SIZE", result)
	}
}
