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

// TODO: add timers for both players

type Session[P comparable] struct {
	session     *gamelogic.GameSession
	playerBlack P
	playerRed   P
	moves       chan PlayerMove[P]
}

func NewSession[P comparable](playerBlack, playerRed P) *Session[P] {
	return &Session[P]{
		session:     gamelogic.NewGameSession(),
		playerBlack: playerBlack,
		playerRed:   playerRed,
		moves:       make(chan PlayerMove[P]),
	}
}

func (s *Session[P]) playerColor(player P) gamelogic.PieceColor {
	if player == s.playerBlack {
		return gamelogic.BLACK
	}
	return gamelogic.RED
}

func (s *Session[P]) Start(updateHandler func(update Update)) {
	go func() {
		for playerMove := range s.moves {
			color := s.playerColor(playerMove.Player)
			verdict := s.session.MakeMove(color, playerMove.Move.From, playerMove.Move.To)
			if !verdict {
				continue
			}
			updateHandler(Update{
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
