package gamelogic_test

import (
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

func Test(t *testing.T) {

}
