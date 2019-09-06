package main

import (
	"flag"

	"go.uber.org/zap"

	. "gohipernetFake"
	"main/protocol"
)

func main() {
	NetLibInitLog()

	netConfig, appConfig := parseAppConfig()
	appConfig.writeServerConfig()
	netConfig.WriteNetworkConfig(true)


	protocol.Init_packet()

	NTELIB_LOG_INFO("[[protocolHeaderSize]]",
		zap.Int16("ClientHeaderSize", protocol.ClientHeaderSize()),
		zap.Int16("ServerHeaderSize", protocol.ServerHeaderSize()))


	// 아래 함수를 호출하면 강제적으로 종료 시킬 때까지 대기 상태가 된다.
	createServer(netConfig, appConfig)
}

func parseAppConfig() (NetworkConfig, configAppServer) {
	client := NetworkConfig{}
	appConfig := configAppServer{}

	flag.BoolVar(&client.IsTcp4Addr,"c_IsTcp4Addr", true, "bool flag")
	flag.StringVar(&client.BindAddress,"c_BindAddress", "127.0.0.1:11021", "string flag")
	flag.IntVar(&client.MaxSessionCount,"c_MaxSessionCount", 0, "int flag")
	flag.IntVar(&client.MaxPacketSize,"c_MaxPacketSize", 0, "int flag")
	flag.IntVar(&client.MaxReceiveBufferSize,"c_MaxReceiveBufferSize", 0, "int flag")


	//
	flag.StringVar(&appConfig.GameName,"GameName", "default", "string flag")
	flag.IntVar(&appConfig.RoomMaxCount,"RoomMaxCount", 0, "int flag")
	flag.IntVar(&appConfig.RoomStartNum,"RoomStartNum", 0, "int flag")
	flag.IntVar(&appConfig.RoomMaxUserCount,"RoomMaxUserCount", 0, "RoomMaxUserCount flag")
	flag.IntVar(&appConfig.RoomMaxProcessBufferCount,"RoomMaxProcessBufferCount", 0, "int flag")
	flag.IntVar(&appConfig.RoomCountByGoroutine,"RoomCountByGoroutine", 0, "int flag")
	flag.IntVar(&appConfig.RoomInternalPacketChanBufferCount,"RoomInternalPacketChanBufferCount", 0, "int flag")

	flag.IntVar(&appConfig.CheckCountAtOnce,"CheckCountAtOnce", 0, "int flag")
	flag.IntVar(&appConfig.CheckReriodMillSec,"CheckReriodMillSec", 0, "int flag")
	flag.IntVar(&appConfig.LoginWaitTimeSec,"LoginWaitTimeSec", 0, "int flag")
	flag.IntVar(&appConfig.DisConnectWaitTimeSec,"DisConnectWaitTimeSec", 0, "int flag")
	flag.IntVar(&appConfig.RoomEnterWaitTimeSec,"RoomEnterWaitTimeSec", 0, "int flag")
	flag.IntVar(&appConfig.PingWaitTimeSec,"PingWaitTimeSec", 0, "int flag")
	flag.IntVar(&appConfig.MaxRequestCountPerSecond,"MaxRequestCountPerSecond", 0, "int flag")

	flag.Parse()
	return client, appConfig
}