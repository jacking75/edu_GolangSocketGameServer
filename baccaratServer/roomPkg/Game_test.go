package roomPkg

import (
	"testing"
	"time"
)

// 카드 순서: 스페이드, 다이아몬드, 클로버, 하트 각각 A,2,3,4,5,6,7,8,9,10,J,Q,K(CARD_ROW_COUNT=13장).
// 정규 바카라 룰의 카드 점수: A=1, 2~9=액면가, 10/J/Q/K=0.
var expectedScoreByRank = [CARD_ROW_COUNT]int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0}

func TestBaccaratCardIndexToScore(t *testing.T) {
	for cardIndex := int8(0); cardIndex < int8(MAX_CARD_CONT); cardIndex++ {
		rank := cardIndex % CARD_ROW_COUNT
		want := expectedScoreByRank[rank]

		got := _baccaratCardIndexToScore(cardIndex)
		if got != want {
			t.Errorf("_baccaratCardIndexToScore(%d) [rank=%d] = %d, want %d", cardIndex, rank, got, want)
		}
	}
}

// 배팅 대기/결과 대기 시간은 Unix 초 단위로 일관되게 계산되어야 한다.
// 예전에는 "밀리초" 단위 상수를 초 단위 타임스탬프에 더해 대기 시간이 1000배로 늘어나는 버그가 있었다.
func TestBaccaratGameBattingWaitIsInSeconds(t *testing.T) {
	var game baccaratGame
	now := time.Now().Unix()

	game.setBattingWaitTime(now)

	if game.isTimeOver(now + BATTING_WAIT_SEC - 1) {
		t.Fatalf("isTimeOver should still be false just before BATTING_WAIT_SEC elapses")
	}
	if !game.isTimeOver(now + BATTING_WAIT_SEC) {
		t.Fatalf("isTimeOver should be true once BATTING_WAIT_SEC seconds have elapsed")
	}
}
