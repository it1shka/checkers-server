package gamelogic_test

import (
	"strings"
	"testing"

	"it1shka.com/checkers-server/gamelogic"
)

func TestNewGameSession(t *testing.T) {
	session := gamelogic.NewGameSession()
	if session.Status() != gamelogic.ACTIVE {
		t.Fatal("at the start session status should be ACTIVE")
	}
	if session.Board().String() != gamelogic.InitBoard().String() {
		t.Fatalf("at the start session should have an initial board")
	}
}

func TestGameSessionSimulation(t *testing.T) {
	session := gamelogic.NewGameSession()
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
	for index, move := range moves {
		color := gamelogic.BLACK
		if index%2 != 0 {
			color = gamelogic.RED
		}
		result := session.MakeMove(color, move.From, move.To)
		if !result {
			t.Fatal("all moves should succeed")
		}
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
	actualBoard := session.Board().String()
	if session.Board().String() != expectedBoard {
		t.Fatalf("%s\nexpected, found\n%s", expectedBoard, actualBoard)
	}
	if session.Status() != gamelogic.ACTIVE {
		t.Fatalf("status of the session should be ACTIVE")
	}
}
