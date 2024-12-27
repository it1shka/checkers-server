package multiplayer

import (
	"sync"

	"it1shka.com/checkers-server/gamelogic"
)

type Session[P comparable] struct {
  mutex sync.Mutex
	session     *gamelogic.GameSession
	playerBlack P
	playerRed   P
}

func NewSession[P comparable](playerBlack, playerRed P) *Session[P] {
  return &Session[P]{
    mutex: sync.Mutex{},
    session: gamelogic.NewGameSession(),
    playerBlack: playerBlack,
    playerRed: playerRed,
  }
}

type MoveResult struct {
  Board gamelogic.Board
  Status gamelogic.GameStatus
}

func (s *Session[P]) MakeMove(player P, move gamelogic.BoardMove) MoveResult {
  
}

type LiveSessions[P comparable] struct {
  mutex sync.RWMutex
	sessions map[P]*Session[P]
}

func NewLiveSessions[P comparable]() *LiveSessions[P] {
  return &LiveSessions[P]{
    mutex: sync.RWMutex{},
    sessions: make(map[P]*Session[P]),
  }
}
