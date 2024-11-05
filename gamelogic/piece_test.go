package gamelogic_test

import (
	"testing"

	"it1shka.com/checkers-server/gamelogic"
)

func TestPieceSquareValidity(t *testing.T) {
	valid := []gamelogic.PieceSquare{1, 5, 31, 32, 13}
	for _, square := range valid {
		if !square.IsValid() {
			t.Fatalf("%d should be a valid square", square)
		}
	}

	invalid := []gamelogic.PieceSquare{-1, 0, 33, 67, 199}
	for _, square := range invalid {
		if square.IsValid() {
			t.Fatalf("%d should be an invalid square", square)
		}
	}
}

func TestPiecePositionValidity(t *testing.T) {
	valid := []gamelogic.PiecePosition{
		{1, 2}, {1, 4},
		{3, 4}, {3, 6},
		{6, 1}, {6, 7},
	}
	for _, position := range valid {
		if !position.IsValid() {
			t.Fatalf("%v should be valid", position)
		}
	}

	invalid := []gamelogic.PiecePosition{
		{-1, 0}, {9, 9}, {1, 1},
		{1, 3}, {2, 2}, {2, 4},
		{4, 4}, {5, 5}, {100, 100},
		{8, 9},
	}
	for _, position := range invalid {
		if position.IsValid() {
			t.Fatalf("%v should be invalid", position)
		}
	}
}

func TestSquareToPosition(t *testing.T) {
	cases := []struct {
		Square   gamelogic.PieceSquare
		Position gamelogic.PiecePosition
	}{
		{1, gamelogic.PiecePosition{1, 2}},
		{3, gamelogic.PiecePosition{1, 6}},
		{10, gamelogic.PiecePosition{3, 4}},
		{16, gamelogic.PiecePosition{4, 7}},
		{25, gamelogic.PiecePosition{7, 2}},
	}
	for _, testCase := range cases {
		output := testCase.Square.ToPosition()
		if output.Row != testCase.Position.Row ||
			output.Column != testCase.Position.Column {
			t.Fatalf(
				"square %d should be at %v, but the output was %v",
				testCase.Square,
				testCase.Position,
				output,
			)
		}
	}
}

func TestPositionToSquare(t *testing.T) {
	cases := []struct {
		Position gamelogic.PiecePosition
		Square   gamelogic.PieceSquare
	}{
		{gamelogic.PiecePosition{1, 2}, 1},
		{gamelogic.PiecePosition{5, 6}, 19},
		{gamelogic.PiecePosition{7, 6}, 27},
		{gamelogic.PiecePosition{8, 7}, 32},
	}
	for _, testCase := range cases {
		output := testCase.Position.ToSquare()
		if output != testCase.Square {
			t.Fatalf(
				"position %v should be at square %d, but the output was %d",
				testCase.Position,
				testCase.Square,
				output,
			)
		}
	}
}

func TestOppositeColor(t *testing.T) {
	testCases := []struct {
		Current  gamelogic.PieceColor
		Opposite gamelogic.PieceColor
	}{
		{gamelogic.BLACK, gamelogic.RED},
		{gamelogic.RED, gamelogic.BLACK},
	}
	for _, testCase := range testCases {
		if testCase.Current.Opposite() != testCase.Opposite {
			t.Fatal("expected different opposite color")
		}
	}
}
