package gamelogic

import "strings"

const (
	repetitionRule     = 3
	withoutCaptureRule = 40
)

type GameStatus int

const (
	ACTIVE GameStatus = iota
	BLACK_WIN
	RED_WIN
	DRAW
)

const (
	emptyTag = 'X'
	redTag   = 'R'
	blackTag = 'B'
	delimTag = ';'
)

type GameSession struct {
	board          Board
	status         GameStatus
	states         map[string]int
	withoutCapture int
}

func identify(board Board) string {
	var builder strings.Builder
	if board.turn == BLACK {
		builder.WriteRune(blackTag)
	} else {
		builder.WriteRune(redTag)
	}
	builder.WriteRune(delimTag)
	for i := redStartSquare; i <= blackEndSquare; i++ {
		square := PieceSquare(i)
		piece, exists := board.PieceAt(square)
		var runeTag rune
		switch {
		case !exists:
			runeTag = emptyTag
		case piece.Color == BLACK:
			runeTag = blackTag
		default:
			runeTag = redTag
		}
		builder.WriteRune(runeTag)
		builder.WriteRune(delimTag)
	}
	return builder.String()
}

func NewGameSession() *GameSession {
	board := InitBoard()
	states := make(map[string]int)
	states[identify(board)] = 1
	return &GameSession{
		board:          board,
		status:         ACTIVE,
		states:         states,
		withoutCapture: 0,
	}
}

func (session *GameSession) Board() Board {
	return session.board
}

func (session *GameSession) Status() GameStatus {
	return session.status
}

func (session *GameSession) MakeMove(player PieceColor, from, to PieceSquare) bool {
	if session.board.turn != player || session.status != ACTIVE {
		return false
	}
	nextBoard, ok := session.board.MakeMove(from, to)
	if !ok {
		return false
	}

	prevBoard := session.board
	session.board = nextBoard
	if len(prevBoard.pieces) == len(session.board.pieces) {
		session.withoutCapture++
	} else {
		session.withoutCapture = 0
	}
	stateHash := identify(session.board)
	session.states[stateHash]++

	if len(session.board.CurrentPossibleMoves()) <= 0 {
		if session.board.turn == RED {
			session.status = BLACK_WIN
		} else {
			session.status = RED_WIN
		}
		return true
	}
	if session.withoutCapture >= withoutCaptureRule ||
		session.states[stateHash] >= repetitionRule {
		session.status = DRAW
	}
	return true
}
