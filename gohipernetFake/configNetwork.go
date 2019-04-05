package gohipernetFake

import (
	"go.uber.org/zap"
)


type NetworkConfig struct {
	IsTcp4Addr                    bool
	BindAddress                   string // 예) localhost:19999
	MaxSessionCount               int // 최대 클라이언트 세션 수. 넉넉하게 많이 해도 괜찮다
	MaxPacketSize                 int // 최대 패킷 크기
	MaxReceiveBufferSize          int // 받기 버퍼 크기. 최소 MaxPacketSize 2배 이상 추천.

}

func (config NetworkConfig) WriteNetworkConfig(isClientSide bool) {
	section := "ClientSide"
	if isClientSide == false {
		section = "ServerSide"
	}

	Logger.Info("config - " + section,
		zap.Bool("IsTcp4Addr", config.IsTcp4Addr),
		zap.String("ClientAddress", config.BindAddress),
		zap.Int("MaxSessionCount", config.MaxSessionCount),
		zap.Int("MaxPacketSize", config.MaxPacketSize),
		zap.Int("MaxReceiveBufferSize", config.MaxReceiveBufferSize))
}






