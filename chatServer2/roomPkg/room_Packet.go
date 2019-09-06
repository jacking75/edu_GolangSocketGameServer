package roomPkg

import (
	"go.uber.org/zap"

	"main/protocol"
	. "gohipernetFake"
)

func (room *baseRoom) _packetProcess_Relay(user *roomUser, packet protocol.Packet) int16 {
	var relayNotify protocol.RoomRelayNtfPacket
	relayNotify.RoomUserUniqueId = user.RoomUniqueId
	relayNotify.Data = packet.Data
	notifySendBuf, packetSize := relayNotify.EncodingPacket(packet.DataSize)
	room.broadcastPacket(packetSize, notifySendBuf, 0)

	NTELIB_LOG_DEBUG("Room Relay", zap.String("Sender", string(user.ID[:])))
	return protocol.ERROR_CODE_NONE
}