package gamelogic

import (
	"strings"
)

const (
	initialPieceAmount   = 24
	redStartSquare       = 1
	redEndSquare         = 12
	blackStartSquare     = 21
	blackEndSquare       = 32
	firstRow             = 1
	lastRow              = 8
	firstColumn          = 1
	lastColumn           = 8
	blackKingStartSquare = 1
	blackKingEndSquare   = 4
	redKingStartSquare   = 29
	redKingEndSquare     = 32
	blackManSymbol       = 'b'
	redManSymbol         = 'r'
	blackKingSymbol      = 'B'
	redKingSymbol        = 'R'
	emptySymbol          = ' '
	squareSymbol         = '*'
)

type Board struct {
	turn          PieceColor
	pieces        []Piece
	multijump     bool
	multijumpFrom PieceSquare
}

type BoardMove struct {
	From PieceSquare
	To   PieceSquare
	Hit  bool
}

func UnsafeInitBoard(turn PieceColor, pieces []Piece) Board {
	return Board{
		turn:          turn,
		pieces:        pieces,
		multijump:     false,
		multijumpFrom: -1,
	}
}

func InitBoard() Board {
	pieces := make([]Piece, initialPieceAmount)
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
		turn:          BLACK,
		pieces:        pieces,
		multijump:     false,
		multijumpFrom: -1,
	}
}

func (board Board) Copy() Board {
	return Board{
		turn:          board.Turn(),
		pieces:        board.Pieces(),
		multijump:     board.multijump,
		multijumpFrom: board.multijumpFrom,
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

func (board Board) UnsafePieces() []Piece {
	return board.pieces
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
	if board.multijump {
		multijumpMoves, ok := board.HypotheticalMovesAt(board.multijumpFrom)
		if !ok {
			return []BoardMove{}
		}
		var hitMoves []BoardMove
		for _, move := range multijumpMoves {
			if !move.Hit {
				continue
			}
			hitMoves = append(hitMoves, move)
		}
		return hitMoves
	}
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

func (board Board) MakeMove(from, to PieceSquare) (Board, bool) {
	possibleMoves := board.CurrentPossibleMoves()
	possible := false
	hit := false
	for _, move := range possibleMoves {
		if move.From == from && move.To == to {
			possible = true
			hit = move.Hit
			break
		}
	}
	if !possible {
		return Board{}, false
	}
	fromPosition := from.ToPosition()
	toPosition := to.ToPosition()
	middleSquare := PiecePosition{
		Row:    (fromPosition.Row + toPosition.Row) / 2,
		Column: (fromPosition.Column + toPosition.Column) / 2,
	}.ToSquare()
	var nextPieces []Piece
	var movingPiece Piece
	for _, piece := range board.pieces {
		if piece.Square == from {
			movingPiece = piece
			continue
		}
		if hit && piece.Square == middleSquare {
			continue
		}
		nextPieces = append(nextPieces, piece)
	}
	movingPiece.Square = to
	if (movingPiece.Color == BLACK && to >= blackKingStartSquare && to <= blackKingEndSquare) ||
		(movingPiece.Color == RED && to >= redKingStartSquare && to <= redKingEndSquare) {
		movingPiece.Type = KING
	}
	nextPieces = append(nextPieces, movingPiece)
	nextBoard := Board{
		turn:          board.turn,
		pieces:        nextPieces,
		multijump:     true,
		multijumpFrom: to,
	}
	nextTurn := true
	if hit {
		nextPossibleMoves := nextBoard.CurrentPossibleMoves()
		if len(nextPossibleMoves) > 0 && nextPossibleMoves[0].Hit {
			nextTurn = false
		}
	}
	if nextTurn {
		nextBoard.turn = nextBoard.turn.Opposite()
		nextBoard.multijump = false
		nextBoard.multijumpFrom = -1
	}
	return nextBoard, true
}
