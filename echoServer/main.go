package main

import (
	"flag"
	. "gohipernetFake"
)


func main() {
	NetLibInitLog(LOG_LEVEL_DEBUG, nil)

	netConfigClient := parseAppConfig()
	netConfigClient.WriteNetworkConfig(true)

	// 아래 함수를 호출하면 강제적으로 종료 시킬 때까지 대기 상태가 된다.
	createServer(netConfigClient)
}

func parseAppConfig() NetworkConfig {
	client := NetworkConfig{}

	flag.BoolVar(&client.IsTcp4Addr,"c_IsTcp4Addr", true, "bool flag")
	flag.StringVar(&client.BindAddress,"c_BindAddress", "127.0.0.1:11021", "string flag")
	flag.IntVar(&client.MaxSessionCount,"c_MaxSessionCount", 0, "int flag")
	flag.IntVar(&client.MaxPacketSize,"c_MaxPacketSize", 0, "int flag")

	flag.Parse()
	return client
}

