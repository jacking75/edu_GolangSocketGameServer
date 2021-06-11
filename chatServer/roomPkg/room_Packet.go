package roomPkg

import (
	"main/protocol"
)

func (room *baseRoom) _packetProcess_Relay(user *roomUser, packet protocol.Packet) int16 {
	var relayNotify protocol.RoomRelayNtfPacket
	relayNotify.RoomUserUniqueId = user.RoomUniqueId
	relayNotify.Data = packet.Data
	notifySendBuf, packetSize := relayNotify.EncodingPacket(packet.DataSize)
	room.broadcastPacket(packetSize, notifySendBuf, 0)

	return protocol.ERROR_CODE_NONE
}
