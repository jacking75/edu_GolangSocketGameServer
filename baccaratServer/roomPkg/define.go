package roomPkg

type RoomConfig struct {
	StartRoomNumber int32
	MaxRoomCount int32
	MaxUserCount int32
}


type addRoomUserInfo struct {
	userID []byte

	netSessionIndex     int32
	netSessionUniqueId  uint64
}

// 방의 상태
const (
	ROOM_STATE_NOE = 1
	ROOM_STATE_GAME_WAIT_BATTING = 2
	ROOM_STATE_GAME_RESULT = 3
)

// 카드 정보
const MAX_CARD_CONT = 52
const CARD_ROW_COUNT = 13
// 카드 순서 스페이드, 다이아몬드, 클로버, 하트 A,2,3,4,5,6,7,8,9,10,J,Q,K
func makeCard() []int8 {
	a := make([]int8, MAX_CARD_CONT)
	for i := range a {
		a[i] = (int8)(i)
	}
	return a
}

const BATTING_WAIT_MILLISEC = 5000
const NEXT_GAME_WAIT_MILLISEC = 10000

const (
	BATTING_SELECT_NONE = 0
	BATTING_SELECT_PLAYER = 1
	BATTING_SELECT_BANKER = 2
)

// 게임 결과
const (
	GAME_RESULT_WIN_PLAYER = 1
	GAME_RESULT_WIN_BANKER = 2
	GAME_RESULT_TIE = 3
)

type baccaratGameResultInfo struct {
	cardsBanker [3]int8
	cardsPlayer [3]int8

	playerScore int8
	bankerScore int8

	result int8
}

