package main

import (
	. "gohipernetFake"
)

type configAppServer struct {
	GameName                   string

	RoomMaxCount               int
	RoomStartNum                   int
	RoomMaxUserCount                  int
	RoomMaxProcessBufferCount                 int
	RoomCountByGoroutine          int
	RoomInternalPacketChanBufferCount     int

	CheckCountAtOnce int
	CheckReriodMillSec int
	LoginWaitTimeSec         int
	DisConnectWaitTimeSec    int
	RoomEnterWaitTimeSec     int
	PingWaitTimeSec          int
	MaxRequestCountPerSecond int
}

func (config configAppServer) writeServerConfig() {
	NTELIB_LOG_INFO("writeServerConfig")
	/*NTELIB_LOG_INFO("writeServerConfig - " + config.GameName,
		zap.Int("RoomMaxCount", config.RoomMaxCount),
		zap.Int("RoomStartNum", config.RoomStartNum),
		zap.Int("RoomMaxUserCount", config.RoomMaxUserCount),
		zap.Int("RoomMaxProcessBufferCount", config.RoomMaxProcessBufferCount),
		zap.Int("RoomCountByGoroutine", config.RoomCountByGoroutine),
		zap.Int("RoomInternalPacketChanBufferCount", config.RoomInternalPacketChanBufferCount),
		zap.Int("CheckCountAtOnce", config.CheckCountAtOnce),
		zap.Int("CheckReriodMillSec", config.CheckReriodMillSec),
		zap.Int("LoginWaitTimeSec", config.LoginWaitTimeSec),
		zap.Int("DisConnectWaitTimeSec", config.DisConnectWaitTimeSec),
		zap.Int("RoomEnterWaitTimeSec", config.RoomEnterWaitTimeSec),
		zap.Int("PingWaitTimeSec", config.PingWaitTimeSec),
		zap.Int("MaxRequestCountPerSecond", config.MaxRequestCountPerSecond))*/
}