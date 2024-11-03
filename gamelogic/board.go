package gamelogic

import "strings"

const (
	BOARD_SIZE         = 32
	RED_START_SQUARE   = 1
	RED_END_SQUARE     = 12
	BLACK_START_SQUARE = 21
	BLACK_END_SQUARE   = 32
	FIRST_ROW          = 1
	LAST_ROW           = 8
	FIRST_COLUMN       = 1
	LAST_COLUMN        = 8
	BLACK_MARK         = 'B'
	RED_MARK           = 'R'
	EMPTY_MARK         = ' '
	SQUARE_MARK        = '*'
)

type Board struct {
	turn   PieceColor
	pieces []Piece
}

func InitBoard() Board {
	pieces := make([]Piece, BOARD_SIZE)
	index := 0
	for i := RED_START_SQUARE; i <= RED_END_SQUARE; i++ {
		pieces[index] = Piece{
			Color:  RED,
			Type:   MAN,
			Square: PieceSquare(i),
		}
		index++
	}
	for i := BLACK_START_SQUARE; i <= BLACK_END_SQUARE; i++ {
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
	for row := FIRST_ROW; row <= LAST_ROW; row++ {
		for col := FIRST_COLUMN; col <= LAST_COLUMN; col++ {
			pos := PiecePosition{row, col}
			if !pos.IsValid() {
				output.WriteRune(EMPTY_MARK)
				continue
			}
			empty := true
			for _, piece := range board.pieces {
				if pos.ToSquare() == piece.Square {
					empty = false
					if piece.Color == BLACK {
						output.WriteRune(BLACK_MARK)
					} else {
						output.WriteRune(RED_MARK)
					}
					break
				}
			}
			if empty {
				output.WriteRune(SQUARE_MARK)
			}
		}
		if row < LAST_ROW {
			output.WriteRune('\n')
		}
	}
	return output.String()
}
