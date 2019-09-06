// 내부에서 사용하는 패킷
package protocol

import (
	. "gohipernetFake"
)

type InternalPacket struct {
	RoomIndex   int32
	Id          int16
	DataSize    int16
	Data        []byte
}


type InternalPacketDisConnectedUserToRoom struct{
	SessionUniqueId uint64
	RoomNum int32
}

func (packet *InternalPacketDisConnectedUserToRoom ) Encoding() ([]byte, int16) {
	totalSize := int16(8 + 4)
	sendBuf := make([]byte, totalSize)

	writer := MakeWriter(sendBuf, true)
	writer.WriteU64(packet.SessionUniqueId)
	writer.WriteS32(packet.RoomNum)
	return sendBuf, totalSize
}

func (packet *InternalPacketDisConnectedUserToRoom) Decoding(Data []byte) bool {
	reader := MakeReader(Data, true)

	packet.SessionUniqueId, _ = reader.ReadU64()
	packet.RoomNum, _ = reader.ReadS32()
	return true
}