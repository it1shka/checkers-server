package gamelogic

import "strings"

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
