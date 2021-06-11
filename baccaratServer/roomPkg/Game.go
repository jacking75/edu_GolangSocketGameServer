package roomPkg

import (
	"math/rand"
	"time"
)

type baccaratGame struct {
	_statusChangeCompletionMillSec int64 // 다음 상태로 바뀔 때까지의 시간. 0 이면 사용하지 않음.
	_rand                          *rand.Rand
	_dillerCardPos                 int
	_cards                         []int8
}

func (game *baccaratGame) init() {
	game._rand = rand.New(rand.NewSource(1))
	game._rand.Seed(time.Now().UTC().UnixNano()) // 다른 방의 시드 값과 같을 듯. 무작위 난수값이 들어가는 것이 좋음
	game._cards = makeCard()
	game.clear()
}

func (game *baccaratGame) clear() {
	game._statusChangeCompletionMillSec = 0
	game._dillerCardPos = 0
	game._cardChuffle()
}

func (game *baccaratGame) isTimeOver(curTimeMilliSec int64) bool {
	return game._statusChangeCompletionMillSec != 0 && game._statusChangeCompletionMillSec <= curTimeMilliSec
}

func (game *baccaratGame) setBattingWaitTime(curMilliSec int64) {
	game._statusChangeCompletionMillSec = curMilliSec + BATTING_WAIT_MILLISEC
}

func (game *baccaratGame) doBaccarat() baccaratGameResultInfo {
	// 첫번째 카드 배포
	var gameResult baccaratGameResultInfo
	gameResult.clear()

	gameResult.TwoCardtoBanker(game._getDillerCard(), game._getDillerCard())
	gameResult.TwoCardtoPlayer(game._getDillerCard(), game._getDillerCard())

	palyerScore := gameResult.playerScore
	bankerScore := gameResult.bankerScore

	//네추럴 상태로 게임 종료
	if palyerScore == 8 || palyerScore == 9 || bankerScore == 8 || bankerScore == 9 {
		_End(&gameResult)
		return gameResult
	}

	// 3번째 카드
	if palyerScore <= 5 {
		gameResult.ThreeCardtoPlayer(game._getDillerCard())
		palyerScore = gameResult.playerScore
	}

	if (palyerScore == 6 || palyerScore == 7) && bankerScore <= 5 {
		gameResult.ThreeCardtoBanker(game._getDillerCard())
	} else if gameResult.cardsPlayer[2] != -1 {
		player3CardScore := gameResult.cardsPlayer[2]

		switch bankerScore {
		case 0, 1, 2:
			gameResult.ThreeCardtoBanker(game._getDillerCard())
		case 3:
			if player3CardScore != 8 {
				gameResult.ThreeCardtoBanker(game._getDillerCard())
			}
		case 4:
			if 2 <= player3CardScore && player3CardScore <= 7 {
				gameResult.ThreeCardtoBanker(game._getDillerCard())
			}
		case 5:
			if 4 <= player3CardScore && player3CardScore <= 7 {
				gameResult.ThreeCardtoBanker(game._getDillerCard())
			}
		case 6:
			if 6 <= player3CardScore && player3CardScore <= 7 {
				gameResult.ThreeCardtoBanker(game._getDillerCard())
			}
		}

	}

	_End(&gameResult)

	game._statusChangeCompletionMillSec = time.Now().Unix() + NEXT_GAME_WAIT_MILLISEC
	return gameResult
}

// 게임 종료
func _End(endResult *baccaratGameResultInfo) {
	endResult.result = GAME_RESULT_WIN_PLAYER

	if endResult.playerScore == endResult.bankerScore {
		endResult.result = GAME_RESULT_TIE
	} else if endResult.playerScore < endResult.bankerScore {
		endResult.result = GAME_RESULT_WIN_BANKER
	}
}

func (game *baccaratGame) _cardChuffle() {
	// Fisher–Yates shuffle 알고리즘
	// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle

	for i := MAX_CARD_CONT - 1; i >= 0; i-- {
		j := game._rand.Intn(i + 1)
		game._cards[i], game._cards[j] = game._cards[j], game._cards[i]
	}
}

func (game *baccaratGame) _getDillerCard() int8 {
	card := game._cards[game._dillerCardPos]
	game._dillerCardPos++
	return card
}

func (result *baccaratGameResultInfo) clear() {
	result.cardsBanker[0], result.cardsBanker[1], result.cardsBanker[2] = -1, -1, -1
	result.cardsPlayer[0], result.cardsPlayer[1], result.cardsPlayer[2] = -1, -1, -1
	result.playerScore = 0
	result.bankerScore = 0
	result.result = 0
}

func (result *baccaratGameResultInfo) TwoCardtoPlayer(card1 int8, card2 int8) {
	result.cardsPlayer[0] = card1
	result.cardsPlayer[1] = card2
	result.playerScore = (_baccaratCardIndexToScore(card1) + _baccaratCardIndexToScore(card2)) % 10
}

func (result *baccaratGameResultInfo) ThreeCardtoPlayer(card3 int8) {
	result.cardsPlayer[2] = card3
	result.playerScore = (_baccaratCardIndexToScore(result.cardsPlayer[0]) +
		_baccaratCardIndexToScore(result.cardsPlayer[1]) +
		_baccaratCardIndexToScore(result.cardsPlayer[2])) % 10
}

func (result *baccaratGameResultInfo) TwoCardtoBanker(card1 int8, card2 int8) {
	result.cardsBanker[0] = card1
	result.cardsBanker[1] = card2
	result.bankerScore = (_baccaratCardIndexToScore(card1) + _baccaratCardIndexToScore(card2)) % 10
}

func (result *baccaratGameResultInfo) ThreeCardtoBanker(card3 int8) {
	result.cardsBanker[2] = card3
	result.bankerScore = (_baccaratCardIndexToScore(result.cardsBanker[0]) +
		_baccaratCardIndexToScore(result.cardsBanker[1]) +
		_baccaratCardIndexToScore(result.cardsBanker[2])) % 10
}

func _baccaratCardIndexToScore(cardIndex int8) int8 {
	score := cardIndex % CARD_ROW_COUNT

	if score == 0 {
		score = 1
	} else if 10 <= score {
		score = 0
	}

	return score
}
