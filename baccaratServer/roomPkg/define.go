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

// 기존에는 "밀리초" 단위 상수(5000, 10000)를 만들어 놓고 실제로는 time.Now().Unix()(초 단위)에
// 더해서 사용하는 바람에 배팅 대기/결과 표시 시간이 1000배(약 83분/166분)로 늘어나는 버그가 있었다.
// 전체 흐름이 초 단위(Unix())로 통일되어 있으므로 상수도 초 단위로 맞춘다.
const BATTING_WAIT_SEC = 5
const NEXT_GAME_WAIT_SEC = 10

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

