package main

import (
	"time"

	. "gohipernetFake"
)

type TimerScheduler struct {
	onDone chan struct{}
}

func (scheduler *TimerScheduler) Start() {
	scheduler.onDone = make(chan struct{})
	go scheduler._periodicLoop_goroutine()
}

func (scheduler *TimerScheduler) End() {
	close(scheduler.onDone)
}

func (scheduler *TimerScheduler) _periodicLoop_goroutine() {
	NTELIB_LOG_INFO("Start TimerScheduler goroutine !!!")

	for {
		if scheduler._periodicLoop_goroutine_Impl() {
			NTELIB_LOG_INFO("Wanted Stop TimerScheduler goroutine !!!")
			break
		}
	}

	NTELIB_LOG_INFO("Stop TimerScheduler goroutine !!!")
}

func (scheduler *TimerScheduler) _periodicLoop_goroutine_Impl() bool {
	IsWantedTermination := false

	secondTimeticker := time.NewTicker(time.Second)

	time.Sleep(2 * time.Second) // 순서 종속성으로 인해 2초 뒤 시작한다.
	defer PrintPanicStack()
	defer secondTimeticker.Stop()

	for {
		select {
		/*case secondTime := <-secondTimeticker.C:
		{
			NetLibSetCurrnetUnixTime(secondTime.Unix())
		}*/
		case <-scheduler.onDone:
			{
				IsWantedTermination = true
				break
			}
		}
	}

	return IsWantedTermination
}
