package gamelogic

import (
	"strings"
)

const (
	boardSize        = 32
	redStartSquare   = 1
	redEndSquare     = 12
	blackStartSquare = 21
	blackEndSquare   = 32
	firstRow         = 1
	lastRow          = 8
	firstColumn      = 1
	lastColumn       = 8
	blackManSymbol   = 'b'
	redManSymbol     = 'r'
	blackKingSymbol  = 'B'
	redKingSymbol    = 'R'
	emptySymbol      = ' '
	squareSymbol     = '*'
)

type Board struct {
	turn   PieceColor
	pieces []Piece
}

type BoardMove struct {
	From PieceSquare
	To   PieceSquare
	Hit  bool
}

func InitBoard() Board {
	pieces := make([]Piece, boardSize)
	index := 0
	for i := redStartSquare; i <= redEndSquare; i++ {
		pieces[index] = Piece{
			Color:  RED,
			Type:   MAN,
			Square: PieceSquare(i),
		}
		index++
	}
	for i := blackStartSquare; i <= blackEndSquare; i++ {
		pieces[index] = Piece{
			Color:  BLACK,
			Type:   MAN,
			Square: PieceSquare(i),
		}
		index++
	}
	return Board{
		turn:   BLACK,
		pieces: pieces,
	}
}

func (board Board) Turn() PieceColor {
	return board.turn
}

func (board Board) Pieces() []Piece {
	output := make([]Piece, len(board.pieces))
	copy(output, board.pieces)
	return output
}

func (board Board) String() string {
	var output strings.Builder
	for row := firstRow; row <= lastRow; row++ {
		for col := firstColumn; col <= lastColumn; col++ {
			pos := PiecePosition{row, col}
			if !pos.IsValid() {
				output.WriteRune(emptySymbol)
				continue
			}
			symbol := squareSymbol
			for _, piece := range board.pieces {
				if piece.Square != pos.ToSquare() {
					continue
				}
				switch {
				case piece.Color == BLACK && piece.Type == MAN:
					symbol = blackManSymbol
				case piece.Color == BLACK && piece.Type == KING:
					symbol = blackKingSymbol
				case piece.Color == RED && piece.Type == MAN:
					symbol = redManSymbol
				default:
					symbol = redKingSymbol
				}
				break
			}
			output.WriteRune(symbol)
		}
		if row < lastRow {
			output.WriteRune('\n')
		}
	}
	return output.String()
}

func (board Board) PieceAt(square PieceSquare) (Piece, bool) {
	for _, piece := range board.pieces {
		if piece.Square == square {
			return piece, true
		}
	}
	return Piece{}, false
}

func (board Board) HypotheticalMovesAt(square PieceSquare) ([]BoardMove, bool) {
	piece, exists := board.PieceAt(square)
	if !exists {
		return nil, false
	}
	var directions []PiecePosition
	switch {
	case piece.Type == MAN && piece.Color == BLACK:
		directions = []PiecePosition{{-1, -1}, {-1, 1}}
	case piece.Type == MAN && piece.Color == RED:
		directions = []PiecePosition{{1, -1}, {1, 1}}
	default:
		directions = []PiecePosition{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	}
	current := piece.Square.ToPosition()
	var moves []BoardMove
	for _, direction := range directions {
		next := PiecePosition{
			Row:    current.Row + direction.Row,
			Column: current.Column + direction.Column,
		}
		if !next.IsValid() {
			continue
		}
		nextSquare := next.ToSquare()
		nextPiece, exists := board.PieceAt(nextSquare)
		if !exists {
			moves = append(moves, BoardMove{
				From: piece.Square,
				To:   nextSquare,
				Hit:  false,
			})
			continue
		}
		if nextPiece.Color == piece.Color {
			continue
		}
		doubleNext := PiecePosition{
			Row:    current.Row + direction.Row*2,
			Column: current.Column + direction.Column*2,
		}
		if !doubleNext.IsValid() {
			continue
		}
		doubleNextSquare := doubleNext.ToSquare()
		if _, exists := board.PieceAt(doubleNextSquare); exists {
			continue
		}
		moves = append(moves, BoardMove{
			From: piece.Square,
			To:   doubleNextSquare,
			Hit:  true,
		})
	}
	return moves, true
}

func (board Board) PossibleMovesFor(color PieceColor) []BoardMove {
	var boardMoves []BoardMove
	for i := redStartSquare; i <= blackEndSquare; i++ {
		square := PieceSquare(i)
		piece, exists := board.PieceAt(square)
		if !exists || piece.Color != color {
			continue
		}
		moves, ok := board.HypotheticalMovesAt(square)
		if !ok {
			continue
		}
		boardMoves = append(boardMoves, moves...)
	}
	var hitMoves []BoardMove
	for _, move := range boardMoves {
		if move.Hit {
			hitMoves = append(hitMoves, move)
		}
	}
	if len(hitMoves) > 0 {
		return hitMoves
	}
	return boardMoves
}

func (board Board) CurrentPossibleMoves() []BoardMove {
	return board.PossibleMovesFor(board.turn)
}
