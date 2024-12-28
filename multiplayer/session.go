package multiplayer

import "it1shka.com/checkers-server/gamelogic"

const SESSION_MAX_TIME = 5 * 60

type SessionUpdate struct {
	Board  gamelogic.Board
	Status gamelogic.GameStatus
}

type Session[ID comparable] struct {
	session                        *gamelogic.GameSession
	playerBlack, playerRed         ID
	playerBlackTime, playerRedTime int
}

func NewSession[ID comparable](playerBlack, playerRed ID) *Session[ID] {
	return &Session[ID]{
		session:         gamelogic.NewGameSession(),
		playerBlack:     playerBlack,
		playerRed:       playerRed,
		playerBlackTime: SESSION_MAX_TIME,
		playerRedTime:   SESSION_MAX_TIME,
	}
}

func (s *Session[ID]) Start(updateHandler func(update SessionUpdate)) {
	// TODO:
}

func (s *Session[ID]) SurrenderAndStop(player ID) {
	// TODO:
}

func (s *Session[ID]) Stop() {
	// TODO:
}
