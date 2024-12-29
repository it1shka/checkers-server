package multiplayer

import "it1shka.com/checkers-server/gamelogic"

const SESSION_MAX_TIME = 5 * 60

type SessionMove struct {
	Player   string
	From, To uint
}

type SessionUpdate struct {
	Board     gamelogic.Board
	Status    gamelogic.GameStatus
	TimeBlack int
	TimeRed   int
}

type Session struct {
	active                         bool
	session                        *gamelogic.GameSession
	playerBlack, playerRed         string
	playerBlackTime, playerRedTime int
	moves                          chan SessionMove
	done                           chan bool
	updateHandler                  func(update SessionUpdate)
}

func NewSession(playerBlack, playerRed string, updateHandler func(update SessionUpdate)) *Session {
	return &Session{
		active:          true,
		session:         gamelogic.NewGameSession(),
		playerBlack:     playerBlack,
		playerRed:       playerRed,
		playerBlackTime: SESSION_MAX_TIME,
		playerRedTime:   SESSION_MAX_TIME,
		moves:           make(chan SessionMove),
		done:            make(chan bool),
		updateHandler:   updateHandler,
	}
}

func (s *Session) playerColor(player string) gamelogic.PieceColor {
	if player == s.playerBlack {
		return gamelogic.BLACK
	}
	return gamelogic.RED
}

func (s *Session) Start() {
	go func() {
		for {
			select {
			case move := <-s.moves:
				playerColor := s.playerColor(move.Player)
				from := gamelogic.PieceSquare(move.From)
				to := gamelogic.PieceSquare(move.To)
				if from.IsValid() && to.IsValid() {
					s.session.MakeMove(playerColor, from, to)
					s.updateHandler(SessionUpdate{})
				}
			case <-s.done:
				return
			}
		}
	}()
}

func (s *Session) Move(move SessionMove) {
	if !s.active {
		return
	}
	s.moves <- move
}

func (s *Session) SurrenderAndStop(player string) {
	if player == s.playerBlack {
		s.updateHandler(SessionUpdate{
			Board:     s.session.Board(),
			Status:    gamelogic.RED_WIN,
			TimeBlack: s.playerBlackTime,
			TimeRed:   s.playerRedTime,
		})
	} else {
		s.updateHandler(SessionUpdate{
			Board:     s.session.Board(),
			Status:    gamelogic.BLACK_WIN,
			TimeBlack: s.playerBlackTime,
			TimeRed:   s.playerRedTime,
		})
	}
	s.Stop()
}

func (s *Session) Stop() {
	s.active = false
	close(s.done)
	close(s.moves)
}
