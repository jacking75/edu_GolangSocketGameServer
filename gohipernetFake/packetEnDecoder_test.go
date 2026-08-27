package gohipernetFake

import (
	"bytes"
	"testing"
)

func TestRawPacketData_ReadWriteRoundTrip(t *testing.T) {
	buf := make([]byte, 32)
	w := MakeWriter(buf, true)
	w.WriteU16(1234)
	w.WriteS32(-999)
	w.WriteBytes([]byte("hi"))

	r := MakeReader(buf, true)

	v16, err := r.ReadU16()
	if err != nil || v16 != 1234 {
		t.Fatalf("ReadU16: got (%v, %v), want (1234, nil)", v16, err)
	}

	v32, err := r.ReadS32()
	if err != nil || v32 != -999 {
		t.Fatalf("ReadS32: got (%v, %v), want (-999, nil)", v32, err)
	}

	got := r.ReadBytes(2)
	if !bytes.Equal(got, []byte("hi")) {
		t.Fatalf("ReadBytes: got %q, want %q", got, "hi")
	}
}

// gohipernetFake/TcpSession.go의 makePacket이 조작된 패킷(예: 길이 필드가 헤더 크기보다 작음)을
// 받았을 때, ReadBytes가 범위를 벗어난 길이로 호출되어도 패닉 없이 안전하게 처리되어야 한다.
func TestRawPacketData_ReadBytes_OutOfRangeReturnsNil(t *testing.T) {
	r := MakeReader([]byte{1, 2, 3}, true)

	if got := r.ReadBytes(10); got != nil {
		t.Fatalf("ReadBytes(10) on a 3-byte buffer: got %v, want nil", got)
	}

	if got := r.ReadBytes(-1); got != nil {
		t.Fatalf("ReadBytes(-1): got %v, want nil", got)
	}
}

func TestRawPacketData_ReadU16_OutOfRangeReturnsError(t *testing.T) {
	r := MakeReader([]byte{1}, true)

	if _, err := r.ReadU16(); err == nil {
		t.Fatalf("ReadU16 on a 1-byte buffer: expected an error, got nil")
	}
}

func TestRawPacketData_ReadString_OutOfRangeReturnsError(t *testing.T) {
	// 길이 헤더(2바이트)는 5를 가리키지만 실제 데이터는 1바이트뿐인 조작된 입력.
	r := MakeReader([]byte{5, 0, 'a'}, true)

	if _, err := r.ReadString(); err == nil {
		t.Fatalf("ReadString with truncated body: expected an error, got nil")
	}
}
