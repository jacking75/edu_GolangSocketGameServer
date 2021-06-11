package gohipernetFake

import (
	"fmt"
)


type NetworkConfig struct {
	IsTcp4Addr                    bool
	BindAddress                   string // 예) localhost:19999
	MaxSessionCount               int // 최대 클라이언트 세션 수. 넉넉하게 많이 해도 괜찮다
	MaxPacketSize                 int // 최대 패킷 크기
	MaxReceiveBufferSize          int // 받기 버퍼 크기. 최소 MaxPacketSize 2배 이상 추천.

}

func (config NetworkConfig) WriteNetworkConfig(isClientSide bool) {
	logInfo("", 0, fmt.Sprintf("config - isClientSide: %t", isClientSide))
	logInfo("", 0, fmt.Sprintf("config - IsTcp4Addr: %t", config.IsTcp4Addr))
	logInfo("", 0, fmt.Sprintf("config - ClientAddress: %s", config.BindAddress))
	logInfo("", 0, fmt.Sprintf("config - MaxSessionCount: %d", config.MaxSessionCount))
	logInfo("", 0, fmt.Sprintf("config - MaxPacketSize: %d", config.MaxPacketSize))
	logInfo("", 0, fmt.Sprintf("config - MaxReceiveBufferSize: %d", config.MaxReceiveBufferSize))
}






