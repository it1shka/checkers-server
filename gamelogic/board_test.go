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

func TestPieceAt(t *testing.T) {
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
		Moves  []gamelogic.PieceSquare
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

func TestPossibleMovesFor(t *testing.T) {
	board := gamelogic.InitBoard()
	testCases := []struct {
		Color gamelogic.PieceColor
		Moves []gamelogic.BoardMove
	}{
		{gamelogic.BLACK, []gamelogic.BoardMove{
			{From: 21, To: 17},
			{From: 22, To: 17},
			{From: 22, To: 18},
			{From: 23, To: 18},
			{From: 23, To: 19},
			{From: 24, To: 19},
			{From: 24, To: 20},
		}},
		{gamelogic.RED, []gamelogic.BoardMove{
			{From: 9, To: 13},
			{From: 9, To: 14},
			{From: 10, To: 14},
			{From: 10, To: 15},
			{From: 11, To: 15},
			{From: 11, To: 16},
			{From: 12, To: 16},
		}},
	}
	for _, testCase := range testCases {
		moves := board.PossibleMovesFor(testCase.Color)
		slices.SortFunc(moves, func(a, b gamelogic.BoardMove) int {
			if a.From == b.From {
				return int(a.To - b.To)
			}
			return int(a.From - b.From)
		})
		if !reflect.DeepEqual(moves, testCase.Moves) {
			t.Fatalf("%v expected, %v found", testCase.Moves, moves)
		}
	}
}

func TestMakeMove(t *testing.T) {
	board := gamelogic.InitBoard()
	moves := []struct {
		From gamelogic.PieceSquare
		To   gamelogic.PieceSquare
	}{
		{21, 17},
		{10, 14},
		{17, 10},
		{7, 14},
		{24, 19},
		{11, 15},
		{19, 10},
	}
	for _, move := range moves {
		nextBoard, ok := board.MakeMove(move.From, move.To)
		t.Logf("Board:\n%s\n", nextBoard)
		t.Logf("Next turn: %s\n", nextBoard.Turn())
		if !ok {
			t.Fatalf(
				"move from %d to %d failed",
				move.From,
				move.To,
			)
		}
		board = nextBoard
	}
	expectedBoard := strings.Join([]string{
		" r r r r",
		"r r * r ",
		" r b * r",
		"* r * * ",
		" * * * *",
		"* b b * ",
		" b b b b",
		"b b b b ",
	}, "\n")
	if board.String() != expectedBoard {
		t.Fatalf("%s\nexpected, found\n%s", expectedBoard, board)
	}
}
