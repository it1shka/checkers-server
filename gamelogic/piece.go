package gamelogic

type PieceColor bool

const (
	BLACK PieceColor = false
	RED   PieceColor = true
)

func (color PieceColor) Opposite() PieceColor {
	return !color
}

func (color PieceColor) String() string {
	switch color {
	case BLACK:
		return "black"
	default:
		return "red"
	}
}

type PieceType bool

const (
	MAN  PieceType = false
	KING PieceType = true
)

func (ptype PieceType) String() string {
	switch ptype {
	case MAN:
		return "man"
	default:
		return "king"
	}
}

type PieceSquare int

func (square PieceSquare) IsValid() bool {
	return square >= 1 && square <= 32
}

func (square PieceSquare) ToPosition() PiecePosition {
	row := (int(square) + 3) / 4
	column := (int(square) - (row-1)*4) * 2
	if row%2 == 0 {
		column--
	}
	return PiecePosition{row, column}
}

type PiecePosition struct {
	Row    int
	Column int
}

func (position PiecePosition) IsValid() bool {
	return (position.Row >= 1 && position.Row <= 8) &&
		(position.Column >= 1 && position.Column <= 8) &&
		(position.Row%2 != position.Column%2)
}

func (position PiecePosition) ToSquare() PieceSquare {
	start := (position.Row-1)*4 + 1
	square := start + (position.Column-1)/2
	return PieceSquare(square)
}

type Piece struct {
	Color  PieceColor
	Type   PieceType
	Square PieceSquare
}
