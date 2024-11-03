package gamelogic_test

import (
	"strings"
	"testing"

	"it1shka.com/checkers-server/gamelogic"
)

func TestBoardInit(t *testing.T) {
	board := gamelogic.InitBoard()
	expected := strings.Join([]string{
		" R R R R",
		"R R R R ",
		" R R R R",
		"* * * * ",
		" * * * *",
		"B B B B ",
		" B B B B",
		"B B B B ",
	}, "\n")
	if board.String() != expected {
		t.Fatalf("%s\nexpected, found\n%s", expected, board)
	}
}
