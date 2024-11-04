package gamelogic_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"it1shka.com/checkers-server/gamelogic"
)

func TestBoardInit(t *testing.T) {
	board := gamelogic.InitBoard()
	expected := strings.Join([]string{
		" r r r r",
		"r r r r ",
		" r r r r",
		"* * * * ",
		" * * * *",
		"b b b b ",
		" b b b b",
		"b b b b ",
	}, "\n")
	if board.String() != expected {
		t.Fatalf("%s\nexpected, found\n%s", expected, board)
	}
}

func TestBoardPieceAt(t *testing.T) {
	board := gamelogic.InitBoard()
	testCases := []struct {
		Square gamelogic.PieceSquare
		Color  gamelogic.PieceColor
		Exists bool
	}{
		{Square: 1, Color: gamelogic.RED, Exists: true},
		{Square: 2, Color: gamelogic.RED, Exists: true},
    {Square: 13, Exists: false},
    {Square: 14, Exists: false},
    {Square: 21, Color: gamelogic.BLACK, Exists: true},
    {Square: 25, Color: gamelogic.BLACK, Exists: true},
	}
  for _, testCase := range testCases {
    piece, exists := board.PieceAt(testCase.Square)
    switch {
    case !testCase.Exists && exists:
      t.Fatalf(
        "expected %d square to be empty", 
        testCase.Square,
      )
    case testCase.Exists && !exists:
      t.Fatalf(
        "expected %d square to contain a piece",
        testCase.Square,
      )
    case testCase.Exists && exists && testCase.Color != piece.Color:
      t.Fatalf(
        "expected different color at %d square",
        testCase.Square,
      )
    }
  }
}

func TestHypotheticalMovesAt(t *testing.T) {
  board := gamelogic.InitBoard()
  testCases := []struct {
    Square gamelogic.PieceSquare
    Moves []gamelogic.PieceSquare
  }{
    {1, []gamelogic.PieceSquare{}},
    {3, []gamelogic.PieceSquare{}},
    {7, []gamelogic.PieceSquare{}},
    {30, []gamelogic.PieceSquare{}},
    {27, []gamelogic.PieceSquare{}},
    {11, []gamelogic.PieceSquare{15, 16}},
    {10, []gamelogic.PieceSquare{14, 15}},
  }
  for _, testCase := range testCases {
    pieceMoves, _ := board.HypotheticalMovesAt(testCase.Square)
    squareMoves := make([]gamelogic.PieceSquare, len(pieceMoves))
    for i := 0; i < len(pieceMoves); i++ {
      squareMoves[i] = pieceMoves[i].To
    }
    slices.Sort(squareMoves)
    if !reflect.DeepEqual(squareMoves, testCase.Moves) {
      t.Fatalf(
        "wrong moves at %d square: found %v, expected %v", 
        testCase.Square, 
        squareMoves,
        testCase.Moves,
      )
    }
  }
}
