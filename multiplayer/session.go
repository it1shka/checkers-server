package multiplayer

import "it1shka.com/checkers-server/gamelogic"

// Helper records

type PlayerMove[P any] struct {
	Player P
	Move   gamelogic.BoardMove
}

type Update struct {
	Board  gamelogic.Board
	Status gamelogic.GameStatus
}

// Session implementation

type Session[P comparable] struct {
	session     *gamelogic.GameSession
	playerBlack P
	playerRed   P
	moves       chan PlayerMove[P]
	updateHook  func(update Update)
}

func NewSession[P comparable](playerBlack, playerRed P, updateHook func(update Update)) *Session[P] {
	return &Session[P]{
		session:     gamelogic.NewGameSession(),
		playerBlack: playerBlack,
		playerRed:   playerRed,
		moves:       make(chan PlayerMove[P]),
		updateHook:  updateHook,
	}
}

func (s *Session[P]) playerColor(player P) gamelogic.PieceColor {
	if player == s.playerBlack {
		return gamelogic.BLACK
	}
	return gamelogic.RED
}

func (s *Session[P]) Start() {
	go func() {
		for playerMove := range s.moves {
			color := s.playerColor(playerMove.Player)
			verdict := s.session.MakeMove(color, playerMove.Move.From, playerMove.Move.To)
			if !verdict {
				continue
			}
			s.updateHook(Update{
				Board:  s.session.Board(),
				Status: s.session.Status(),
			})
		}
	}()
}

func (s *Session[P]) Stop() {
	close(s.moves)
}

func (s *Session[P]) MakeMove(player P, move gamelogic.BoardMove) {
	s.moves <- PlayerMove[P]{
		Player: player,
		Move:   move,
	}
}
