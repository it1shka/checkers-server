package gamelogic

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestIdentifyBoardFormat(t *testing.T) {
	board := InitBoard()
	identifier := identify(board)
	expected := "1rrrrrrrrrrrrXXXXXXXXbbbbbbbbbbbb"
	if identifier != expected {
		t.Fatalf("initial hash should be %s, found %s", expected, identifier)
	}
}

func TestIdentifyBoard(t *testing.T) {
	dictionary := make(map[string]string)
	var board Board
	setBoard := func(newBoard Board) {
		board = newBoard
		turn := "black"
		if board.turn == RED {
			turn = "red"
		}
		identifier := identify(board)
		currentPrint := fmt.Sprintf("%s\n%s", board.String(), turn)
		previousPrint, exists := dictionary[identifier]
		if !exists {
			dictionary[identifier] = currentPrint
			return
		}
		if previousPrint != currentPrint {
			t.Fatalf("Boards:\n\n%s\n\nand\n\n%s\n\nhave the same hash:\n%s", previousPrint, currentPrint, identifier)
		}
	}

	setBoard(InitBoard())

	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 1_000_000; i++ {
		available := board.CurrentPossibleMoves()
		if len(available) <= 0 {
			setBoard(InitBoard())
			continue
		}
		move := available[rnd.Intn(len(available))]
		nextBoard, ok := board.MakeMove(move.From, move.To)
		if !ok {
			setBoard(InitBoard())
			continue
		}
		setBoard(nextBoard)
	}

	t.Logf("checked %d hashes for boards", len(dictionary))
}
