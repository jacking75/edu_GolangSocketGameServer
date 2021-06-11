package main

import (
	"flag"
	. "gohipernetFake"
)


func main() {
	NetLibInitLog(LOG_LEVEL_DEBUG, nil)

	netConfig, appConfig := parseAppConfig()
	netConfig.WriteNetworkConfig(true)

	// 아래 함수를 호출하면 강제적으로 종료 시킬 때까지 대기 상태가 된다.
	createAnsStartServer(netConfig, appConfig)
}

func parseAppConfig() (NetworkConfig, configAppServer) {
	OutPutLog(LOG_LEVEL_INFO,"", 0,"[[Setting NetworkConfig]]")

	//TODO flag 사용하기
	appConfig := configAppServer {
		"chatServer",
		1000,
		0,
		4,
	}


	netConfig := NetworkConfig{}

	flag.BoolVar(&netConfig.IsTcp4Addr,"c_IsTcp4Addr", true, "bool flag")
	flag.StringVar(&netConfig.BindAddress,"c_BindAddress", "127.0.0.1:11021", "string flag")
	flag.IntVar(&netConfig.MaxSessionCount,"c_MaxSessionCount", 0, "int flag")
	flag.IntVar(&netConfig.MaxPacketSize,"c_MaxPacketSize", 0, "int flag")

	flag.Parse()
	return netConfig, appConfig
}

